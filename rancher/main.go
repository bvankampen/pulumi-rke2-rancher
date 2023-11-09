package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"os"
	"rancher/components"
	"strings"
)

var nodes []components.Node
var sshuser components.SshUser
var appconfig AppConfig

func LoadConfig(ctx *pulumi.Context) error {
	cfg := config.New(ctx, "rancher")
	cfg.RequireObject("nodes", &nodes)
	cfg.RequireObject("sshuser", &sshuser)
	cfg.RequireObject("config", &appconfig)

	if strings.HasPrefix(sshuser.PrivateKeyFile, "~/") {
		sshuser.PrivateKeyFile, _ = homedir.Expand(sshuser.PrivateKeyFile)
	}

	f, err := os.ReadFile(sshuser.PrivateKeyFile)

	if err != nil {
		return err
	}

	sshuser.PrivateKey = string(f)

	return nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := LoadConfig(ctx)
		if err != nil {
			_ = fmt.Errorf(err.Error())
			return err
		}
		err = InstallRKE2(ctx)
		if err != nil {
			_ = fmt.Errorf(err.Error())
			return err
		}
		return nil

	})
}
