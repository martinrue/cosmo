package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/drawing"
)

// CommandTasks lists all tasks for a server.
type CommandTasks struct {
	Config config.Config
	Server string
}

// Exec runs the subcommand.
func (cmd *CommandTasks) Exec() {
	if cmd.Server != "" {
		server, ok := cmd.Config.Servers[cmd.Server]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: server '%s' not found, check config\n", cmd.Server)
			return
		}

		for name := range server.Tasks {
			fmt.Println(name)
		}

		return
	}

	table := &drawing.Table{}

	for server, s := range cmd.Config.Servers {
		if len(s.Tasks) == 0 {
			continue
		}

		lines := make([]string, 0)

		for task := range s.Tasks {
			lines = append(lines, task)
		}

		table.AddRow(fmt.Sprintf("%s:", server), strings.Join(lines, "\n"))
	}

	fmt.Println(table)
}

// NewCommandTasks creates a new 'tasks' subcommand.
func NewCommandTasks(config config.Config, args []string) Command {
	flags := flag.NewFlagSet("tasks", flag.ExitOnError)
	server := flags.String("server", "", "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo tasks [--server=<id>]")
	}

	flags.Parse(args)

	return &CommandTasks{config, *server}
}
