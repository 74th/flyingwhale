package commands

import (
	"flag"
	"fmt"

	"github.com/74th/flyingwhale/lib"
)

// Install module
type Install struct {
	pmName string
}

// check arguments
func (cmd *Install) checkArgs() {
	args := flag.Args()
	if len(args) < 3 {
		panic("too few arguments")
	}
	cmd.pmName = args[1]
}

// Execute this commands
func (cmd *Install) Execute() {

	cmd.checkArgs()

	docker := lib.CreateDockerClient()
	pm := lib.CreatePackageManager(cmd.pmName, docker)

	// Init
	pm.InitializeInstall()

	// Pull
	// fmt.Println("Pulling " + pm.GetImageName())
	// docker.PullImage(pm.GetImageName())

	// TODO check dupricate container name

	// CreateContainer
	fmt.Println("Creating container ")
	pm.CreateContainer()

	fmt.Println("Updating package manager")
	pm.UpdatePackageManager()
}
