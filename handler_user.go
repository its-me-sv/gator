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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument <name>")
	}

	name := cmd.args[0]
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	dbUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatalf("user with name \"%s\" already exists\n", name)
	}
	if err := s.cfg.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("user \"%s\" was created\n", name)
	fmt.Printf("%+v\n", dbUser)

	return nil
}

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

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("too many arguments")
	}

	if err := s.db.DeleteUsers(context.Background()); err != nil {
		log.Fatalf("unable to delete all users: %v\n", err)
	}

	fmt.Println("reset success")
	return nil
}
