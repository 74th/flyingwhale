package lib

import (
	"flag"
	"strings"
)

// AptGet : for nodejs
type AptGet struct {
	AbstractPackageManager
	PackageName string
}

// CreateAptGet : create package manager
func CreateAptGet(client *DockerClient) IPackageManager {
	apm := AbstractPackageManager{client}
	return &AptGet{
		apm,
		""}
}

// CheckExecutableCommands :
func (pm *AptGet) CheckExecutableCommands() {
}

// GetImageName gets official image name
func (pm *AptGet) GetImageName() string {
	return "library/ubuntu"
}

// GetContainerName is const
func (pm *AptGet) GetContainerName() string {
	return "whale-apt-get-" + pm.PackageName
}

// InitializeInstall checks arguments
func (pm *AptGet) InitializeInstall() {
	args := flag.Args()
	if len(args) < 3 {
		panic("Require 3 arguments: whale install apt-get <PackageName>")
	}
	pm.PackageName = args[2]
}

// CreateContainer create a container for apt-get
func (pm *AptGet) CreateContainer() {
	pm.AbstractPackageManager.Client.CreateContainer(
		"whale-apt-get-"+pm.PackageName,
		"library/ubuntu:latest",
		[]string{},
		[]string{})
}

// UpdatePackageManager call apt-get update
func (pm *AptGet) UpdatePackageManager() {
	pm.AbstractPackageManager.Client.ExecWithShowingStdout([]string{"apt-get", "update", "-y"}, true)
}

// GetBinList is
func (pm *AptGet) GetBinList() []string {
	output, _ := pm.AbstractPackageManager.Client.ExecShortCommand([]string{"ls", "-1", "/usr/bin/"}, false)
	list := strings.Split(strings.TrimSpace(output), "\n")
	return list
}

// Install install node package
func (pm *AptGet) Install() {
	pm.AbstractPackageManager.Client.ExecWithShowingStdout([]string{"apt-get", "install", "-y", pm.PackageName}, false)
}

// CreateCommandScript create a command script
func (pm *AptGet) CreateCommandScript(command string) bool {
	pm.AbstractPackageManager.CreateExecuteCommand(pm.GetContainerName(), command, []string{})
	return true
}
