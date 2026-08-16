package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/its-me-sv/gator/internal/config"
	"github.com/its-me-sv/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalln(err)
	}
	dbQueriues := database.New(db)

	appState := state{
		cfg: &cfg,
		db:  dbQueriues,
	}

	appCmds := commands{
		cmdHandlers: make(map[string]func(*state, command) error),
	}
	appCmds.register("login", handlerLogin)
	appCmds.register("register", handlerRegister)
	appCmds.register("reset", handlerReset)
	appCmds.register("users", handlerListUsers)
	appCmds.register("agg", handlerAgg)
	appCmds.register("addfeed", handlerAddFeed)
	appCmds.register("feeds", handlerListFeeds)
	appCmds.register("follow", handlerFeedFollow)
	appCmds.register("following", handlerFeedsFollowedByUser)

	args := os.Args
	if len(args) < 2 {
		log.Fatalln("not enough arguments were provided")
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	if err = appCmds.run(&appState, cmd); err != nil {
		log.Fatalln(err)
	}
}
