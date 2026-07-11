package main

import (
	"github.com/ashy558/bootdev-gator/internal/config"
	"github.com/ashy558/bootdev-gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
