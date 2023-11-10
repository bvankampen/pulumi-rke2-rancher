package components

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewRunRKE2Installer(ctx *pulumi.Context, name string, nodes []Node, sshuser SshUser, opts ...pulumi.ResourceOption) (*RunRKE2Installer, error) {
	RunRKE2Installer := &RunRKE2Installer{}
	err := ctx.RegisterComponentResource("rancher:rke2:RunRKE2Installer", name, RunRKE2Installer, opts...)
	if err != nil {
		return RunRKE2Installer, err
	}

	var dependsOn []pulumi.Resource

	for _, node := range nodes {
		connection := remote.ConnectionArgs{
			Host:       pulumi.String(node.IP),
			User:       pulumi.String(sshuser.Name),
			PrivateKey: pulumi.String(sshuser.PrivateKey),
		}

		runRKE2Install, err := remote.NewCommand(ctx, "install-rke2-"+node.Name, &remote.CommandArgs{
			Create:     pulumi.String("sudo bash /opt/rke2/install/run-install-rke2.sh"),
			Connection: &connection,
		}, pulumi.Parent(RunRKE2Installer), pulumi.DependsOn(dependsOn))
		if err != nil {
			return RunRKE2Installer, err
		}

		waitForRKE2, err := remote.NewCommand(ctx, "wait-for-rke2-"+node.Name, &remote.CommandArgs{
			Create:     pulumi.String("bash /opt/rke2/install/wait-for-rke2.sh"),
			Connection: &connection,
		}, pulumi.Parent(RunRKE2Installer), pulumi.DependsOn([]pulumi.Resource{runRKE2Install}))
		if err != nil {
			return RunRKE2Installer, err
		}

		dependsOn = []pulumi.Resource{waitForRKE2}
	}

	return RunRKE2Installer, nil
}
