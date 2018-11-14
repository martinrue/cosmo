package script

import "github.com/martinrue/cosmo/config"

// Writer is a type capabale of writing steps to a script.
type Writer interface {
	Write(steps []config.Step) (string, error)
}
