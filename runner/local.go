package runner

// Local is a local script runner.
type Local struct{}

// Run runs a script locally.
func (l *Local) Run(script string) error {
	return Exec("bash", "-c", script)
}
