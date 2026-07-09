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
	currentConfig, err := config.Read()
	if err != nil {
		fmt.Printf("could not read current configuration: %s", err)
		return 1
	}
	if err = currentConfig.SetUser("ashy558"); err != nil {
		fmt.Printf("could not set user: %s", err)
		return 1
	}
	newConfig, err := config.Read()
	if err != nil {
		fmt.Printf("could not read updated configuration: %s", err)
		return 1
	}
	fmt.Println("Updated configuration:")
	fmt.Println(newConfig.String())
	return 0
}
