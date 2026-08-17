package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/its-me-sv/gator/internal/database"
)

func handlerListPostsForUser(s *state, cmd command, currUser database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		if parsedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = parsedLimit
		}
	}

	userPostParams := database.PostsForUserParams{
		UserID: currUser.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.PostsForUser(context.Background(), userPostParams)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%+v\n", post)
	}
	return nil
}
