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

type RemoteCommand struct {
	pulumi.ResourceState
}

type RemoteCommandList struct {
	pulumi.ResourceState
}

type RemoteCommandArguments struct {
	Name          string
	Node          Node
	CreateCommand string
	DeleteCommand string
	UpdateCommand string
	ExportOutput  bool
	SshUser       string
	SshPrivateKey string
}

type RemoteCommandListItem struct {
	Name          string
	CreateCommand string
	DeleteCommand string
	UpdateCommand string
	ExportOutput  bool
}

type DownloadFile struct {
	pulumi.ResourceState
}

type DownloadRKE2Files struct {
	pulumi.ResourceState
}

type DownloadFileArgs struct {
	Name      string
	BaseURL   string
	Version   string
	LocalPath string
}
