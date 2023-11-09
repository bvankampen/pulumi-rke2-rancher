package main

import (
	"bytes"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
	"text/template"
)

func ParseTemplate(fileName string, data any) (string, error) {
	f, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	templ := template.Must(template.New("template").Parse(string(f)))
	var buf bytes.Buffer
	err = templ.Execute(&buf, data)
	return buf.String(), err
}

func CopyFile(ctx *pulumi.Context, node Node, localPath string, remotePath string) error {
	_, err := remote.NewCopyFile(ctx, "", &remote.CopyFileArgs{
		LocalPath:  pulumi.String(localPath),
		RemotePath: pulumi.String(remotePath),
		Connection: &remote.ConnectionArgs{
			Host:       pulumi.String(node.IP),
			User:       pulumi.String(sshuser.Name),
			PrivateKey: pulumi.String(sshuser.PrivateKey),
		},
	})

	if err != nil {
		return err
	}

	return nil
}
