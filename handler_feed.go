package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/its-me-sv/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("not enough arguments")
	}
	if len(cmd.args) == 1 {
		return errors.New("missing argument <url>")
	}

	currUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
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

	fmt.Printf("feed \"%s\" was created\n", cmd.args[0])
	fmt.Printf("%+v\n", dbFeed)

	return nil
}
