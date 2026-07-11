package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ashy558/bootdev-gator/internal/database"
)

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
