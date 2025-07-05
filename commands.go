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
	//fmt.Printf("command '%v' successfully registered\n", name)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatalf("missing username for command %v\n", cmd.name)
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatalf("failed to get user %v from db: %v\n", cmd.args[0], err)
	}
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("%v now logged in\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatalf("missing username for command %v\n", cmd.name)
	}

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		log.Fatalf("failed to create user %v: %v\n", cmd.args[0], err)
	}
	fmt.Printf("%v has been registered\n", cmd.args[0])

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Printf("%v now logged in\n", cmd.args[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		log.Fatalf("failed to reset database: %v\n", err)
	}
	fmt.Println("database reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	if len(users) > 0 {
		fmt.Println("registered users:")
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
	// TODO: better logging
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalf("failed to fetch feed: %v\n", err)
	}
	fmt.Printf("%v", *feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		log.Fatalf("missing name, url for command %v\n", cmd.name)
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
	feedfollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatalf("failed to create feed follow: %v\n", err)
	}
	fmt.Printf("%v added and followed %v\n", feedfollow.UserName, feedfollow.FeedName)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.Feeds(context.Background())
	if err != nil {
		log.Fatalf("failed to retrieve feeds from db: %v\n", err)
	}

	if len(feeds) > 0 {
		fmt.Println("aggregated feeds:")
	}
	for i := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feeds[i].UserID)
		if err != nil {
			log.Fatalf("failed to retrieve user from db: %v\n", err)
		}
		fmt.Printf("%v @ %v added by %v\n", feeds[i].Name, feeds[i].Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		log.Fatalf("missing url for command %v\n", cmd.name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatalf("failed to retrieve feed from db: %v\n", err)
	}
	feedfollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatalf("failed to create feed follow: %v\n", err)
	}
	fmt.Printf("%v followed %v\n", feedfollow.UserName, feedfollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatalf("failed to retrieve feed follows from db: %v\n", err)
	}
	if len(follows) > 0 {
		fmt.Printf("%v currently following:\n", user.Name)
	}
	for i := range follows {
		feed, err := s.db.GetFeedByID(context.Background(), follows[i].FeedID)
		if err != nil {
			log.Fatalf("failed to retrieve feed from db: %v\n", err)
		}
		fmt.Printf("%v\n", feed.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		log.Fatalf("missing url for command %v\n", cmd.name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		log.Fatalf("failed to retrieve feed from db: %v\n", err)
	}
	if err := s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		log.Fatalf("failed to remove feed follow from db: %v\n", err)
	}
	return nil
}
