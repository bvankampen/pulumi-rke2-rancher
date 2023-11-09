package main

type Node struct {
	Name string
	IP   string
}

type AppConfig struct {
	Version        string
	Airgapped      bool
	SourceBasePath string
	SourceBaseURL  string
	FilesBasePath  string
}
