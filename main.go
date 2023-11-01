package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/bep/simplecobra"
	"github.com/lmittmann/tint"
)

const CLIVersion = "v0.1.0-pre"
const exitOk = 0
const exitSystemFailure = 1
const exitCommandFailure = 2

func main() {
	os.Exit(Run(os.Stderr, os.Stdout))
}

func Run(logWriter io.Writer, outWriter io.Writer) (exitCode int) {
	logLevelController := new(slog.LevelVar)
	logLevelController.Set(slog.LevelWarn)

	h := tint.NewHandler(logWriter, &tint.Options{
		Level:      logLevelController,
		TimeFormat: time.TimeOnly,
	})
	slog.SetDefault(slog.New(h))

	command := NewRootCommand(logLevelController)
	command.Version = CLIVersion

	x, err := simplecobra.New(command)
	if err != nil {
		slog.Error("failed to initialize command", "err", err)
		return exitSystemFailure
	}

	if _, err := x.Execute(context.Background(), os.Args[1:]); err != nil {
		slog.Error("command execution failure", "err", err)
		return exitCommandFailure
	}

	return exitOk
}
