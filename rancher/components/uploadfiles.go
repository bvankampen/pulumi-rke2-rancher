package components

import (
	"bytes"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
	"text/template"
)

func (f *UploadFilesArguments) ParseTemplateToFile() error {
	if f.TemplateFile != "" {
		file, err := os.ReadFile(f.TemplateFile)
		if err != nil {
			return err
		}
		templ := template.Must(template.New("template").Parse(string(file)))
		var buf bytes.Buffer
		err = templ.Execute(&buf, f.TemplateData)
		if err != nil {
			return err
		}
		err = os.WriteFile(f.TemplateTempPath+"/"+f.Name, buf.Bytes(), 0644)
		f.LocalPath = f.TemplateTempPath
		return err
	} else {
		return nil
	}
}

func (f *UploadFilesArguments) GetLocalPath() string {
	if f.Version == "" {
		return f.LocalPath + "/" + f.Name
	} else {
		return f.LocalPath + "/" + f.Version + "/" + f.Name
	}
}

func NewUploadFiles(ctx *pulumi.Context, name string, nodes []Node, files []UploadFilesArguments, sshuser SshUser, opts ...pulumi.ResourceOption) (*UploadFiles, error) {

	tempPath := "/tmp/rancher-sources"

	UploadFiles := &UploadFiles{}
	err := ctx.RegisterComponentResource("rancher:rke2:Uploadfiles", name, UploadFiles, opts...)
	if err != nil {
		return UploadFiles, err
	}

	for _, node := range nodes {
		connection := remote.ConnectionArgs{
			Host:       pulumi.String(node.IP),
			User:       pulumi.String(sshuser.Name),
			PrivateKey: pulumi.String(sshuser.PrivateKey),
		}

		makeTempPath, err := remote.NewCommand(ctx, "make-temp-path-"+name+"-"+node.Name, &remote.CommandArgs{
			Create:     pulumi.String("mkdir -p " + tempPath),
			Connection: &connection,
		}, pulumi.Parent(UploadFiles))
		if err != nil {
			return UploadFiles, err
		}

		for _, file := range files {

			err := file.ParseTemplateToFile()
			if err != nil {
				return UploadFiles, err
			}

			if file.UseSudo {

				makeRemotePath, err := remote.NewCommand(ctx, "make-remote-path-"+file.Name+"-"+node.Name, &remote.CommandArgs{
					Create:     pulumi.String("sudo mkdir -p " + file.RemotePath),
					Connection: &connection,
				}, pulumi.Parent(UploadFiles), pulumi.DependsOn([]pulumi.Resource{makeTempPath}))
				if err != nil {
					return UploadFiles, err
				}

				copyFile, err := remote.NewCopyFile(ctx, file.Name+"-"+node.Name, &remote.CopyFileArgs{
					LocalPath:  pulumi.String(file.GetLocalPath()),
					RemotePath: pulumi.String(tempPath + "/" + file.Name),
					Connection: &connection,
				}, pulumi.Parent(UploadFiles), pulumi.DependsOn([]pulumi.Resource{makeRemotePath}))
				if err != nil {
					return UploadFiles, err
				}

				moveFile, err := remote.NewCommand(ctx, "move-"+file.Name+"-"+node.Name, &remote.CommandArgs{
					Connection: &connection,
					Create:     pulumi.String("sudo  mv " + tempPath + "/" + file.Name + " " + file.RemotePath + " && sudo chown root:root " + file.RemotePath + "/" + file.Name),
				}, pulumi.Parent(UploadFiles), pulumi.DependsOn([]pulumi.Resource{copyFile}))
				if err != nil {
					return UploadFiles, err
				}

				if file.PostUploadCommand != "" {
					_, err = remote.NewCommand(ctx, "post-upload-command-"+file.Name+"-"+node.Name, &remote.CommandArgs{
						Connection: &connection,
						Create:     pulumi.String(file.PostUploadCommand),
					}, pulumi.Parent(UploadFiles), pulumi.DependsOn([]pulumi.Resource{moveFile}))
					if err != nil {
						return UploadFiles, err
					}
				}

			} else {
				makeRemotePath, err := remote.NewCommand(ctx, "make-remote-path-"+file.Name+"-"+node.Name, &remote.CommandArgs{
					Create:     pulumi.String("mkdir -p " + file.RemotePath),
					Connection: &connection,
				}, pulumi.Parent(UploadFiles), pulumi.DependsOn([]pulumi.Resource{makeTempPath}))
				if err != nil {
					return UploadFiles, err
				}

				copyFile, err := remote.NewCopyFile(ctx, file.Name+"-"+node.Name, &remote.CopyFileArgs{
					LocalPath:  pulumi.String(file.GetLocalPath()),
					RemotePath: pulumi.String(file.RemotePath + "/" + file.Name),
					Connection: &connection,
				}, pulumi.Parent(UploadFiles), pulumi.DependsOn([]pulumi.Resource{makeRemotePath}))
				if err != nil {
					return UploadFiles, err
				}

				if file.PostUploadCommand != "" {
					_, err = remote.NewCommand(ctx, "post-upload-command-"+file.Name+"-"+node.Name, &remote.CommandArgs{
						Connection: &connection,
						Create:     pulumi.String(file.PostUploadCommand),
					}, pulumi.Parent(UploadFiles), pulumi.DependsOn([]pulumi.Resource{copyFile}))
					if err != nil {
						return UploadFiles, err
					}
				}
			}

		}
	}

	return UploadFiles, nil
}
