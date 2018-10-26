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

// CommandTasks lists all tasks for a server.
type CommandTasks struct {
	Config *config.Config
	Server string
	All    bool
}

func (cmd *CommandTasks) truncate(text string, length int) string {
	if len(text) > length {
		return fmt.Sprintf("%s...", text[:length])
	}

	return text
}

func (cmd *CommandTasks) listTasks(server string, tasks map[string]config.Task, table *drawing.Table) {
	serverName := func() string {
		if server != "" {
			prefix := fmt.Sprintf("%s:", server)
			server = ""
			return prefix
		}

		return ""
	}

	for name, task := range tasks {
		lines := make([]string, 0)

		for _, step := range task.Steps {
			lines = append(lines, fmt.Sprintf("%s", cmd.truncate(step.Exec, maxCommandLength)))
		}

		taskName := fmt.Sprintf("<%s>", name)
		taskSteps := strings.Join(lines, "\n")

		if server == "-" {
			table.AddRow(taskName, taskSteps)
		} else {
			table.AddRow(serverName(), taskName, taskSteps)
		}
	}
}

// Exec runs the subcommand.
func (cmd *CommandTasks) Exec() {
	table := &drawing.Table{}

	if !cmd.All {
		server, ok := cmd.Config.Servers[cmd.Server]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: server '%s' not found, check config\n", cmd.Server)
			return
		}

		cmd.listTasks("-", server.Tasks, table)

		fmt.Println(table)
		return
	}

	for name, server := range cmd.Config.Servers {
		cmd.listTasks(name, server.Tasks, table)
	}

	fmt.Println(table)
}

// NewCommandTasks creates a new 'tasks' subcommand.
func NewCommandTasks(config *config.Config, args []string) Command {
	flags := flag.NewFlagSet("tasks", flag.ExitOnError)
	server := flags.String("server", "", "")
	all := flags.Bool("all", false, "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo tasks --server=<id> [--all]")
	}

	flags.Parse(args)

	if *server == "" && !*all {
		flags.Usage()
		os.Exit(1)
	}

	return &CommandTasks{config, *server, *all}
}
