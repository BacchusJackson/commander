package main

import (
	"context"
	"example/pkg/cli"
	"fmt"
	"os"

	"log/slog"

	"github.com/bep/simplecobra"
)

type rootCmd struct {
	Version            string
	logLevelController *slog.LevelVar
}

func NewRootCommand(logLevelController *slog.LevelVar) *rootCmd {
	return &rootCmd{logLevelController: logLevelController}
}
func (r *rootCmd) Name() string { return "cmdr" }

func (r *rootCmd) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{&initCmd{}, &execCmd{}}
}

func (r *rootCmd) Init(c *simplecobra.Commandeer) error {
	c.CobraCommand.Version = r.Version
	c.CobraCommand.Short = "A CLI tool for templating shell commands"
	c.CobraCommand.PersistentFlags().BoolP("verbose", "v", false, "verbose logging output")
	c.CobraCommand.InitDefaultVersionFlag()
	return nil
}

func (r *rootCmd) PreRun(cmdr *simplecobra.Commandeer, _ *simplecobra.Commandeer) error {
	if verbose, _ := cmdr.CobraCommand.Flags().GetBool("verbose"); verbose {
		r.logLevelController.Set(slog.LevelDebug)
	}
	return nil
}

func (r *rootCmd) Run(ctx context.Context, cmdr *simplecobra.Commandeer, s []string) error {
	fmt.Fprintf(cmdr.CobraCommand.OutOrStderr(), "Commander\n%s\nversion: %s\n\nuse `-h` flag for help on any command\n%s\n",
		cmdr.CobraCommand.Short, cmdr.CobraCommand.Version, cmdr.CobraCommand.Use)

	return nil
}

type initCmd struct {
	flagEncode string
}

func (c *initCmd) Name() string {
	return "init"
}
func (c *initCmd) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{}
}

func (c *initCmd) Init(cmdr *simplecobra.Commandeer) error {
	cmdr.CobraCommand.Short = "Print an example configuration file to stdout"
	cmdr.CobraCommand.Flags().StringP("encode", "e", "json", "set text encoding for output can be [json,yaml,toml]")
	return nil
}

func (c *initCmd) PreRun(cmdr *simplecobra.Commandeer, _ *simplecobra.Commandeer) error {
	c.flagEncode, _ = cmdr.CobraCommand.Flags().GetString("encode")
	return nil
}

func (c *initCmd) Run(ctx context.Context, cmdr *simplecobra.Commandeer, _ []string) error {
	slog.Debug("run", "cmd", "init", "encode", c.flagEncode)
	w := cli.NewTypedWriter(cmdr.CobraCommand.OutOrStdout(), cli.ParseEncoding(c.flagEncode))
	return cli.EncodeExample(w)
}

type execCmd struct {
	flagEncode   cli.Encoding
	flagTarget   string
	flagDryRun   bool
	configReader *cli.TypedReader
}

func (c *execCmd) Name() string {
	return "exec"
}
func (c *execCmd) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{}
}

func (c *execCmd) Init(cmdr *simplecobra.Commandeer) error {
	cmdr.CobraCommand.Short = "execute a command after applying templating"
	cmdr.CobraCommand.Flags().StringP("encode", "e", "json", "set text encoding for input can be [json,yaml,toml]")
	cmdr.CobraCommand.Flags().StringP("target", "t", "", "the command to run in the commander file")
	cmdr.CobraCommand.Flags().StringP("file", "f", "commander.json", "the commander configuration file with templated commands and values")
	cmdr.CobraCommand.Flags().BoolP("stdin", "s", false, "read commander file from stdin")
	cmdr.CobraCommand.Flags().BoolP("dry-run", "n", false, "print rendered command but don't execute")

	exampleMsg := `
	cmdr exec --file commander.json --target echo
	cmdr exec --file commander.yaml -e yaml --dry-run --target echo
	cat commander.yaml | cmdr exec --stdin -e yaml --target echo
	`

	cmdr.CobraCommand.Example = fmt.Sprintf("%s%s", cmdr.CobraCommand.Example, exampleMsg)

	cmdr.CobraCommand.MarkFlagFilename("file")
	cmdr.CobraCommand.MarkFlagRequired("target")
	cmdr.CobraCommand.MarkFlagsMutuallyExclusive("file", "stdin")
	return nil
}

func (c *execCmd) PreRun(cmdr *simplecobra.Commandeer, _ *simplecobra.Commandeer) error {
	encodeStr, _ := cmdr.CobraCommand.Flags().GetString("encode")
	c.flagEncode = cli.ParseEncoding(encodeStr)
	c.flagTarget, _ = cmdr.CobraCommand.Flags().GetString("target")
	c.flagDryRun, _ = cmdr.CobraCommand.Flags().GetBool("dry-run")

	log := slog.Default().With("encode", c.flagEncode, "target", c.flagTarget)
	if useStdin, _ := cmdr.CobraCommand.Flags().GetBool("stdin"); useStdin {
		log.Debug("set read from stdin")
		c.configReader = cli.NewTypedReader(cmdr.CobraCommand.InOrStdin(), c.flagEncode)
		return nil
	}

	filename, _ := cmdr.CobraCommand.Flags().GetString("file")
	log = log.With("filename", filename)
	log.Debug("open commander file")

	f, err := os.Open(filename)
	if err != nil {
		log.Error("failed to open commander file", "err", err)
		return err
	}

	c.configReader = cli.NewTypedReader(f, c.flagEncode)
	return nil
}

func (c *execCmd) Run(ctx context.Context, cmdr *simplecobra.Commandeer, _ []string) error {
	log := slog.Default().With("cmd", "exec", "encode", c.flagEncode, "target", c.flagTarget, "dry_run", c.flagDryRun)
	log.Debug("run")
	execCmd, err := cli.NewSystemCmd(c.configReader, c.flagTarget)
	if err != nil {
		log.Error("failed to generate system command", "err", err)
		return err
	}
	log.Debug("execute system command", "exec_cmd", execCmd.String())
	if c.flagDryRun {
		cmdr.CobraCommand.Println(execCmd.String())
		return nil
	}
	return execCmd.Run()
}
