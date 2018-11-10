package runner

import "io"

// Runner is a type capable of running a script.
type Runner interface {
	Run(script string, writer io.Writer) error
}
