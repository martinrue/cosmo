package commands

import (
	"fmt"
	"io"
	"sort"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/table"
)

// CommandServers lists all servers from config.
type CommandServers struct {
	Config config.Config
	Writer io.Writer
}

// Exec runs the subcommand.
func (cmd *CommandServers) Exec() error {
	table := &table.Table{}

	keys := make([]string, 0)

	for key := range cmd.Config.Servers {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, serverName := range keys {
		server := cmd.Config.Servers[serverName]
		table.AddRow(serverName, server.String(), fmt.Sprintf("tasks: %d", len(server.Tasks)))
	}

	fmt.Fprintln(cmd.Writer, table)

	return nil
}

// NewCommandServers creates a new 'servers' subcommand.
func NewCommandServers(config config.Config, args []string, writer io.Writer) (Command, error) {
	return &CommandServers{
		Config: config,
		Writer: writer,
	}, nil
}
