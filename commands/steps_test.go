package commands_test

import (
	"bytes"
	"testing"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
)

func TestCommandStepsErrors(t *testing.T) {
	tests := []struct {
		Name     string
		Args     []string
		Expected string
	}{
		{"invalid flag", []string{"task", "--invalid"}, "Usage: cosmo steps <task> [--server=<name>]\n"},
		{"missing task", []string{}, "Usage: cosmo steps <task> [--server=<name>]\n"},
		{"ambiguous task", []string{"task-1"}, "task 'task-1' is ambiguous, specify server\n"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg, err := config.Read("testdata/steps/steps.toml")
			if err != nil {
				t.Fatalf("expected config read err to be nil, got (%v)", err)
			}

			buffer := &bytes.Buffer{}
			if _, err := commands.NewCommandSteps(cfg, test.Args, buffer); err == nil {
				t.Fatalf("expected ctor to return err")
			}

			if buffer.String() != test.Expected {
				t.Fatalf("expected ctor to display error (%v), got (%v)", test.Expected, buffer.String())
			}
		})
	}
}
