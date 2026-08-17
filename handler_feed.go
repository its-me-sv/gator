package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/its-me-sv/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, currUser database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("not enough arguments")
	}
	if len(cmd.args) == 1 {
		return errors.New("missing argument <url>")
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currUser.ID,
	}

	dbFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("unable to add feed, error: %v", err)
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

	fmt.Printf("feed \"%s\" was created and followed\n", cmd.args[0])
	fmt.Printf("%+v\n", dbFeed)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("too many arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("unable to fetch all feeds: %v\n", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s | Url: %s | Created by: %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}
