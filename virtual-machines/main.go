package main

import (
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := createVirtualMachines(ctx)
		if err != nil {
			_ = fmt.Errorf(err.Error())
			return err
		}
		return nil

	})
}
