package lib

import (
	"flag"
	"strings"
)

// Npm : for nodejs
type Npm struct {
	AbstractPackageManager
	PackageName string
}

// CreateNpm : create package manager
func CreateNpm(client *DockerClient) IPackageManager {
	apm := AbstractPackageManager{client}
	return &Npm{
		apm,
		""}
}

// CheckExecutableCommands :
func (pm *Npm) CheckExecutableCommands() {
}

// GetImageName gets official image name
func (pm *Npm) GetImageName() string {
	return "library/node"
}

// GetContainerName is const
func (pm *Npm) GetContainerName() string {
	return "whale-npm-" + pm.PackageName
}

// InitializeInstall checks arguments
func (pm *Npm) InitializeInstall() {
	args := flag.Args()
	if len(args) < 3 {
		panic("Require 3 arguments: whale install npm <PackageName>")
	}
	pm.PackageName = args[2]
}

// CreateContainer create a container for npm
func (pm *Npm) CreateContainer() {
	pm.AbstractPackageManager.Client.CreateContainer(
		"whale-npm-"+pm.PackageName,
		"library/node:latest",
		[]string{},
		[]string{})
}

// UpdatePackageManager call npm update
func (pm *Npm) UpdatePackageManager() {
	pm.AbstractPackageManager.Client.ExecWithShowingStdout([]string{"npm", "update"}, true)
}

// GetBinList is
func (pm *Npm) GetBinList() []string {
	output, _ := pm.AbstractPackageManager.Client.ExecShortCommand([]string{"ls", "-1", "/usr/local/bin/"}, false)
	return strings.Split(strings.TrimSpace(output), "\n")
}

// Install install node package
func (pm *Npm) Install() {
	pm.AbstractPackageManager.Client.ExecWithShowingStdout([]string{"npm", "install", "-yg", pm.PackageName}, true)
}

// CreateCommandScript create a command script
func (pm *Npm) CreateCommandScript(command string) bool {
	pm.AbstractPackageManager.CreateExecuteCommand(pm.GetContainerName(), command, []string{})
	return true
}
