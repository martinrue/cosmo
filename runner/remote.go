package runner

// Remote is a remote script runner.
type Remote struct {
	Host string
}

// Run runs a script remotely.
func (r *Remote) Run(script string) error {
	return Exec("ssh", r.Host, "bash -s <<COSMO", script, "\nCOSMO")
}
