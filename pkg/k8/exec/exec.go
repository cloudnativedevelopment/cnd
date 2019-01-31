package exec

import (
	"context"
	"fmt"
	"io"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	kexec "k8s.io/kubernetes/pkg/kubectl/cmd/exec"
)

// Exec executes the command in the cnd container
func Exec(ctx context.Context, c *kubernetes.Clientset, config *rest.Config, pod *apiv1.Pod, container string, tty bool, stdin io.Reader, stdout, stderr io.Writer, command []string) error {

	if pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodFailed {
		return fmt.Errorf("cannot exec into a cloud native development environment that's not active; current status is %s", pod.Status.Phase)
	}

	p := &kexec.ExecOptions{}
	p.Config = config
	p.Command = command
	p.Executor = &kexec.DefaultRemoteExecutor{}
	p.IOStreams = genericclioptions.IOStreams{In: stdin, Out: stdout, ErrOut: stderr}
	p.Stdin = true
	p.TTY = tty

	t := p.SetupTTY()

	var sizeQueue remotecommand.TerminalSizeQueue
	if t.Raw {
		// this call spawns a goroutine to monitor/update the terminal size
		sizeQueue = t.MonitorSize(t.GetSize())

		// unset p.Err if it was previously set because both stdout and stderr go over p.Out when tty is
		// true
		p.ErrOut = nil
	}

	fn := func() error {
		req := c.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(pod.Name).
			Namespace(pod.Namespace).
			SubResource("exec").
			Param("container", container)
		req.VersionedParams(&apiv1.PodExecOptions{
			Container: container,
			Command:   p.Command,
			Stdin:     p.Stdin,
			Stdout:    p.Out != nil,
			Stderr:    p.ErrOut != nil,
			TTY:       t.Raw,
		}, scheme.ParameterCodec)

		req.Context(ctx)
		return p.Executor.Execute("POST", req.URL(), config, p.In, p.Out, p.ErrOut, t.Raw, sizeQueue)
	}

	if err := t.Safe(fn); err != nil {
		return err
	}

	return nil
}
