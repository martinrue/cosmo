package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/drawing"
)

const maxCommandLength = 50

// CommandCmds lists all commands for a task.
type CommandCmds struct {
	Task config.Task
}

func (cmd *CommandCmds) addCommands(prefix string, commands []config.Command, table *drawing.Table) {
	trunc := func(text string) string {
		if len(text) > maxCommandLength {
			return fmt.Sprintf("%s...", text[:maxCommandLength])
		}

		return text
	}

	if len(commands) > 0 {
		lines := make([]string, 0)

		for _, command := range commands {
			lines = append(lines, trunc(command.Exec))
		}

		table.AddRow(prefix, strings.Join(lines, "\n"))
	}
}

// Exec runs the subcommand.
func (cmd *CommandCmds) Exec() {
	table := &drawing.Table{}

	cmd.addCommands("local:", cmd.Task.Local, table)
	cmd.addCommands("remote:", cmd.Task.Remote, table)

	fmt.Println(table)
}

// NewCommandCmds creates a new 'cmds' subcommand.
func NewCommandCmds(config config.Config, args []string) Command {
	flags := flag.NewFlagSet("tasks", flag.ExitOnError)
	server := flags.String("server", "", "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo cmds <task> [--server=<id>]")
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

	return &CommandCmds{task}
}
