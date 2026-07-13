package main

import (
	"context"
	"fmt"

	"github.com/ashy558/bootdev-gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		username := s.cfg.CurrentUserName
		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("could not fetch current user info: %s", err)
		}
		return handler(s, c, user)
	}
}
