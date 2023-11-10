package components

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewRemoteCommands(ctx *pulumi.Context, name string, nodes []Node, commandList []RemoteCommandsArguments, sshuser SshUser, opts ...pulumi.ResourceOption) (*RemoteCommands, error) {
	RemoteCommandList := &RemoteCommands{}

	err := ctx.RegisterComponentResource("rancher:rke2:RemoteCommands", name, RemoteCommandList, opts...)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		var dependsOn []pulumi.Resource
		for _, command := range commandList {

			RemoteCommand, err := remote.NewCommand(ctx, command.Name+"-"+node.Name, &remote.CommandArgs{
				Create: pulumi.String(command.CreateCommand),
				Delete: pulumi.String(command.DeleteCommand),
				Update: pulumi.String(command.UpdateCommand),
				Connection: &remote.ConnectionArgs{
					Host:       pulumi.String(node.IP),
					User:       pulumi.String(sshuser.Name),
					PrivateKey: pulumi.String(sshuser.PrivateKey),
				},
			}, pulumi.Parent(RemoteCommandList), pulumi.DependsOn(dependsOn))

			if err != nil {
				return RemoteCommandList, err
			}

			dependsOn = []pulumi.Resource{RemoteCommand}

			if command.ExportOutput {
				ctx.Export(command.Name+"-"+node.Name, RemoteCommand.Stdout)
			}
		}
	}

	return RemoteCommandList, nil
}
