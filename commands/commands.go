package commands

import (
	"github.com/martinrue/cosmo/config"
)

// Command describes a runnable command.
type Command interface {
	Exec()
}

// Ctor describes a command constructor function.
type Ctor func(config.Config, []string) Command
