package main

import (
	"github.com/mirmeep/gator-blog/internal/config"
	"fmt"
	"os"
)

type state struct {
	conf *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("cannot read config file: %v", err)
		os.Exit(1)
	}

	st := state{
		conf: &cfg,
	}

	cmds := commands{
		cmds: map[string]func(*state, command) error{
			"login": handlerLogin,
		},
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println("cannot read config file:", err)	
		os.Exit(1)
	}
}