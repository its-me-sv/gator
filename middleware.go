package main

import (
	"context"
	"errors"

	"github.com/its-me-sv/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, currUser database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return errors.New("user not found, please login first")
		}
		return handler(s, cmd, currUser)
	}
}
