package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/table"
)

// CommandTasks lists all tasks for a server.
type CommandTasks struct {
	Config     config.Config
	ServerName string
}

func (cmd *CommandTasks) addTasks(serverName string, server config.Server, table *table.Table) {
	tasks := make([]string, 0)

	for name := range server.Tasks {
		tasks = append(tasks, name)
	}

	table.AddRow(serverName, strings.Join(tasks, ", "))
}

// Exec runs the subcommand.
func (cmd *CommandTasks) Exec() error {
	table := &table.Table{}

	for serverName, server := range cmd.Config.Servers {
		if cmd.ServerName != "" && cmd.ServerName != serverName {
			continue
		}

		if len(server.Tasks) > 0 {
			cmd.addTasks(serverName, server, table)
		}
	}

	fmt.Println(table)

	return nil
}

// NewCommandTasks creates a new 'tasks' subcommand.
func NewCommandTasks(config config.Config, args []string) Command {
	flags := flag.NewFlagSet("tasks", flag.ExitOnError)
	server := flags.String("server", "", "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo tasks [--server=<name>]")
	}

	flags.Parse(args)

	if *server != "" {
		if _, ok := config.Servers[*server]; !ok {
			fmt.Fprintf(os.Stderr, "error: server '%s' not found, check config\n", *server)
			os.Exit(1)
		}
	}

	return &CommandTasks{
		Config:     config,
		ServerName: *server,
	}
}
