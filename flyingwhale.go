package main

import (
	"flag"
	"fmt"
)
import "github.com/74th/flyingwhale/commands"

func main() {

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("whale <command>")
		return
	}
	command := args[0]
	if command == "install" {
		install := commands.Install{}
		install.Execute()
	}
}
