package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
	"rancher/components"
)

func ConfigureNodes(ctx *pulumi.Context) (*components.RemoteCommands, error) {
	var list = []components.RemoteCommandsArguments{
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

	remoteCommands, err := components.NewRemoteCommands(ctx, "configure-nodes", nodes, list, sshuser)
	if err != nil {
		return remoteCommands, err
	}
	return remoteCommands, nil
}

func DownloadRKE2Sources(ctx *pulumi.Context, dependsOn []pulumi.Resource) (*components.DownloadFiles, error) {

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
		{
			Name:      "install.sh",
			BaseURL:   "https://get.rke2.io",
			Version:   appconfig.Version,
			LocalPath: appconfig.SourceBasePath,
		},
	}

	downloadFiles, err := components.NewDownloadFiles(ctx, "download-rke2-files", files, pulumi.DependsOn(dependsOn))

	return downloadFiles, err
}

func HardenRKE2(ctx *pulumi.Context, dependsOn []pulumi.Resource) (*components.RemoteCommands, error) {
	var list = []components.RemoteCommandsArguments{
		{
			Name:          "add-etcd-group",
			CreateCommand: "sudo groupadd -f --system etcd",
		}, {
			Name:          "add-user-etcd",
			CreateCommand: "id -u etcd &>/dev/null || sudo useradd -M --system -g etcd -s /sbin/nologin etcd",
		}, {
			Name:          "create-etcd-directory",
			CreateCommand: "sudo mkdir -p /var/lib/rancher/rke2/server/db/etcd && sudo chown etcd:etcd /var/lib/rancher/rke2/server/db/etcd && sudo chmod 0700 /var/lib/rancher/rke2/server/db/etcd",
		},
	}

	remoteCommands, err := components.NewRemoteCommands(ctx, "harden-nodes-for-rke2", nodes, list, sshuser, pulumi.DependsOn(dependsOn))
	if err != nil {
		return remoteCommands, err
	}
	return remoteCommands, nil
}

func UploadFiles(ctx *pulumi.Context, dependsOn []pulumi.Resource) (*components.UploadFiles, error) {

	files := []components.UploadFilesArguments{
		{
			Name:              "00-suse-rancher.conf",
			LocalPath:         appconfig.FilesBasePath,
			RemotePath:        "/etc/sysctl.d",
			UseSudo:           true,
			PostUploadCommand: "sudo sysctl -f --system",
		}, {
			Name:       "run-install-rke2.sh",
			LocalPath:  appconfig.FilesBasePath,
			RemotePath: "/opt/rke2/install",
			UseSudo:    true,
		}, {
			Name:       "wait-for-rke2.sh",
			LocalPath:  appconfig.FilesBasePath,
			RemotePath: "/opt/rke2/install",
			UseSudo:    true,
		}, {
			Name:              "bashrc",
			LocalPath:         appconfig.FilesBasePath,
			RemotePath:        "/tmp",
			PostUploadCommand: "cat /tmp/bashrc | sudo tee -a /root/.bashrc",
		},
	}

	if appconfig.Airgapped {
		files = append(files, []components.UploadFilesArguments{
			{
				Name:       "sha256sum-amd64.txt",
				LocalPath:  appconfig.SourceBasePath,
				Version:    appconfig.Version,
				RemotePath: "/opt/rke2/install",
				UseSudo:    true,
			}, {
				Name:       "rke2-images.linux-amd64.tar.zst",
				LocalPath:  appconfig.SourceBasePath,
				Version:    appconfig.Version,
				RemotePath: "/opt/rke2/install",
				UseSudo:    true,
			}, {
				Name:       "rke2.linux-amd64.tar.gz",
				LocalPath:  appconfig.SourceBasePath,
				Version:    appconfig.Version,
				RemotePath: "/opt/rke2/install",
				UseSudo:    true,
			}, {
				Name:       "install.sh",
				LocalPath:  appconfig.SourceBasePath,
				Version:    appconfig.Version,
				RemotePath: "/opt/rke2/install",
				UseSudo:    true,
			}, {
				Name:       "run-install-rke2-airgapped.sh",
				LocalPath:  appconfig.FilesBasePath,
				RemotePath: "/opt/rke2/install",
				UseSudo:    true,
			},
		}...)
	}

	if appconfig.CISProfile != "" {
		files = append(files, []components.UploadFilesArguments{
			{
				Name:       "rancher-psact.yaml",
				LocalPath:  appconfig.FilesBasePath,
				RemotePath: "/etc/rancher/rke2",
				UseSudo:    true,
			},
		}...)
	}

	uploadFiles, err := components.NewUploadFiles(ctx, "upload-files", nodes, files, sshuser, pulumi.DependsOn(dependsOn))
	if err != nil {
		return uploadFiles, err
	}
	return uploadFiles, nil
}

func CreateAndUploadRKE2Config(ctx *pulumi.Context, dependsOn []pulumi.Resource) ([]pulumi.Resource, error) {
	rke2config.AppConfig = appconfig
	rke2config.FirstNode = true
	rke2config.FirstNodeIP = nodes[0].IP
	var allUploadFiles []pulumi.Resource

	for _, node := range nodes {

		rke2config.NodeName = node.Name

		tempPath := appconfig.TempFilesBasePath + "/" + node.Name

		err := os.MkdirAll(tempPath, 0750)
		if err != nil {
			return nil, err
		}

		files := []components.UploadFilesArguments{{
			Name:             "config.yaml",
			LocalPath:        tempPath,
			UseSudo:          true,
			TemplateData:     rke2config,
			TemplateFile:     "./templates/config.yaml.gotmpl",
			RemotePath:       "/etc/rancher/rke2/",
			TemplateTempPath: "../tmp/" + node.Name,
		}}

		if appconfig.Airgapped {
			files = append(files, []components.UploadFilesArguments{{
				Name:             "registries.yaml",
				LocalPath:        tempPath,
				UseSudo:          true,
				TemplateData:     rke2config,
				TemplateFile:     "./templates/registries.yaml.gotmpl",
				RemotePath:       "/etc/rancher/rke2/",
				TemplateTempPath: "../tmp/" + node.Name,
			}}...)
		}

		rke2config.FirstNode = false

		uploadFiles, err := components.NewUploadFiles(ctx, "upload-rke2-config-files-"+node.Name, []components.Node{node}, files, sshuser, pulumi.DependsOn(dependsOn))
		allUploadFiles = append(allUploadFiles, uploadFiles)
		if err != nil {
			return allUploadFiles, err
		}
	}

	return allUploadFiles, nil
}

func RunRKE2Installer(ctx *pulumi.Context, dependsOn []pulumi.Resource) (*components.RunRKE2Installer, error) {
	runRKE2Install, err := components.NewRunRKE2Installer(ctx, "run-rke2-installer", nodes, sshuser, appconfig.Airgapped, pulumi.DependsOn(dependsOn))
	return runRKE2Install, err
}

func GetKubeConfig(ctx *pulumi.Context, dependsOn []pulumi.Resource) (*components.GetKubeConfig, error) {
	getKubeConfig, err := components.NewGetKubeConfig(ctx, "get-kubeconfig", nodes[0], appconfig.KubeConfigPath, appconfig.ClusterName+".yaml", sshuser, pulumi.DependsOn(dependsOn))
	return getKubeConfig, err
}

func InstallRKE2(ctx *pulumi.Context) error {
	configureNodes, err := ConfigureNodes(ctx)
	if err != nil {
		return err
	}
	hardenRKE2, err := HardenRKE2(ctx, []pulumi.Resource{configureNodes})
	if err != nil {
		return err
	}
	downloadRKE2Sources, err := DownloadRKE2Sources(ctx, []pulumi.Resource{hardenRKE2})
	if err != nil {
		return err
	}
	uploadFiles, err := UploadFiles(ctx, []pulumi.Resource{downloadRKE2Sources})
	if err != nil {
		return err
	}
	allUploadFiles, err := CreateAndUploadRKE2Config(ctx, []pulumi.Resource{uploadFiles})
	if err != nil {
		return err
	}
	runRKEInstaller, err := RunRKE2Installer(ctx, allUploadFiles)
	if err != nil {
		return err
	}
	_, err = GetKubeConfig(ctx, []pulumi.Resource{runRKEInstaller})
	if err != nil {
		return err
	}
	return nil
}
