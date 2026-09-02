package main

import (
	"fmt"
)

type state struct {
	conf *Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("command login requires 1 argument, username\n")
	}

	username = cmd.args[0]
	err := conf.SetUser(username) {
		return fmt.Errorf("cannot set username: %v\n", err)
	}

	fmt.Println("user %v has been set", username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	err := c.cmds[cmd.name]
	if err != nil {
		return fmt.Errof("cannot run command: %v\n", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
	return nil
}