package cmd

import (
	"fmt"

	"github.com/okteto/cnd/pkg/storage"
	"github.com/okteto/cnd/pkg/syncthing"

	"github.com/okteto/cnd/pkg/k8/client"
	"github.com/okteto/cnd/pkg/k8/deployments"
	"github.com/okteto/cnd/pkg/model"
	"github.com/spf13/cobra"
)

//Down stops a cloud native environment
func Down() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Deactivates your cloud native development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDown(devPath)
		},
	}

	return cmd
}

func executeDown(devPath string) error {
	fmt.Println("Deactivating your cloud native development environment...")
	dev, err := model.ReadDev(devPath)
	if err != nil {
		return err
	}

	namespace, client, _, err := client.Get()
	if err != nil {
		return err
	}

	name, err := deployments.Deploy(dev, namespace, client)
	if err != nil {
		return err
	}

	syncthing, err := syncthing.NewSyncthing(name, namespace, "")
	if err != nil {
		return err
	}

	storage.Delete(namespace, name)

	err = syncthing.Stop()
	if err != nil {
		return err
	}

	fmt.Println("Cloud native development environment deactivated...")
	return nil
}
