package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	cmdHandlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdHandlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, handlerExist := c.cmdHandlers[cmd.name]
	if !handlerExist {
		return fmt.Errorf("handler missing for command \"%s\"", cmd.name)
	}
	return handler(s, cmd)
}
