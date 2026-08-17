package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/its-me-sv/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("missing argument <time_between_reqs>")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to parse <time_between_reqs>. error: %v", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqs.String())
	var lastError error

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			lastError = err
			break
		}
	}

	return lastError
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to fetch next feed. error: %v", err)
	}

	if err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		return fmt.Errorf("unable to mark feed as fetched. error: %v", err)
	}

	feedContent, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed content. error: %v", err)
	}

	for _, feedItem := range feedContent.Channel.Item {
		pubDate := sql.NullTime{}
		if parsedPubDate, err := time.Parse(time.RFC1123Z, feedItem.PubDate); err != nil {
			pubDate = sql.NullTime{
				Time:  parsedPubDate,
				Valid: true,
			}
		}

		postParams := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     feedItem.Title,
			Url:       feedItem.Link,
			Description: sql.NullString{
				String: feedItem.Description,
				Valid:  true,
			},
			PublishedAt: pubDate,
			FeedID:      nextFeed.ID,
		}
		if _, err = s.db.CreatePost(context.Background(), postParams); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				continue
			}
			fmt.Printf("failed to create post: %v\n", err)
		}
	}

	return nil
}
