package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/ashy558/bootdev-gator/internal/database"
)

var (
	ErrInvalidURL = errors.New("must enter a valid URL")
	ErrNotFound   = errors.New("not found")
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	inputName := cmd.args[0]
	inputURL := cmd.args[1]

	_, err := url.Parse(inputURL)
	if err != nil {
		return ErrInvalidURL
	}

	ctx := context.Background()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   inputName,
		Url:    inputURL,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %s", err)
	}

	follows, err := s.db.CreateFeedFollow(
		ctx,
		database.CreateFeedFollowParams{
			ID:     uuid.New(),
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		// s.db.DeleteFeed(ctx, feed.ID)
		return fmt.Errorf("could not create feed follow: %s", err)
	}
	fmt.Println("Successfully crated new feed:")
	printFeed(feed, user.Name)
	fmt.Println()
	fmt.Println("Successfully following new feed:")
	printFeedFollow(database.GetFeedFollowsForUserRow(follows))
	fmt.Println("=====================================")

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("usage: agg <time_between_reqs>")
	}
	parsedDuration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return errors.New("<time_between_reqs> must be a duration")
	}
	fmt.Printf("Collecting feeds every %v\n", parsedDuration)
	ticker := time.NewTicker(parsedDuration)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return fmt.Errorf("could not fetch feed: %s", err)
		}
	}
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.args) > 1 {
		return errors.New("usage: browse [LIMIT]")
	}
	limit := 2
	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return errors.New("limit must be a valid number")
		}
		limit = parsedLimit
	}
	queryParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(ctx, queryParams)
	if err != nil {
		return fmt.Errorf("could not fetch posts for user: %s", err)
	}
	fmt.Printf("Fetched %d posts for user %s:\n", len(posts), user.Name)
	for i, post := range posts {
		fmt.Println()
		fmt.Printf("%d.\n", i+1)
		fmt.Println(stringifyPost(post))
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("usage: feeds")
	}
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not list feeds: %s", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		fmt.Printf("* Name: %s\n", feed.Name)
		fmt.Printf("* URL: %s\n", feed.Url)
		fmt.Printf("* User: %s\n", feed.Username)
	}
	fmt.Println("=====================================")
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("usage: follow URL")
	}
	inputURL := cmd.args[0]
	_, err := url.Parse(inputURL)
	if err != nil {
		return errors.New("must enter a valid URL")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), inputURL)
	if err != nil {
		return fmt.Errorf("could not fetch feed info: %s", err)
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	},
	)
	if err != nil {
		return fmt.Errorf("could not create feed follows entry: %s", err)
	}
	fmt.Println("Feed Follow entry successfully created!")
	printFeedFollow(database.GetFeedFollowsForUserRow(follow))
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("usage: following")
	}
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("could not fetch following feeds for user: %s", err)
	}
	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, follow := range feedFollows {
		printFeedFollow(follow)
	}
	return nil
}

func handlerHelp(cmds commands) func(*state, command) error {
	return func(s *state, c command) error {
		if len(c.args) != 0 {
			return errors.New("usage: help")
		}
		fmt.Println("Gator CLI")
		cmds.printHelp()
		return nil
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
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
	if len(cmd.args) != 1 {
		return errors.New("usage: register <username>")
	}
	username := cmd.args[0]
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:   uuid.New(),
			Name: username,
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
	if len(cmd.args) != 0 {
		return errors.New("usage: reset")
	}
	err := s.db.TruncateUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: could not truncate users table: %s", err)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("usage: unfollow URL")
	}
	ctx := context.Background()
	rawURL := cmd.args[0]
	_, err := url.Parse(rawURL)
	if err != nil {
		return ErrInvalidURL
	}
	feed, err := s.db.GetFeedByURL(ctx, rawURL)
	if err != nil {
		return ErrNotFound
	}
	_, err = s.db.DeleteFeedFollow(ctx, feed.ID)
	if err != nil {
		return ErrNotFound
	}
	fmt.Println("Feed unfollowed successfully!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("usage: users")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: could not get users: %s", err)
	}
	fmt.Println("Registered Users:")
	printUsers(s, users)
	return nil
}

func printFeed(feed database.Feed, username string) {
	fmt.Printf("* ID: %s\n", feed.ID)
	fmt.Printf("* Created: %v\n", feed.CreatedAt)
	fmt.Printf("* Updated: %v\n", feed.UpdatedAt)
	fmt.Printf("* Name: %s\n", feed.Name)
	fmt.Printf("* URL: %s\n", feed.Url)
	fmt.Printf("* User: %s\n", username)
	fmt.Println("=====================================")
}

func printFeedFollow(follow database.GetFeedFollowsForUserRow) {
	fmt.Printf("* ID: %s\n", follow.ID)
	fmt.Printf("* Created: %v\n", follow.CreatedAt)
	fmt.Printf("* Updated: %v\n", follow.UpdatedAt)
	fmt.Printf("* Feed ID: %s\n", follow.FeedID)
	fmt.Printf("* Name: %s\n", follow.FeedName)
	fmt.Printf("* User ID: %s\n", follow.UserID)
	fmt.Printf("* User: %s\n", follow.UserName)
	fmt.Println("=====================================")
}

func printUsers(s *state, users []database.User) {
	for _, user := range users {
		if user.Name != s.cfg.CurrentUserName {
			fmt.Printf("* %s\n", user.Name)
		} else {
			fmt.Printf("* %s (current)\n", user.Name)
		}
	}
	fmt.Println("=====================================")
}
