package main

import (
	"context"
	"fmt"

	"github.com/cyberfly100/bootdev_gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	modified_handler := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Failed to get user from config: %w", err)
		}
		if user == (database.User{}) {
			return fmt.Errorf("No user logged in. Please login first.")
		}
		return handler(s, cmd, user)
	}
	return modified_handler
}
