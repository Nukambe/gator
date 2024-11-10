package commands

import (
	"fmt"
	cfg "github.com/Nukambe/gator/internal/config"
)

type State struct {
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

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login expects a username")
	}
	if err := s.Config.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("unable to set user: %w", err)
	}
	fmt.Println("User has been set:", cmd.Args[0])
	return nil
}

func InitCommands() Commands {
	commands := Commands{cmds: map[string]commandHandler{}}
	commands.register("login", handlerLogin)
	return commands
}
