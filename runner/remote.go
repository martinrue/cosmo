package runner

import "os/exec"

// Remote is a remote script runner.
type Remote struct {
	Exec Executor
	Host string
}

// Run runs a script remotely.
func (r *Remote) Run(script string) error {
	cmd := exec.Command("ssh", r.Host, "bash -s <<COSMO", script, "\nCOSMO")
	return r.Exec(cmd)
}
