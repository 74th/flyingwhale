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

// show list additional commands
func (cmd *Install) showListCommand(list []string) []bool {
	input := make([]bool, len(list))
	for i := range list {
		input[i] = false
	}
	for true {
		for i, command := range list {
			if input[i] {
				fmt.Print("[*]")
			} else {
				fmt.Print("[ ]")
			}
			fmt.Println(i+1, ":", command)
		}
		fmt.Print("select number (exit for 0):")
		var in int
		_, err := fmt.Scan(&in)
		if err != nil {
			panic("cannnot load intenger")
		}
		if in == 0 {
			break
		}
		input[in-1] = !input[in-1]
	}
	return input
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

	// list additional command
	installList := cmd.showListCommand(addedCommands)

	for i, command := range addedCommands {
		if installList[i] {
			succ := pm.CreateCommandScript(command)
			if succ {
				fmt.Println("ready for:" + command)
			}
		}
	}
}
