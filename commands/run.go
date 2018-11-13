package commands

import (
	"flag"
	"fmt"
	"io"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/runner"
	"github.com/martinrue/cosmo/script"
)

// CommandRun runs a task.
type CommandRun struct {
	LocalRunner  runner.Runner
	RemoteRunner runner.Runner
	Writer       io.Writer
	Task         config.Task
	Verbose      bool
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
			return fmt.Errorf("failed to write bash script: %s", err)
		}

		if err := runner.Run(script, cmd.Writer); err != nil {
			return fmt.Errorf("script run failed: %s", err)
		}

		return nil
	}

	if err := execSteps(cmd.Task.Local, cmd.LocalRunner); err != nil {
		return err
	}

	if err := execSteps(cmd.Task.Remote, cmd.RemoteRunner); err != nil {
		return err
	}

	return nil
}

// NewCommandRun creates a new 'run' subcommand.
func NewCommandRun(config config.Config, local runner.Runner, remote runner.Runner, args []string, writer io.Writer) (Command, error) {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	flagServer := flags.String("server", "", "")
	flagVerbose := flags.Bool("v", false, "")

	flags.Usage = func() {
		fmt.Fprintln(writer, "Usage: cosmo run <task> [--server=<name>] [-v]")
	}

	if len(args) == 0 {
		flags.Usage()
		return nil, ErrNoTask
	}

	if err := flags.Parse(args[1:]); err != nil {
		return nil, ErrFlagParse
	}

	task, server, err := config.Servers.FindTask(args[0], *flagServer)
	if err != nil {
		fmt.Fprintln(writer, err)
		return nil, ErrFindTask
	}

	remoteRunner := remote.(*runner.Remote)
	remoteRunner.Host = server.String()

	return &CommandRun{
		LocalRunner:  local,
		RemoteRunner: remote,
		Writer:       writer,
		Task:         task,
		Verbose:      *flagVerbose,
	}, nil
}
