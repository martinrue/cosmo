package commands

// Command describes a runnable command.
type Command interface {
	Exec() error
}

// // Ctor describes a command constructor function.
// type Ctor func(config.Config, []string) Command
