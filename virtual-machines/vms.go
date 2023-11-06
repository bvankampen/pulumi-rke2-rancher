package main

import (
	"bytes"
	"github.com/pulumi/pulumi-libvirt/sdk/go/libvirt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"os"
	"text/template"
)

type User struct {
	Name              string
	SSHAuthorizedKeys []string
	Password          string
}

type UserData struct {
	Users    []User
	Longhorn bool
	Hostname string
}

type VMConfig struct {
	BaseImageName string
	BasePoolName  string
	DiskSize      int
	PoolName      string
	Network       NetworkInfo
}

type VM struct {
	Name    string
	IP      string
	Network NetworkInfo
}

type NetworkInfo struct {
	Interface  string
	Name       string
	Gateway    string
	SubnetMask int
	DNS        []string
}

const GB = 1024 * 1024 * 1024 * 1024

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

func createVirtualMachines(ctx *pulumi.Context) error {

	var vmConfig VMConfig
	var vms []VM
	var userdata UserData

	cfg := config.New(ctx, "virtual-machines")
	cfg.RequireObject("vmconfig", &vmConfig)
	cfg.RequireObject("vms", &vms)
	cfg.RequireObject("users", &userdata.Users)

	provider, err := libvirt.NewProvider(ctx, "provider", &libvirt.ProviderArgs{
		Uri: pulumi.String("qemu:///system"),
	})
	if err != nil {
		return err
	}

	for _, vm := range vms {

		vm.Network = vmConfig.Network

		if vm.Network.Interface == "" {
			vm.Network.Interface = "eth0"
		}

		fileSystem, err := libvirt.NewVolume(ctx, vm.Name+"-disk-1", &libvirt.VolumeArgs{
			Pool:           pulumi.String(vmConfig.PoolName),
			BaseVolumeName: pulumi.String(vmConfig.BaseImageName),
			BaseVolumePool: pulumi.String(vmConfig.BasePoolName),
			Size:           pulumi.Int(GB * vmConfig.DiskSize),
		}, pulumi.Provider(provider),
		)

		if err != nil {
			return err
		}

		userdata.Hostname = vm.Name

		cloudInitUserData, err := ParseTemplate("./templates/cloud_init_user_data.yaml.gotmpl", userdata)
		cloudInitMetaData, err := ParseTemplate("./templates/cloud_init_meta_data.yaml.gotmpl", vm)

		cloudInit, err := libvirt.NewCloudInitDisk(ctx, vm.Name+"cloud-init", &libvirt.CloudInitDiskArgs{
			MetaData:      pulumi.String(string(cloudInitMetaData)),
			NetworkConfig: pulumi.String(string(cloudInitMetaData)),
			Pool:          pulumi.String("default"),
			UserData:      pulumi.String(string(cloudInitUserData)),
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		_, err = libvirt.NewDomain(ctx, vm.Name, &libvirt.DomainArgs{
			Memory:    pulumi.Int(1024),
			Vcpu:      pulumi.Int(2),
			Cloudinit: cloudInit.ID(),
			Disks: libvirt.DomainDiskArray{
				libvirt.DomainDiskArgs{
					VolumeId: fileSystem.ID(),
				},
			},
			NetworkInterfaces: libvirt.DomainNetworkInterfaceArray{
				libvirt.DomainNetworkInterfaceArgs{
					NetworkName:  pulumi.String(vmConfig.Network.Name),
					WaitForLease: pulumi.Bool(false),
				},
			},
			Consoles: libvirt.DomainConsoleArray{
				libvirt.DomainConsoleArgs{
					Type:       pulumi.String("pty"),
					TargetPort: pulumi.String("0"),
					TargetType: pulumi.String("serial"),
				},
			},
		},
			pulumi.Provider(provider),
			pulumi.ReplaceOnChanges([]string{"*"}),
			pulumi.DeleteBeforeReplace(true),
		)

		if err != nil {
			return err
		}

	}

	return nil
}
