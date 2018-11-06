package runner

// Runner is a type capable of running a script.
type Runner interface {
	Run(script string) error
}
