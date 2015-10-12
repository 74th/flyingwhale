package lib

import "os"

// IPackageManager is the package manager interface
type IPackageManager interface {
	InitializeInstall()

	GetImageName() string

	GetContainerName() string

	CreateContainer()

	UpdatePackageManager()

	GetBinList() []string

	Install()

	CreateCommandScript(commandName string) bool
}

// AbstractPackageManager is the abstract package manager
type AbstractPackageManager struct {
	Client *DockerClient
}

// CreatePackageManager :
func CreatePackageManager(name string, client *DockerClient) IPackageManager {
	if name == "npm" {
		return CreateNpm(client)
	} else if name == "yum" {
		return CreateYum(client)
	} else if name == "apt-get" {
		return CreateAptGet(client)
	}
	panic("cannot use " + name)
}

// CreateExecuteCommand :
func (*AbstractPackageManager) CreateExecuteCommand(imageName string, commandName string, entryPoint []string) {

	// TODO for root
	// TODO for Windows

	file, err := os.OpenFile("/usr/local/bin/"+commandName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("#!/bin/sh\n")
	file.WriteString("# This script was created by flying docker 0.1 \n")
	//file.WriteString("eval $(whale env)")
	file.WriteString("docker run -it --rm -v `pwd`:/src --workdir=/src --entrypoint=" + commandName + " " + imageName + " $*\n")
}
