package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registered map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registered[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registered[cmd.name]
	if !ok {
		return errors.New("not found")
	}
	if err := handler(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) printHelp() {
	fmt.Println(`Usage: gator COMMAND

Examples:
  gator addfeed HackerNews https://hnrss.org/
  gator agg 5m
  gator browse 10
  gator feeds
  gator follow https://hnrss.org/
  gator following
  gator help
  gator login boots
  gator register boots
  gator reset
  gator unfollow https://hnrss.org/
  gator users

Commands:
  addfeed NAME URL  Add a new feed
  agg INTERVAL      Fetch new posts from feeds
  browse [LIMIT]    Browse posts from followed feeds
  feeds             Print the feeds added
  follow URL        Follow an existing feed
  following         Print the feeds followed
  help              Print this help message
  login USERNAME    Login as a registered user
  register USERNAME Register a new user
  reset             Truncate all user data (development only)
  unfollow URL      Unfollow a followed feed
  users             Print all registered users
	`)
}
