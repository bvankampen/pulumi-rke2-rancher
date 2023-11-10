package main

type Node struct {
	Name string
	IP   string
}

type AppConfig struct {
	Version           string
	Airgapped         bool
	SourceBasePath    string
	TempFilesBasePath string
	SourceBaseURL     string
	FilesBasePath     string
	KubeConfigPath    string
	ClusterName       string
	CNI               string
	CISProfile        string
	RKE2Token         string
	TLSSANHostname    string
	Registries        any
}

type RKE2Config struct {
	AppConfig   AppConfig
	FirstNode   bool
	FirstNodeIP string
	NodeName    string
}
