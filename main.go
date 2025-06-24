package main

import (
	"log"
	"os"

	"github.com/brendenwelch/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("no command specified")
	}

	s := state{}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	s.cfg = &cfg

	cmds := commands{
		callbacks: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	cmd := command{}
	cmd.name = os.Args[1]
	if len(os.Args) > 2 {
		cmd.args = os.Args[2:]
	}
	cmds.run(&s, cmd)
}
