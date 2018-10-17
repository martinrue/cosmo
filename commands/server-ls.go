package commands

import (
	"fmt"

	"github.com/martinrue/cosmo/config"
)

// ServerLs is a server-ls command.
type ServerLs struct {
	Config *config.Config
}

// Exec handles running of the sub command.
func (cmd *ServerLs) Exec() {
	for name, server := range cmd.Config.Servers {
		fmt.Printf("[%s] %s\n", name, server.String())
	}
}

// NewServerLs creates a new server-ls command.
func NewServerLs(config *config.Config) *ServerLs {
	return &ServerLs{
		Config: config,
	}
}
