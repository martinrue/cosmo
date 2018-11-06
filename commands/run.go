package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/runner"
	"github.com/martinrue/cosmo/script"
)

// CommandRun runs a task.
type CommandRun struct {
	Task    config.Task
	Host    string
	Verbose bool
}

// Exec runs the subcommand.
func (cmd *CommandRun) Exec() error {
	execSteps := func(steps []config.Step, runner runner.Runner) error {
		if len(steps) == 0 {
			return nil
		}

		bash := &script.Bash{
			Template:   script.BashTemplate,
			NoEcho:     cmd.Task.NoEcho,
			SkipErrors: cmd.Task.SkipErrors,
			Verbose:    cmd.Verbose,
		}

		script, err := bash.Write(steps)
		if err != nil {
			return fmt.Errorf("error: failed to write bash script: %s", err)
		}

		if err := runner.Run(script); err != nil {
			return fmt.Errorf("error: script run failed: %s", err)
		}

		return nil
	}

	if err := execSteps(cmd.Task.Local, &runner.Local{}); err != nil {
		return err
	}

	if err := execSteps(cmd.Task.Remote, &runner.Remote{Host: cmd.Host}); err != nil {
		return err
	}

	return nil
}

// NewCommandRun creates a new 'run' subcommand.
func NewCommandRun(config config.Config, args []string) Command {
	flags := flag.NewFlagSet("run", flag.ExitOnError)
	flagServer := flags.String("server", "", "")
	flagVerbose := flags.Bool("v", false, "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo run <task> [--server=<name>] [-v]")
	}

	if len(args) == 0 {
		flags.Usage()
		os.Exit(1)
	}

	flags.Parse(args[1:])

	task, server, err := config.Servers.FindTask(args[0], *flagServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	return &CommandRun{
		Task:    task,
		Host:    server.String(),
		Verbose: *flagVerbose,
	}
}
