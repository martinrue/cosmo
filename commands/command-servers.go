package commands

import (
	"fmt"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/drawing"
)

// CommandServers lists all servers from config.
type CommandServers struct {
	Config *config.Config
}

// Exec runs the subcommand.
func (cmd *CommandServers) Exec() {
	table := &drawing.Table{}

	for name, server := range cmd.Config.Servers {
		table.AddRow(fmt.Sprintf("%s:", name), server.String())
	}

	fmt.Println(table)
}

// NewCommandServers creates a new 'servers' subcommand.
func NewCommandServers(config *config.Config, args []string) Command {
	return &CommandServers{config}
}
