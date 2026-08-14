package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/its-me-sv/gator/internal/config"
)

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

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument <username>")
	}

	username := cmd.args[0]
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("user has been set to \"%s\"\n", username)
	return nil
}

type commands struct {
	cmdHandlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, handlerExist := c.cmdHandlers[cmd.name]
	if !handlerExist {
		return fmt.Errorf("handler missing for command \"%s\"", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdHandlers[name] = f
}
