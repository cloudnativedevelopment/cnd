package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/okteto/app/cli/pkg/config"
	"github.com/okteto/app/cli/pkg/k8s/exec"
	"github.com/okteto/app/cli/pkg/k8s/pods"
	"github.com/okteto/app/cli/pkg/log"
	"github.com/okteto/app/cli/pkg/model"

	k8Client "github.com/okteto/app/cli/pkg/k8s/client"

	"github.com/spf13/cobra"
)

//Exec executes a command on the CND container
func Exec() *cobra.Command {
	var devPath string
	var pod string
	var container string
	var namespace string
	var port int

	cmd := &cobra.Command{
		Use:   "exec COMMAND",
		Short: "Execute a command in your Okteto Environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			devPath = getFullPath(devPath)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if port != 0 {
				go func() {
					http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						log.Debug("canceling process due to a request")
						cancel()
					})

					log.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
					panic("webserver stopped handling requests")
				}()
			}

			if _, err := os.Stat(devPath); os.IsNotExist(err) {
				return fmt.Errorf("'%s' does not exist", devPath)
			}

			dev, err := model.Get(devPath)
			if err != nil {
				return err
			}
			if namespace != "" {
				dev.Namespace = namespace
			}
			err = executeExec(ctx, pod, container, dev, args)
			return err
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("exec requires the COMMAND argument")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&devPath, "file", "f", config.ManifestFileName(), "path to the manifest file")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace where the exec command is executed")
	cmd.Flags().StringVarP(&pod, "pod", "p", "", "pod where it is executed")
	cmd.Flags().MarkHidden("pod")
	cmd.Flags().StringVarP(&container, "container", "c", "", "container where it is executed")
	cmd.Flags().MarkHidden("container")
	cmd.Flags().IntVar(&port, "port", 0, "port to listen to signals")
	cmd.Flags().MarkHidden("port")

	return cmd
}

func executeExec(ctx context.Context, pod, container string, dev *model.Dev, args []string) error {
	client, cfg, namespace, err := k8Client.GetLocal()
	if err != nil {
		return err
	}

	if dev.Namespace == "" {
		dev.Namespace = namespace
	}

	if len(pod) == 0 {
		p, err := pods.GetDevPod(ctx, dev, client)
		if err != nil {
			return err
		}

		pod = p.Name
		if len(dev.Container) == 0 {
			dev.Container = p.Spec.Containers[0].Name
		}
	}

	if len(container) > 0 {
		dev.Container = container
	}
	return exec.Exec(ctx, client, cfg, dev.Namespace, pod, dev.Container, true, os.Stdin, os.Stdout, os.Stderr, args)
}
