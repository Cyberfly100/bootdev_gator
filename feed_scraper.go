package main

import (
	"context"
	"fmt"

	"github.com/cyberfly100/bootdev_gator/internal/database"
)

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get next feed to fetch: %w", err)
	}

	err = scrapeFeed(s.db, feed)
	if err != nil {
		return fmt.Errorf("Failed to scrape feed: %w", err)
	}
	return nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) error {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("Failed to mark feed as fetched: %w", err)
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Failed to fetch feed: %w", err)
	}
	fmt.Printf("%s\n", feedData.Channel.Title)
	for _, item := range feedData.Channel.Item {
		fmt.Printf("  %s\n", item.Title)
	}
	return nil
}
