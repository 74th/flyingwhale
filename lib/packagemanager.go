package lib

// IPackageManager is the package manager interface
type IPackageManager interface {
	InitializeInstall()

	GetImageName() string

	GetContainerName() string

	CreateContainer()

	UpdatePackageManager()

	Install()
}

// AbstractPackageManager is the abstract package manager
type AbstractPackageManager struct {
	Client *DockerClient
}

// CheckExecutableCommands exec ls
func (pm *AbstractPackageManager) CheckExecutableCommands() {
}

// CreatePackageManager :
func CreatePackageManager(name string, client *DockerClient) IPackageManager {
	if name == "npm" {
		return CreateNpm(client)
	}
	panic("cannot use " + name)
}
