package main

import (
	"fmt"
	"github.com/Nukambe/gator/internal/commands"
	cfg "github.com/Nukambe/gator/internal/config"
	"os"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		fmt.Println("unable to read config file:", err)
		os.Exit(1)
	}
	mainState := commands.State{Config: &config}
	cmds := commands.InitCommands()
	args := os.Args
	if len(args) < 2 {
		fmt.Println("too few arguments")
		os.Exit(1)
	}

	var gatorArgs []string
	if len(args) < 3 {
		gatorArgs = nil
	} else {
		gatorArgs = args[2:]
	}

	err = cmds.Run(&mainState, commands.Command{Name: args[1], Args: gatorArgs})
	if err != nil {
		fmt.Println("unable to run command:", err)
		os.Exit(1)
	}
}
