package components

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewRemoteCommand(ctx *pulumi.Context, command RemoteCommandArguments, opts ...pulumi.ResourceOption) (*RemoteCommand, error) {
	RemoteCommand := &RemoteCommand{}

	err := ctx.RegisterComponentResource("rancher:rke2:RemoteCommand", command.Name+"-"+command.Node.Name, RemoteCommand, opts...)
	if err != nil {
		return nil, err
	}

	result, err := remote.NewCommand(ctx, command.Name+"-"+command.Node.Name, &remote.CommandArgs{
		Create: pulumi.String(command.CreateCommand),
		Delete: pulumi.String(command.DeleteCommand),
		Update: pulumi.String(command.UpdateCommand),
		Connection: &remote.ConnectionArgs{
			Host:       pulumi.String(command.Node.IP),
			User:       pulumi.String(command.SshUser),
			PrivateKey: pulumi.String(command.SshPrivateKey),
		},
	})

	if err != nil {
		return RemoteCommand, err
	}

	if command.ExportOutput {
		ctx.Export(command.Name+"-"+command.Node.Name, result.Stdout)
	}

	return RemoteCommand, nil
}

func NewRemoteCommandList(ctx *pulumi.Context, name string, nodes []Node, commandList []RemoteCommandListItem, sshuser SshUser, opts ...pulumi.ResourceOption) (*RemoteCommandList, error) {
	RemoteCommandList := &RemoteCommandList{}

	err := ctx.RegisterComponentResource("rancher:rke2:RemoteCommandList", name, RemoteCommandList, opts...)
	if err != nil {
		return nil, err
	}

	for _, command := range commandList {
		for _, node := range nodes {
			_, err := NewRemoteCommand(ctx, RemoteCommandArguments{
				Name:          command.Name,
				Node:          node,
				CreateCommand: command.CreateCommand,
				UpdateCommand: command.UpdateCommand,
				DeleteCommand: command.DeleteCommand,
				ExportOutput:  command.ExportOutput,
				SshUser:       sshuser.Name,
				SshPrivateKey: sshuser.PrivateKey,
			},
				pulumi.Parent(RemoteCommandList))
			if err != nil {
				return nil, err
			}
		}
	}

	return RemoteCommandList, nil
}
