package lib

import (
	"flag"
	"fmt"
	"strings"
)

// Yum : for nodejs
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
	output, errstr := pm.AbstractPackageManager.Client.Exec([]string{"yum", "update", "-y"})
	fmt.Println(output)
	fmt.Println(errstr)
}

// GetBinList is
func (pm *Yum) GetBinList() []string {
	output, _ := pm.AbstractPackageManager.Client.Exec([]string{"ls", "-1", "/usr/bin/"})
	list := strings.Split(strings.TrimSpace(output), "\n")
	return list
}

// Install install node package
func (pm *Yum) Install() {
	output, errstr := pm.AbstractPackageManager.Client.Exec([]string{"yum", "install", "-y", pm.PackageName})
	fmt.Println(output)
	fmt.Println(errstr)
}

// CreateCommandScript create a command script
func (pm *Yum) CreateCommandScript(command string) bool {
	if command != pm.PackageName {
		return false
	}
	pm.AbstractPackageManager.CreateExecuteCommand(pm.GetContainerName(), command, []string{})
	return true
}