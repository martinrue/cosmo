package runner

import (
	"io"
	"os/exec"
)

// Local is a local script runner.
type Local struct {
	Exec Executor
}

// Run runs a script locally.
func (l *Local) Run(script string, writer io.Writer) error {
	cmd := exec.Command("bash", "-c", script)
	return l.Exec(cmd, writer)
}
