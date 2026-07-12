package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"

	"github.com/ashy558/bootdev-gator/internal/database"
)

const (
	testFeedURL = "https://www.wagslane.dev/index.xml"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}
	username := s.cfg.CurrentUserName
	inputName := cmd.args[0]
	inputURL := cmd.args[1]
	_, err := url.Parse(inputURL)
	if err != nil {
		return errors.New("must enter a valid URL")
	}
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("could not fetch current user info: %s", err)
	}
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      inputName,
		Url:       inputURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %s", err)
	}
	fmt.Println("Successfully crated new feed:")
	fmt.Println(feed)
	follows, err := s.db.CreateFeedFollow(
		ctx,
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		// s.db.DeleteFeed(ctx, feed.ID)
		return fmt.Errorf("could not create feed following entry: %s", err)
	}
	fmt.Println("Successfully following new feed:")
	fmt.Println(follows)
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), testFeedURL)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %s", err)
	}
	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not list feeds: %s", err)
	}
	fmt.Println("Listing feeds:")
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: follow <url>")
	}
	username := s.cfg.CurrentUserName
	inputURL := cmd.args[0]
	_, err := url.Parse(inputURL)
	if err != nil {
		return errors.New("must enter a valid URL")
	}
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("could not fetch current user info: %s", err)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), inputURL)
	if err != nil {
		return fmt.Errorf("could not fetch feed info: %s", err)
	}
	feedFollows, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	},
	)
	if err != nil {
		return fmt.Errorf("could not create feed follows entry: %s", err)
	}
	fmt.Println("feed follows entry successfully created!")
	fmt.Printf("Feed Name: %s\n", feedFollows.FeedName)
	fmt.Printf("User Name: %s\n", feedFollows.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	username := s.cfg.CurrentUserName
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("could not fetch following feeds for user: %s", err)
	}
	fmt.Printf("Feeds followed by %s:\n", username)
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: login <username>")
	}
	username := cmd.args[0]
	user, err := s.db.GetUser(
		context.Background(),
		username,
	)
	if err != nil {
		return fmt.Errorf("error: username %s does not exist", username)
	}
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error: could not set username in config file: %s", err)
	}
	fmt.Println("login: user logged in successfully!")
	fmt.Println(user)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("usage: register <username>")
	}
	username := cmd.args[0]
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		},
	)
	if err != nil {
		return fmt.Errorf("error: could not create user: %s", err)
	}
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error: could not set username in config file: %s", err)
	}
	fmt.Println("register: user created successfully!")
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.TruncateUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: could not truncate users table: %s", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: could not get users: %s", err)
	}
	fmt.Println("Registered Users:")
	for _, user := range users {
		if user.Name != s.cfg.CurrentUserName {
			fmt.Printf("* %s\n", user.Name)
		} else {
			fmt.Printf("* %s (current)\n", user.Name)
		}
	}
	return nil
}
