package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/its-me-sv/gator/internal/database"
)

const (
	timeLayout       = "2006-01-02 15:04:05 MST"
	descriptionLimit = 200
)

func printUser(user database.User) {
	fmt.Println("User:")
	fmt.Printf("  ID:      %s\n", user.ID)
	fmt.Printf("  Name:    %s\n", user.Name)
	fmt.Printf("  Created: %s\n", user.CreatedAt.Format(timeLayout))
	fmt.Printf("  Updated: %s\n", user.UpdatedAt.Format(timeLayout))
}

func printFeed(feed database.Feed) {
	fmt.Println("Feed:")
	fmt.Printf("  ID:           %s\n", feed.ID)
	fmt.Printf("  Name:         %s\n", feed.Name)
	fmt.Printf("  Url:          %s\n", feed.Url)
	fmt.Printf("  Created:      %s\n", feed.CreatedAt.Format(timeLayout))
	fmt.Printf("  Updated:      %s\n", feed.UpdatedAt.Format(timeLayout))
	fmt.Printf("  Last fetched: %s\n", formatNullTime(feed.LastFetchedAt, "never"))
}

func printPost(post database.PostsForUserRow) {
	fmt.Printf("%s\n", post.Title)
	fmt.Printf("  Feed:        %s\n", post.FeedName)
	fmt.Printf("  Url:         %s\n", post.Url)
	fmt.Printf("  Published:   %s\n", formatNullTime(post.PublishedAt, "unknown"))
	if description := formatDescription(post.Description); description != "" {
		fmt.Printf("  Description: %s\n", description)
	}
}

func formatNullTime(nullTime sql.NullTime, fallback string) string {
	if !nullTime.Valid {
		return fallback
	}
	return nullTime.Time.Format(timeLayout)
}

// formatDescription collapses the whitespace an RSS description usually carries
// and caps it, so a long entry can't wreck the layout of a browse listing.
func formatDescription(nullDescription sql.NullString) string {
	if !nullDescription.Valid {
		return ""
	}

	collapsed := strings.Join(strings.Fields(nullDescription.String), " ")
	runes := []rune(collapsed)
	if len(runes) <= descriptionLimit {
		return collapsed
	}
	return string(runes[:descriptionLimit]) + "..."
}
