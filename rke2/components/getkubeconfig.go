package components

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
	"strings"
)

func NewGetKubeConfig(ctx *pulumi.Context, name string, node Node, path string, fileName string, sshuser SshUser, opts ...pulumi.ResourceOption) (*GetKubeConfig, error) {
	GetKubeConfig := &GetKubeConfig{}
	err := ctx.RegisterComponentResource("rancher:rke2:GetKubeConfig", name, GetKubeConfig, opts...)
	if err != nil {
		return GetKubeConfig, err
	}

	connection := remote.ConnectionArgs{
		Host:       pulumi.String(node.IP),
		User:       pulumi.String(sshuser.Name),
		PrivateKey: pulumi.String(sshuser.PrivateKey),
	}

	kubeConfig, err := remote.NewCommand(ctx, "get-kubeconfig-"+node.Name, &remote.CommandArgs{
		Create:     pulumi.String("cat /etc/rancher/rke2/rke2.yaml"),
		Update:     pulumi.String("cat /etc/rancher/rke2/rke2.yaml"),
		Connection: &connection,
	}, pulumi.Parent(GetKubeConfig))
	if err != nil {
		return GetKubeConfig, err
	}

	kubeConfig.Stdout.ApplyT(func(kubeConfig string) error {
		kubeConfig = strings.Replace(kubeConfig, "127.0.0.1", node.IP, -1)
		err := os.MkdirAll(path, 0750)
		if err != nil {
			return err
		}
		err = os.WriteFile(path+"/"+fileName, []byte(kubeConfig), 0644)
		return err
	})
	if err != nil {
		return GetKubeConfig, err
	}

	return GetKubeConfig, nil
}
