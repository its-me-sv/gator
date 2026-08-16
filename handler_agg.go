package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("too many arguments")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("feed: %+v\n", feed)
	return nil
}
