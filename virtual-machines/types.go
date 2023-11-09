package main

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
