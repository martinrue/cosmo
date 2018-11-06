package script

import (
	"fmt"
	"strconv"

	"github.com/martinrue/cosmo/config"
)

type helperFunc func(step config.Step) string

func echo(noEcho bool, verbose bool) helperFunc {
	return func(step config.Step) string {
		if !verbose && (noEcho || step.NoEcho) {
			return ""
		}

		return fmt.Sprintf("echo %s", strconv.Quote(step.Exec))
	}
}

func run(skipErrors bool, verbose bool) helperFunc {
	return func(step config.Step) string {
		command := step.Exec

		if step.Exec[0] == '@' {
			command = fmt.Sprintf("%s >/dev/null", step.Exec[1:])
		}

		if step.Exec[0] == '!' {
			command = fmt.Sprintf("%s >/dev/null 2>&1", step.Exec[1:])
		}

		if skipErrors || step.SkipError {
			command = fmt.Sprintf("%s || true", command)
		}

		if verbose {
			if step.Exec[0] == '@' || step.Exec[0] == '!' {
				command = step.Exec[1:]
			} else {
				command = step.Exec
			}

		}

		return command
	}
}
