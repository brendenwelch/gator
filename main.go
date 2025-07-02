package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/brendenwelch/gator/internal/config"
	"github.com/brendenwelch/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
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
	db, err := sql.Open("postgres", s.cfg.Db_url)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	s.db = database.New(db)

	cmds := commands{
		callbacks: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	cmds.register("reset", handlerReset)
	cmd := command{}
	cmd.name = os.Args[1]
	if len(os.Args) > 2 {
		cmd.args = os.Args[2:]
	}
	cmds.run(&s, cmd)
}
