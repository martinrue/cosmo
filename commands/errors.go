package commands

import "errors"

var (
	// ErrFlagParse occurs when an invalid flag is passed to a command.
	ErrFlagParse = errors.New("flag parse error")

	// ErrNoTask occurs when no task is specified in the 'steps' command.
	ErrNoTask = errors.New("no task supplied")

	// ErrFindTask occurs when a task can not be found, or is ambiguous.
	ErrFindTask = errors.New("could not find task")

	// ErrFindServer occurs when a server can not be found.
	ErrFindServer = errors.New("could not find server")
)
