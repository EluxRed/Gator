package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if c == nil || c.registeredCommands == nil {
		return fmt.Errorf("commands registry not initialized")
	}
	if s == nil {
		return fmt.Errorf("state struct is not initialized")
	}
	cmd_handler, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	if err := cmd_handler(s, cmd); err != nil {
		return fmt.Errorf("running command %s: %w", cmd.Name, err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
