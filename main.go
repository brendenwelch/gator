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
		log.Fatalf("no command specified\n")
	}

	s := state{}
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}
	s.cfg = &cfg
	db, err := sql.Open("postgres", s.cfg.Db_url)
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	s.db = database.New(db)

	cmds := commands{
		callbacks: map[string]func(*state, command) error{},
	}
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("register", handlerRegister)
	cmds.register("login", handlerLogin)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmd := command{}
	cmd.name = os.Args[1]
	if len(os.Args) > 2 {
		cmd.args = os.Args[2:]
	}
	cmds.run(&s, cmd)
}
