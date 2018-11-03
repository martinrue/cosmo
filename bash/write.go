package bash

import (
	"bytes"
	"text/template"

	"github.com/martinrue/cosmo/config"
)

// Write builds up a BASH script from a set of commands.
func Write(commands []config.Step, noEcho bool, skipErrors bool, verbose bool) (string, error) {
	funcs := template.FuncMap{
		"echo": echo(noEcho, verbose),
		"run":  run(skipErrors, verbose),
	}

	t, err := template.New("script").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}
	t.Execute(buffer, commands)

	return buffer.String(), nil
}
