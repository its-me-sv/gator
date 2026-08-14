package main

import (
	"log"
	"os"

	"github.com/its-me-sv/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	appState := state{
		cfg: &cfg,
	}

	appCmds := commands{
		cmdHandlers: make(map[string]func(*state, command) error),
	}
	appCmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalln("not enough arguments were provided")
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = appCmds.run(&appState, cmd)
	if err != nil {
		log.Fatalln(err)
	}
}
