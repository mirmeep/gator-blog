package main

import "fmt"


func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("command login requires 1 argument, username")
	}

	username := cmd.args[0]
	err := s.conf.SetUser(username)
	if err != nil {
		return fmt.Errorf("cannot set username: %v", err)
	}

	fmt.Printf("user %v has been set\n", username)
	return nil
}