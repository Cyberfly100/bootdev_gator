package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s.cfg == nil {
		return fmt.Errorf("Config not loaded")
	}
	function, ok := c.registeredCmds[cmd.name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", cmd.name)
	}
	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("Failed to execute command: %w", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCmds[name] = f
}
