package commands

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/table"
)

// CommandTasks lists all tasks for a server.
type CommandTasks struct {
	Config     config.Config
	ServerName string
	Writer     io.Writer
}

func (cmd *CommandTasks) addTasks(serverName string, server config.Server, table *table.Table) {
	keys := make([]string, 0)

	for key := range server.Tasks {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	tasks := make([]string, 0)

	for _, taskName := range keys {
		tasks = append(tasks, taskName)
	}

	table.AddRow(serverName, strings.Join(tasks, ", "))
}

// Exec runs the subcommand.
func (cmd *CommandTasks) Exec() error {
	keys := make([]string, 0)

	for key := range cmd.Config.Servers {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	table := &table.Table{}

	for _, serverName := range keys {
		if cmd.ServerName != "" && cmd.ServerName != serverName {
			continue
		}

		server := cmd.Config.Servers[serverName]

		if len(server.Tasks) > 0 {
			cmd.addTasks(serverName, server, table)
		}
	}

	fmt.Fprintln(cmd.Writer, table)

	return nil
}

// NewCommandTasks creates a new 'tasks' subcommand.
func NewCommandTasks(config config.Config, args []string, writer io.Writer) (Command, error) {
	flags := flag.NewFlagSet("tasks", flag.ContinueOnError)
	server := flags.String("server", "", "")

	flags.Usage = func() {
		fmt.Fprintln(writer, "Usage: cosmo tasks [--server=<name>]")
	}

	if err := flags.Parse(args); err != nil {
		return nil, ErrFlagParse
	}

	if *server != "" {
		if _, ok := config.Servers[*server]; !ok {
			fmt.Fprintf(writer, "server '%s' not found, check config\n", *server)
			return nil, ErrFindServer
		}
	}

	return &CommandTasks{
		Config:     config,
		ServerName: *server,
		Writer:     writer,
	}, nil
}
