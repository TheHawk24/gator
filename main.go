package main

import (
	"log"
	"os"

	commandHanlder "github.com/TheHawk24/gator/internal/commands"
	"github.com/TheHawk24/gator/internal/config"
)

func main() {

	var state commandHanlder.State

	config_file := config.Read()
	state.Config = &config_file

	commands := commandHanlder.Commands{}
	commands.Commands_all = make(map[string]func(*commandHanlder.State, commandHanlder.Command) error)
	commands.Register("login", commandHanlder.HandlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Please specify command to use")
	}

	var command_name_handle commandHanlder.Command
	commnand_name := args[1]
	command_line_args := args[2:]
	command_name_handle.Name = commnand_name
	command_name_handle.Args = command_line_args

	err := commands.Run(&state, command_name_handle)
	if err != nil {
		log.Fatal(err)
	}

}
