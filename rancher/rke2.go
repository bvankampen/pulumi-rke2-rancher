package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"rancher/components"
)

func ConfigureNodes(ctx *pulumi.Context) (*components.RemoteCommandList, error) {
	var list = []components.RemoteCommandListItem{
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

func DownloadRKE2Sources(ctx *pulumi.Context, version string, dependsOn []pulumi.Resource) error {

	localPath := appconfig.SourceBasePath
	baseURL := appconfig.SourceBaseURL

	downloadFiles := []components.DownloadFileArgs{
		{
			Name:      "rke2.linux-amd64.tar.gz",
			BaseURL:   baseURL,
			Version:   version,
			LocalPath: localPath,
		},
		{
			Name:      "rke2-images.linux-amd64.tar.zst",
			BaseURL:   baseURL,
			Version:   version,
			LocalPath: localPath,
		},
		{
			Name:      "sha256sum-amd64.txt",
			BaseURL:   baseURL,
			Version:   version,
			LocalPath: localPath,
		},
	}

	_, err := components.NewDownloadRKE2Files(ctx, "download-rke2-files", downloadFiles, pulumi.DependsOn(dependsOn))

	return err
}

func InstallRKE2(ctx *pulumi.Context) error {
	configureNodes, err := ConfigureNodes(ctx)
	err = DownloadRKE2Sources(ctx, appconfig.Version, []pulumi.Resource{configureNodes})
	if err != nil {
		return err
	}
	return nil
}
