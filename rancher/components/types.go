package components

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type Node struct {
	Name string
	IP   string
}

type SshUser struct {
	Name           string
	PrivateKeyFile string
	PrivateKey     string
}

type RemoteCommands struct {
	pulumi.ResourceState
}

type RemoteCommandsArguments struct {
	Name          string
	CreateCommand string
	DeleteCommand string
	UpdateCommand string
	ExportOutput  bool
}

type DownloadFile struct {
	pulumi.ResourceState
}

type DownloadFiles struct {
	pulumi.ResourceState
}

type DownloadFileArguments struct {
	Name      string
	BaseURL   string
	Version   string
	LocalPath string
}

type UploadFiles struct {
	pulumi.ResourceState
}

type UploadFilesArguments struct {
	LocalPath         string
	Version           string
	Name              string
	RemotePath        string
	UseSudo           bool
	PostUploadCommand string
	TemplateFile      string
	TemplateData      any
	TemplateTempPath  string
}

type RunRKE2Installer struct {
	pulumi.ResourceState
}
