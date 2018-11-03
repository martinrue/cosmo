package commands

import (
	"fmt"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/text"
)

// CommandServers lists all servers from config.
type CommandServers struct {
	Config config.Config
}

// Exec runs the subcommand.
func (cmd *CommandServers) Exec() error {
	table := &text.Table{}

	for name, server := range cmd.Config.Servers {
		table.AddRow(name, server.String(), fmt.Sprintf("tasks: %d", len(server.Tasks)))
	}

	table.Print()

	return nil
}

// NewCommandServers creates a new 'servers' subcommand.
func NewCommandServers(config config.Config, args []string) Command {
	return &CommandServers{
		Config: config,
	}
}
