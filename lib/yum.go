package lib

import (
	"flag"
	"strings"
)

// Yum is the package manager of centos
type Yum struct {
	AbstractPackageManager
	PackageName string
}

// CreateYum : create package manager
func CreateYum(client *DockerClient) IPackageManager {
	apm := AbstractPackageManager{client}
	return &Yum{
		apm,
		""}
}

// CheckExecutableCommands :
func (pm *Yum) CheckExecutableCommands() {
}

// GetImageName gets official image name
func (pm *Yum) GetImageName() string {
	return "library/centos"
}

// GetContainerName is const
func (pm *Yum) GetContainerName() string {
	return "whale-yum-" + pm.PackageName
}

// InitializeInstall checks arguments
func (pm *Yum) InitializeInstall() {
	args := flag.Args()
	if len(args) < 3 {
		panic("Require 3 arguments: whale install yum <PackageName>")
	}
	pm.PackageName = args[2]
}

// CreateContainer create a container for yum
func (pm *Yum) CreateContainer() {
	pm.AbstractPackageManager.Client.CreateContainer(
		"whale-yum-"+pm.PackageName,
		"library/centos:latest",
		[]string{},
		[]string{})
}

// UpdatePackageManager call yum update
func (pm *Yum) UpdatePackageManager() {
	pm.AbstractPackageManager.Client.ExecWithShowingStdout([]string{"yum", "update", "-y"}, false)
}

// GetBinList is
func (pm *Yum) GetBinList() []string {
	output, _ := pm.AbstractPackageManager.Client.ExecShortCommand([]string{"ls", "-1", "/usr/bin/"}, false)
	list := strings.Split(strings.TrimSpace(output), "\n")
	return list
}

// Install install node package
func (pm *Yum) Install() {
	pm.AbstractPackageManager.Client.ExecWithShowingStdout([]string{"yum", "install", "-y", pm.PackageName}, false)
}

// CreateCommandScript create a command script
func (pm *Yum) CreateCommandScript(command string) bool {
	pm.AbstractPackageManager.CreateExecuteCommand(pm.GetContainerName(), command, []string{})
	return true
}
