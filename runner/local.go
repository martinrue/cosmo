package runner

// Local is a local script runner.
type Local struct {
	Exec Executor
}

// Run runs a script locally.
func (l *Local) Run(script string) error {
	return l.Exec("bash", "-c", script)
}
