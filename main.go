package main

import (
	"log"
	"os"

	"github.com/cyberfly100/bootdev_gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Failed to read config:", err)
	}

	progState := &state{cfg: &cfg}

	commands := commands{registeredCmds: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("No command provided")
	}
	cmd := command{name: args[0], args: args[1:]}
	err = commands.run(progState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
