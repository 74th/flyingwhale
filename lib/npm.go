package lib

import (
	"flag"
	"fmt"
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
	output, errstr := pm.AbstractPackageManager.Client.Exec([]string{"npm", "update"})
	fmt.Println(output)
	fmt.Println(errstr)
	if !strings.HasSuffix(strings.TrimSpace(errstr), "ok") {
		panic("an error ocurred!")
	}
}

// GetBinList is
func (pm *Npm) GetBinList() []string {
	output, errstr := pm.AbstractPackageManager.Client.Exec([]string{"ls", "-1", "/usr/local/bin/"})
	if len(strings.TrimSpace(errstr)) > 0 {
		fmt.Println(output)
		fmt.Println(errstr)
		panic("an error ocurred!")
	}
	return strings.Split(strings.TrimSpace(output), "\n")
}

// Install install node package
func (pm *Npm) Install() {
	output, errstr := pm.AbstractPackageManager.Client.Exec([]string{"npm", "install", "-yg", pm.PackageName})
	fmt.Println(output)
	fmt.Println(errstr)
	if !strings.HasSuffix(strings.TrimSpace(errstr), "ok") {
		panic("an error ocurred!")
	}
}

// CreateCommandScript create a command script
func (pm *Npm) CreateCommandScript(command string) {
	pm.AbstractPackageManager.CreateExecuteCommand(pm.GetContainerName(), command, []string{})
}
