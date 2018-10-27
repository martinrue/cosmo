package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/martinrue/cosmo/config"
)

// CommandRun runs a task.
type CommandRun struct {
	Task config.Task
}

// Exec runs the subcommand.
func (cmd *CommandRun) Exec() {
	// TODO: run cmd.Task's Local and Remote commands
}

// NewCommandRun creates a new 'run' subcommand.
func NewCommandRun(config config.Config, args []string) Command {
	flags := flag.NewFlagSet("run", flag.ExitOnError)
	server := flags.String("server", "", "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo run <task> --server=<id>")
	}

	if len(args) == 0 {
		flags.Usage()
		os.Exit(1)
	}

	flags.Parse(args[1:])

	task, err := config.Servers.FindTask(args[0], *server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	return &CommandRun{task}
}
