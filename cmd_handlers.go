package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Username is required")
	}
	username := cmd.args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to set user: %w", err)
	}
	fmt.Println("User set to:", username)
	return nil
}
