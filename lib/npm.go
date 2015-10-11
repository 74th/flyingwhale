package lib

import "flag"

// Npm : for nodejs
type Npm struct {
	AbstractPackageManager
	PackageName string
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

// UpdatePackageManager : UpdatePackageManager
func (pm *Npm) UpdatePackageManager() {
	pm.AbstractPackageManager.Client.Exec([]string{"npm", "update"})
}

// Install : install
func (pm *Npm) Install() {
}

// CreateNpm : create package manager
func CreateNpm(client *DockerClient) IPackageManager {
	apm := AbstractPackageManager{client}
	return &Npm{
		apm,
		""}
}
