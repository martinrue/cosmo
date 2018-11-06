package runner

import "os/exec"

// Local is a local script runner.
type Local struct {
	Exec Executor
}

// Run runs a script locally.
func (l *Local) Run(script string) error {
	cmd := exec.Command("bash", "-c", script)
	return l.Exec(cmd)
}
