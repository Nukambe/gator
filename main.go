package main

import (
	"database/sql"
	"fmt"
	"github.com/Nukambe/gator/internal/commands"
	cfg "github.com/Nukambe/gator/internal/config"
	"github.com/Nukambe/gator/internal/database"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	// read config
	config, err := cfg.Read()
	if err != nil {
		fmt.Println("unable to read config file:", err)
		os.Exit(1)
	}

	// connect to db
	db, dbErr := sql.Open("postgres", config.DbUrl)
	if dbErr != nil {
		fmt.Println("unable to connect to database:", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	// init state
	mainState := commands.State{Db: dbQueries, Config: &config}

	// init commands
	cmds := commands.InitCommands()

	// read user input
	args := os.Args
	if len(args) < 2 {
		fmt.Println("too few arguments")
		os.Exit(1)
	}

	// parse user input
	var gatorArgs []string
	if len(args) < 3 {
		gatorArgs = nil
	} else {
		gatorArgs = args[2:]
	}

	// run command
	err = cmds.Run(&mainState, commands.Command{Name: args[1], Args: gatorArgs})
	if err != nil {
		fmt.Println("unable to run command:", err)
		os.Exit(1)
	}
}
