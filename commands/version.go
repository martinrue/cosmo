package commands

import (
	"fmt"

	"github.com/martinrue/cosmo/config"
)

// Version is a version command.
type Version struct {
	Config *config.Config
}

// Exec handles running of the sub command.
func (cmd *Version) Exec() {
	fmt.Println("v0.0.1")
}

// NewVersion creates a new version command.
func NewVersion(config *config.Config) *Version {
	return &Version{
		Config: config,
	}
}
