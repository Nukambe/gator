package commands

import (
	"context"
	"fmt"
	cfg "github.com/Nukambe/gator/internal/config"
	"github.com/Nukambe/gator/internal/database"
	"github.com/google/uuid"
	"time"
)

type State struct {
	Db     *database.Queries
	Config *cfg.Config
}

type Command struct {
	Name string
	Args []string
}

type commandHandler func(*State, Command) error

type Commands struct {
	cmds map[string]commandHandler
}

func (c *Commands) register(name string, handler commandHandler) {
	c.cmds[name] = handler
}

func (c *Commands) Run(s *State, cmd Command) error {
	if s == nil {
		return fmt.Errorf("state is nil")
	}
	handler, ok := c.cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("command '%s' does not exist", cmd.Name)
	}
	if err := handler(s, cmd); err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}
	return nil
}

func InitCommands() Commands {
	commands := Commands{cmds: map[string]commandHandler{}}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handleUsers)
	return commands
}

// login
func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login expects a username")
	}
	// check if user exists
	user, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to login: %w", err)
	}
	// login
	if err = s.Config.SetUser(user.Name); err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	// show message if the login command was used
	if cmd.Name == "login" {
		fmt.Println("logged in as:", user.Name)
	}
	return nil
}

// register
func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("register expects a name")
	}

	// create user
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}

	// login with new user
	if err = handlerLogin(s, cmd); err != nil {
		return fmt.Errorf("unable to login after register: %w", err)
	}
	fmt.Printf("User created: %v\n", user)
	return nil
}

// reset
func handlerReset(s *State, cmd Command) error {
	if err := s.Db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("unable to delete users: %w", err)
	}
	fmt.Println("deleted all users")
	return nil
}

// list users
func handleUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %w", err)
	}
	fmt.Println("registered users:")
	for _, user := range users {
		fmt.Printf("	* %s", user.Name)
		if s.Config.CurrentUserName == user.Name {
			fmt.Println(" (current)")
		} else {
			fmt.Println()
		}
	}
	return nil
}
