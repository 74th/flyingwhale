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
		fmt.Println("whale <package manager> install <package>")
		panic("too few arguments")
	}
	cmd.pmName = args[0]
}

func (cmd *Install) compareNewCommand(before []string, after []string) []string {
	result := make([]string, 0, 0)
	for _, a := range after {
		isNew := true
		for _, b := range before {
			if a == b {
				isNew = false
				break
			}
		}
		if isNew {
			result = append(result, a)
		}
	}
	return result
}

// Execute this commands
func (cmd *Install) Execute() {

	cmd.checkArgs()

	docker := lib.CreateDockerClient()
	pm := lib.CreatePackageManager(cmd.pmName, docker)

	// Init
	pm.InitializeInstall()

	// Pull
	fmt.Print("Pulling " + pm.GetImageName() + "...")
	docker.PullImage(pm.GetImageName())
	fmt.Println("Done")

	// TODO check dupricate container name

	pm.CreateContainer()
	defer docker.RemoveContainer()

	before := pm.GetBinList()

	pm.UpdatePackageManager()

	pm.Install()

	after := pm.GetBinList()

	addedCommands := cmd.compareNewCommand(before, after)
	if len(addedCommands) == 0 {
		fmt.Println("cannot found an additional command")
		docker.StopContainer()
		docker.RemoveContainer()
		return
	}

	docker.CommitContainer(pm.GetContainerName())

	for _, command := range addedCommands {
		succ := pm.CreateCommandScript(command)
		if succ {
			fmt.Println("ready for:" + command)
		}
	}
}
