package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: login <username>")
	}
	username := cmd.args[0]
	if err := s.config.SetUser(username); err != nil {
		return fmt.Errorf("could not set username: %s", err)
	}
	fmt.Println("user has been set successfully!")
	return nil
}
