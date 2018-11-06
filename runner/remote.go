package runner

// Remote is a remote script runner.
type Remote struct {
	Exec Executor
	Host string
}

// Run runs a script remotely.
func (r *Remote) Run(script string) error {
	return r.Exec("ssh", r.Host, "bash -s <<COSMO", script, "\nCOSMO")
}
