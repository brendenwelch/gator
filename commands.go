package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brendenwelch/gator/internal/database"
	"github.com/brendenwelch/gator/internal/rss"
	"github.com/google/uuid"
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
		log.Fatalf("missing username for command '%v'\n", cmd.name)
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatalf("user %v not registered\n", cmd.args[0])
	}
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User has been set to '%v'\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatalf("missing username for command '%v'\n", cmd.name)
	}

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		log.Fatalf("user %v already exists\n", cmd.args[0])
	}
	fmt.Printf("%v has been registered\n", cmd.args[0])

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("User has been set to '%v'\n", cmd.args[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		log.Fatalf("failed to reset database: %v\n", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for i := range users {
		if users[i].Name == s.cfg.Current_user_name {
			fmt.Printf("* %v (current)\n", users[i].Name)
		} else {
			fmt.Printf("* %v\n", users[i].Name)
		}
	}
	return nil
}

func handlerAgg(_ *state, _ command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalf("failed to fetch feed: %v\n", err)
	}
	fmt.Printf("%v", *feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		log.Fatalf("not enough arguments. expecting name, url\n")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		log.Fatalf("failed to retrived user from db: %v\n", err)
	}
	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		log.Fatalf("failed to add feed to db: %v\n", err)
	}
	fmt.Printf("%v\n", feed)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.Feeds(context.Background())
	if err != nil {
		log.Fatalf("failed to retrieve feeds from db: %v", err)
	}

	for i := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feeds[i].UserID)
		if err != nil {
			log.Fatalf("failed to retrieve user from db: %v", err)
		}
		fmt.Printf("%v @ %v added by %v\n", feeds[i].Name, feeds[i].Url, user.Name)
	}
	return nil
}
