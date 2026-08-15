package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument <username>")
	}

	username := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		log.Fatalf("user with name \"%s\" doesn't exist\n", username)
	}
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("user has been set to \"%s\"\n", username)
	return nil
}
