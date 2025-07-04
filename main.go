package main

import (
	"database/sql"
	"log"
	"os"

	commandHanlder "github.com/TheHawk24/gator/internal/commands"
	"github.com/TheHawk24/gator/internal/config"
	"github.com/TheHawk24/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	var state commandHanlder.State

	config_file := config.Read()

	// Connect to the database
	dbURL := config_file.Db_url

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	dbQueries := database.New(db)

	state.Config = &config_file
	state.Db = dbQueries

	//Store commands
	commands_storage := commandHanlder.Commands{}
	commands_storage.Commands_all = make(map[string]func(*commandHanlder.State, commandHanlder.Command) error)
	commands_storage.Register("register", commandHanlder.HandlerRegister)
	commands_storage.Register("login", commandHanlder.HandlerLogin)
	commands_storage.Register("reset", commandHanlder.HandlerReset)
	commands_storage.Register("users", commandHanlder.HandlerListUsers)
	commands_storage.Register("agg", commandHanlder.HandlerAgg)
	commands_storage.Register("addfeed", commandHanlder.MiddlewareLoggedIn(commandHanlder.HandlerAddFeed))
	commands_storage.Register("feeds", commandHanlder.HandlerFeeds)
	commands_storage.Register("follow", commandHanlder.MiddlewareLoggedIn(commandHanlder.HandlerFollow))
	commands_storage.Register("following", commandHanlder.MiddlewareLoggedIn(commandHanlder.HandlerFollowing))
	commands_storage.Register("unfollow", commandHanlder.MiddlewareLoggedIn(commandHanlder.HandlerUnfollow))
	commands_storage.Register("browse", commandHanlder.MiddlewareLoggedIn(commandHanlder.HandlerBrowse))

	//Get command line arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Please speicfy command to use")
	}

	var run_command commandHanlder.Command
	run_command.Name = args[1]
	run_command.Args = args[2:]
	//
	err = commands_storage.Run(&state, run_command)
	if err != nil {
		//log.Fatal(err)
		os.Exit(1)
	}

}
