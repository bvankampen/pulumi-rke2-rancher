package components

import (
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"io"
	"net/http"
	"os"
)

func (d *DownloadFileArguments) Exists() bool {
	if _, err := os.Stat(d.LocalFilePath()); err == nil {
		return true
	} else {
		return false
	}
}

func (d *DownloadFileArguments) LocalFilePath() string {
	return d.LocalPath + "/" + d.Version + "/" + d.Name
}

func (d *DownloadFileArguments) GetURL() string {
	return d.BaseURL + "/" + d.Version + "/" + d.Name
}

func (d *DownloadFile) Update() {
	fmt.Print("update")
}

func NewDownloadFile(ctx *pulumi.Context, name string, file DownloadFileArguments, opts ...pulumi.ResourceOption) (*DownloadFile, error) {
	DownloadFile := &DownloadFile{}

	err := ctx.RegisterComponentResource("rancher:rke2:DownloadFile", name, DownloadFile, opts...)
	if err != nil {
		return nil, err
	}

	// Do not download file if this is a DryRun (preview)
	if ctx.DryRun() {
		return DownloadFile, nil
	}

	err = os.MkdirAll(file.LocalPath+"/"+file.Version, 0750)
	if err != nil {
		return nil, err
	}

	// If file exists do not download file
	if file.Exists() {
		ctx.Log.Warn("localFile: "+file.Version+"/"+file.Name+" exists", nil)
		return DownloadFile, nil
	}

	response, err := http.Get(file.GetURL())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	localFile, err := os.Create(file.LocalFilePath())
	if err != nil {
		return nil, err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, response.Body)

	return DownloadFile, err
}

func NewDownloadFiles(ctx *pulumi.Context, name string, files []DownloadFileArguments, opts ...pulumi.ResourceOption) (*DownloadFiles, error) {
	DownloadRKE2Files := &DownloadFiles{}
	err := ctx.RegisterComponentResource("rancher:rke2:DownloadFiles", name, DownloadRKE2Files, opts...)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		_, err := NewDownloadFile(ctx, "download-"+file.Name, file, pulumi.Parent(DownloadRKE2Files))
		if err != nil {
			return nil, err
		}
	}
	return DownloadRKE2Files, nil
}
