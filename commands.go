package main

import (
	"fmt"
	"log"
)

type command struct {
	name string
	args []string
}

type commands struct {
	callbacks map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	return c.callbacks[cmd.name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.callbacks[name] = f
	fmt.Printf("command '%v' successfully registered\n", name)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatalf("missing username for command '%v'", cmd.name)
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User has been set to '%v'\n", cmd.args[0])

	return nil
}
