package main

import (
	"blog_aggregator/internal/config"
	"errors"
	"fmt"
	"os"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commandHandler func(s *state, cmd command) error

type commands struct {
	handlers map[string]commandHandler
}

func newCommands() *commands {
	return &commands{
		handlers: make(map[string]commandHandler),
	}
}

func (c *commands) register(name string, handler commandHandler) {
	c.handlers[name] = handler
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: '%s'", cmd.name)
	}
	return handler(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username argument is required for login")
	}
	username := cmd.args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user in config: %w", err)
	}

	// Print success message to standard output
	fmt.Printf("Logged in as user: %s\n", username)
	return nil
}

func main() {
	initialCfg, err := config.Read()
	if err != nil {
		// If config read fails critically (not just file not found), exit.
		fmt.Fprintf(os.Stderr, "Error reading initial config: %v\n", err)
		os.Exit(1)
	}

	appState := state{
		cfg: &initialCfg,
	}

	cmds := newCommands()
	cmds.register("login", handlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: No command provided.")
		fmt.Fprintln(os.Stderr, "Usage: gator <command> [arguments...]")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&appState, cmd)
	if err != nil {
		// Print command execution errors to standard error
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
