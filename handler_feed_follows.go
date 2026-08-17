package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/its-me-sv/gator/internal/database"
)

func handlerFeedFollow(s *state, cmd command, currUser database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("missing argument <url>")
	}

	url := cmd.args[0]
	dbFeed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return errors.New("feed not found")
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currUser.ID,
		FeedID:    dbFeed.ID,
	}
	if _, err = s.db.CreateFeedFollow(context.Background(), feedFollow); err != nil {
		return fmt.Errorf("unable to follow feed, error: %v", err)
	}

	fmt.Printf("feed: %s | user: %s \n", dbFeed.Name, currUser.Name)
	return nil
}

func handlerFeedsFollowedByUser(s *state, cmd command, currUser database.User) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), currUser.ID)
	if err != nil {
		return errors.New("unable to get following feeds")
	}

	for _, followedFeed := range followedFeeds {
		fmt.Printf("Feed name: %s | User name: %s\n", followedFeed.FeedName, followedFeed.UserName)
	}

	return nil
}
