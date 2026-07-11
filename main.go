package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ashy558/bootdev-gator/internal/config"
	"github.com/ashy558/bootdev-gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	code := run()
	os.Exit(code)
}

func run() int {
	var s state
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("could not read current configuration: %s\n", err)
		return 1
	}
	s.cfg = &cfg
	fmt.Printf("Read config: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("could not connect to database: %s\n", err)
		return 1
	}
	s.db = database.New(db)

	cmds := commands{registered: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	input := os.Args
	if len(input) < 2 {
		fmt.Println("error: not enough args provided")
		return 1
	}
	inputCmd := command{
		name: input[1],
		args: input[2:],
	}
	if err := cmds.run(&s, inputCmd); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
