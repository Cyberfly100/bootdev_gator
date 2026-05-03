package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/cyberfly100/bootdev_gator/internal/config"
	"github.com/cyberfly100/bootdev_gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Failed to read config:", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	progState := &state{db: dbQueries, cfg: &cfg}

	commands := commands{registeredCmds: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegisterUser)
	commands.register("reset", handlerReset)

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
