package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cyberfly100/bootdev_gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Username is required")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Failed to get user from db: %w", err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to set user: %w", err)
	}
	fmt.Println("User set to:", username)
	return nil
}

func handlerRegisterUser(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Username is required")
	}
	username := cmd.args[0]

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	fmt.Printf("User created.ID: %s\nName: %s\nCreatedAt: %s\nUpdatedAt: %s", user.ID, user.Name, user.CreatedAt, user.UpdatedAt)
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to set user: %w", err)
	}
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments expected")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get users: %w", err)
	}
	for _, user := range users {
		fmt.Printf("* %s%s\n", user.Name, func() string {
			if user.Name == s.cfg.CurrentUserName {
				return " (current)"
			}
			return ""
		}())
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments expected")
	}
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset database: %w", err)
	}
	fmt.Println("Database reset")
	return nil
}
