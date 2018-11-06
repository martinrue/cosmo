package script

import (
	"bytes"
	"text/template"

	"github.com/martinrue/cosmo/config"
)

// Bash is a type that can render bash scripts.
type Bash struct {
	Template   string
	NoEcho     bool
	SkipErrors bool
	Verbose    bool
}

// Write writes a Bash script from a set of steps.
func (b *Bash) Write(steps []config.Step) (string, error) {
	funcs := template.FuncMap{
		"echo": echo(b.NoEcho, b.Verbose),
		"run":  run(b.SkipErrors, b.Verbose),
	}

	t, err := template.New("script").Funcs(funcs).Parse(b.Template)
	if err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}
	_ = t.Execute(buffer, steps)

	return buffer.String(), nil
}
