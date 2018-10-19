package commands

import (
	"fmt"

	"github.com/martinrue/cosmo/config"
)

// CommandVersion displays the current cosmo version.
type CommandVersion struct {
	Config *config.Config
}

// Exec runs the subcommand.
func (cmd *CommandVersion) Exec() {
	fmt.Println("v0.0.1")
}

// NewCommandVersion creates a new 'version' subcommand.
func NewCommandVersion(config *config.Config) *CommandVersion {
	return &CommandVersion{
		Config: config,
	}
}
