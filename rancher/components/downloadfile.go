package components

import (
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"io"
	"net/http"
	"os"
)

func (d *DownloadFileArgs) Exists() bool {
	if _, err := os.Stat(d.LocalFilePath()); err == nil {
		return true
	} else {
		return false
	}
}

func (d *DownloadFileArgs) LocalFilePath() string {
	return d.LocalPath + "/" + d.Version + "/" + d.Name
}

func (d *DownloadFileArgs) GetURL() string {
	return d.BaseURL + "/" + d.Version + "/" + d.Name
}

func (d *DownloadFile) Update() {
	fmt.Print("update")
}

func NewDownloadFile(ctx *pulumi.Context, name string, args DownloadFileArgs, opts ...pulumi.ResourceOption) (*DownloadFile, error) {
	DownloadFile := &DownloadFile{}

	err := ctx.RegisterComponentResource("rancher:rke2:DownloadFile", name, DownloadFile, opts...)
	if err != nil {
		return nil, err
	}

	if ctx.DryRun() {
		return DownloadFile, nil
	}

	err = os.MkdirAll(args.LocalPath+"/"+args.Version, 0750)
	if err != nil {
		return nil, err
	}

	if args.Exists() {
		ctx.Log.Warn("file: "+args.Version+"/"+args.Name+" exists", nil)
		return DownloadFile, nil
	}

	response, err := http.Get(args.GetURL())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	file, err := os.Create(args.LocalFilePath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)

	return DownloadFile, err
}

func NewDownloadRKE2Files(ctx *pulumi.Context, name string, args []DownloadFileArgs, opts ...pulumi.ResourceOption) (*DownloadRKE2Files, error) {
	DownloadRKE2Files := &DownloadRKE2Files{}
	err := ctx.RegisterComponentResource("rancher:rke2:DownloadRKE2Files", name, DownloadRKE2Files, opts...)
	if err != nil {
		return nil, err
	}
	for _, arg := range args {
		_, err := NewDownloadFile(ctx, "download-"+arg.Name, arg, pulumi.Parent(DownloadRKE2Files))
		if err != nil {
			return nil, err
		}
	}
	return DownloadRKE2Files, nil
}
