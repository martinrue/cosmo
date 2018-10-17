package commands

// Command describes a runnable command.
type Command interface {
	Exec()
}
