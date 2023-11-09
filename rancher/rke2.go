package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"rancher/components"
)

func ConfigureNodes(ctx *pulumi.Context) (*components.RemoteCommandList, error) {
	var list = []components.RemoteCommandListArguments{
		{
			Name:          "disable-and-stop-firewalld",
			CreateCommand: "if systemctl | grep firewalld.service; then sudo systemctl disable --now firewalld; fi",
		}, {
			Name:          "disable-and-stop-swap",
			CreateCommand: "sudo systemctl disable --now swap.target",
		}, {
			Name:          "disable-swap",
			CreateCommand: "sudo swapoff -a",
		},
	}

	remoteCommandList, err := components.NewRemoteCommandList(ctx, "configure-nodes", nodes, list, sshuser)
	if err != nil {
		return remoteCommandList, err
	}
	return remoteCommandList, nil
}

func DownloadRKE2Sources(ctx *pulumi.Context, dependsOn []pulumi.Resource) (*components.DownloadRKE2Files, error) {

	files := []components.DownloadFileArguments{
		{
			Name:      "rke2.linux-amd64.tar.gz",
			BaseURL:   appconfig.SourceBaseURL,
			Version:   appconfig.Version,
			LocalPath: appconfig.SourceBasePath,
		},
		{
			Name:      "rke2-images.linux-amd64.tar.zst",
			BaseURL:   appconfig.SourceBaseURL,
			Version:   appconfig.Version,
			LocalPath: appconfig.SourceBasePath,
		},
		{
			Name:      "sha256sum-amd64.txt",
			BaseURL:   appconfig.SourceBaseURL,
			Version:   appconfig.Version,
			LocalPath: appconfig.SourceBasePath,
		},
	}

	downloadFiles, err := components.NewDownloadRKE2Files(ctx, "download-rke2-files", files, pulumi.DependsOn(dependsOn))

	return downloadFiles, err
}

func UploadFiles(ctx *pulumi.Context, dependsOn []pulumi.Resource) error {
	return nil
}

func InstallRKE2(ctx *pulumi.Context) error {
	configureNodes, err := ConfigureNodes(ctx)
	if err != nil {
		return err
	}
	downloadRKE2Sources, err := DownloadRKE2Sources(ctx, []pulumi.Resource{configureNodes})
	if err != nil {
		return err
	}
	err = UploadFiles(ctx, []pulumi.Resource{downloadRKE2Sources})
	if err != nil {
		return err
	}
	return nil
}
