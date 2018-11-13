package commands_test

import (
	"bytes"
	"testing"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
)

func TestCommandRunErrors(t *testing.T) {
	tests := []struct {
		Name     string
		Args     []string
		Expected string
		Err      error
	}{
		{"invalid flag", []string{"task-1", "--invalid"}, "Usage: cosmo run <task> [--server=<name>] [-v]\n", commands.ErrFlagParse},
		{"missing task", []string{}, "Usage: cosmo run <task> [--server=<name>] [-v]\n", commands.ErrNoTask},
		{"ambiguous task", []string{"task-1"}, "task 'task-1' is ambiguous, specify server\n", commands.ErrFindTask},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg, err := config.Read("testdata/run/run.toml")
			if err != nil {
				t.Fatalf("expected config read err to be nil, got (%v)", err)
			}

			buffer := &bytes.Buffer{}
			_, err = commands.NewCommandRun(cfg, nil, nil, test.Args, buffer)

			if err == nil {
				t.Fatalf("expected ctor to return err")
			}

			if err != test.Err {
				t.Fatalf("expected ctor to return err (%v), got (%v)", test.Err, err)
			}

			if buffer.String() != test.Expected {
				t.Fatalf("expected ctor to display error (%v), got (%v)", test.Expected, buffer.String())
			}
		})
	}
}
