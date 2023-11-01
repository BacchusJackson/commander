package main

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/bep/simplecobra"
)

type RootCommand struct {
	Version            string
	logLevelController *slog.LevelVar
}

func NewRootCommand(logLevelController *slog.LevelVar) *RootCommand {
	return &RootCommand{logLevelController: logLevelController}
}

func (r RootCommand) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{}
}

func (r RootCommand) Init(c *simplecobra.Commandeer) error {
	c.CobraCommand.Version = r.Version
	c.CobraCommand.Short = "A CLI tool for templating shell commands"
	c.CobraCommand.PersistentFlags().BoolP("verbose", "v", false, "verbose logging output")
	c.CobraCommand.InitDefaultVersionFlag()
	return nil
}

func (r RootCommand) Name() string {
	return "commander"
}

func (r RootCommand) PreRun(c *simplecobra.Commandeer, c1 *simplecobra.Commandeer) error {
	if verbose, _ := c.CobraCommand.Flags().GetBool("verbose"); verbose {
		r.logLevelController.Set(slog.LevelDebug)
	}
	return nil
}

func (r RootCommand) Run(ctx context.Context, c *simplecobra.Commandeer, s []string) error {
	fmt.Fprintf(c.CobraCommand.OutOrStderr(), "%s\nversion: %s\n", c.CobraCommand.Short, c.CobraCommand.Version)

	return nil
}
