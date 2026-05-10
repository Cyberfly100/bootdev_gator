package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cyberfly100/bootdev_gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		PubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Failed to parse published date: %s", err)
		}
		publishedAt := sql.NullTime{
			Time:  PubDate,
			Valid: !PubDate.IsZero(),
		}
		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}
		createPostParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		}
		_, err = db.CreatePost(context.Background(), createPostParams)
		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation
				// Silently ignore duplicate key errors
				continue
			}
			log.Printf("Failed to create post: %s", err)
		}
	}
	return nil
}
