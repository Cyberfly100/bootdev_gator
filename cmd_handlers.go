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
	fmt.Printf("User created.\n  ID: %s\n  Name: %s\n  CreatedAt: %s\n  UpdatedAt: %s", user.ID, user.Name, user.CreatedAt, user.UpdatedAt)
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments expected")
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Failed to fetch feed: %w", err)
	}
	// fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	// fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	// fmt.Printf("Number of Items: %d\n", len(feed.Channel.Item))
	// for _, item := range feed.Channel.Item {
	// 	fmt.Printf("- %s\n", item.Title)
	// 	fmt.Printf("  Description: %s\n", item.Description)
	// }
	fmt.Printf("Feed: %+v\n", feed)
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Feed name and URL are required")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Failed to get current user: %w", err)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUser.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Failed to create feed: %w", err)
	}
	fmt.Printf("Feed created.\n  ID: %s\n  Name: %s\n  URL: %s\n  CreatedAt: %s\n  UpdatedAt: %s\n  User ID: %s\n", feed.ID, feed.Name, feed.Url, feed.CreatedAt, feed.UpdatedAt, feed.UserID)
	// add feedfollow for current user
	feedfollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), feedfollowParams)
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: %w", err)
	}
	fmt.Printf("Feed follow created.\n  Feed name: %s\n  Current user: %s\n", feedFollowRow.FeedName, feedFollowRow.UserName)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments expected")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds: %w", err)
	}
	for _, feed := range feeds {
		username, err := s.db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Failed to get user from ID: %w", err)
		}
		fmt.Printf("=== Feed %s ===\n  URL: %s\n  User name: %s\n", feed.Name, feed.Url, username)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Feed url is required")
	}
	feedURL := cmd.args[0]

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Failed to get current user: %w", err)
	}

	feed, err := s.db.GetFeedFromURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Failed to get feed from URL: %w", err)
	}

	feedfollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), feedfollowParams)
	if err != nil {
		return fmt.Errorf("Failed to create feed follow: %w", err)
	}
	fmt.Printf("Feed follow created.\n  Feed name: %s\n  Current user: %s\n", feedFollowRow.FeedName, feedFollowRow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments expected")
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Failed to get current user: %w", err)
	}
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("Failed to get feed follows for user: %w", err)
	}
	fmt.Printf("Follows of user %s\n ====================\n", currentUser.Name)
	for _, feedFollow := range feedFollows {
		fmt.Printf("  %s\n", feedFollow.FeedName)
	}
	return nil
}
