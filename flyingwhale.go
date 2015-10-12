package main

import (
	"flag"
	"fmt"
)
import "github.com/74th/flyingwhale/commands"

func main() {

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("whale <package manager> <command>")
		return
	}
	command := args[1]
	if command == "install" {
		install := commands.Install{}
		install.Execute()
	}
}
