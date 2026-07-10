package main

import (
	"fmt"
	"os"

	"github.com/ashy558/bootdev-gator/internal/config"
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
	s.config = &cfg
	fmt.Printf("Read config: %+v\n", cfg)

	cmds := commands{registered: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)

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
