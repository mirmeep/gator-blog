package main

import (
	"github.com/mirmeep/gator-blog/internal/config"
	"github.com/mirmeep/gator-blog/internal/database"
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	conf *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("cannot read config file: %v", err)
		os.Exit(1)
	}

	dbURL := cfg.Url
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("cannot open database: %v", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	st := state{
		conf: &cfg,
		db: dbQueries,
	}

	cmds := commands{
		cmds: map[string]func(*state, command) error{
			"login": handlerLogin,
			"register": handlerRegister,
			"reset": handlerReset,
			"users": handlerUsers,
			"agg": handlerAgg,
			"addfeed": handlerAddFeed,
			"feeds": handlerFeeds,
			"follow": handlerFollow,
			"following": handlerFollowing,
		},
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Println("error with running command: ", err)	
		os.Exit(1)
	}
}