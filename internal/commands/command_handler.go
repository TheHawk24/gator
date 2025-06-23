package commands

import (
	"errors"
	"fmt"

	"github.com/TheHawk24/gator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands_all map[string]func(*State, Command) error
}

// Run a command handler
func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.Commands_all[cmd.Name]
	if !ok {
		return errors.New("Command not found")
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands_all[name] = f
}

func HandlerLogin(s *State, cmd Command) error {

	if len(cmd.Args) == 0 {
		err := fmt.Sprintf("%v expects a single argument, a username", cmd.Name)
		return errors.New(err)
	}

	s.Config.Current_Username = cmd.Args[0]
	s.Config.SetUser()
	fmt.Println("User is logged in\n")
	return nil
}
