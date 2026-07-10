package main

import (
	"errors"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registered map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registered[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registered[cmd.name]
	if !ok {
		return errors.New("not found")
	}
	if err := handler(s, cmd); err != nil {
		return err
	}
	return nil
}
