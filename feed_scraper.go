package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cyberfly100/bootdev_gator/internal/database"
)

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get next feed to fetch: %w", err)
	}
	nullableTime := sql.NullTime{Time: time.Now(), Valid: true}
	markFeedFetchedParams := database.MarkFeedFetchedParams{
		LastFetchedAt: nullableTime,
		ID:            feed.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), markFeedFetchedParams)
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
