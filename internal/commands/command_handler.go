package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TheHawk24/gator/internal/config"
	"github.com/TheHawk24/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Db     *database.Queries
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

	name := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return err
	}

	s.Config.Current_Username = cmd.Args[0]
	s.Config.SetUser()
	fmt.Printf("User %v is logged in\n", user.Name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {

	if len(cmd.Args) == 0 {
		err := fmt.Sprintf("%v expects a single argument, a username", cmd.Name)
		return errors.New(err)
	}

	var register_user database.CreateUserParams
	register_user.ID = uuid.New()
	register_user.CreatedAt = time.Now()
	register_user.UpdatedAt = time.Now()
	register_user.Name = cmd.Args[0]

	user, err := s.Db.CreateUser(context.Background(), register_user)
	if err != nil {
		return err
	}

	s.Config.Current_Username = cmd.Args[0]
	s.Config.SetUser()

	fmt.Printf("User %v successfully created\n", cmd.Args[0])
	fmt.Println("ID: ", user.ID)
	fmt.Println("Name: ", user.Name)
	fmt.Println("Created At: ", user.CreatedAt)
	fmt.Println("Updated At: ", user.UpdatedAt)

	return nil

}
