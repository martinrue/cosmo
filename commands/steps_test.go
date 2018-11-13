package commands_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
)

func TestCommandStepsErrors(t *testing.T) {
	tests := []struct {
		Name     string
		Args     []string
		Expected string
		Err      error
	}{
		{"invalid flag", []string{"task", "--invalid"}, "Usage: cosmo steps <task> [--server=<name>]\n", commands.ErrFlagParse},
		{"missing task", []string{}, "Usage: cosmo steps <task> [--server=<name>]\n", commands.ErrNoTask},
		{"ambiguous task", []string{"task-1"}, "task 'task-1' is ambiguous, specify server\n", commands.ErrFindTask},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg, err := config.Read("testdata/steps/steps.toml")
			if err != nil {
				t.Fatalf("expected config read err to be nil, got (%v)", err)
			}

			buffer := &bytes.Buffer{}
			_, err = commands.NewCommandSteps(cfg, test.Args, buffer)

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

func TestCommandStepsOutput(t *testing.T) {
	stringMatchesGoldenFile := func(t *testing.T, str string, filename string) bool {
		golden, err := ioutil.ReadFile(path.Join("testdata", "steps", filename))
		if err != nil {
			t.Fatal(err)
		}

		return str == string(golden)
	}

	tests := []struct {
		Name       string
		Args       []string
		GoldenFile string
	}{
		{"first task", []string{"task-1", "--server=server-1"}, "first-task.golden"},
		{"second task", []string{"task-1", "--server", "server-2"}, "second-task.golden"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg, err := config.Read("testdata/steps/steps.toml")
			if err != nil {
				t.Fatalf("expected config read err to be nil, got (%v)", err)
			}

			buffer := &bytes.Buffer{}
			task, err := commands.NewCommandSteps(cfg, test.Args, buffer)
			if err != nil {
				t.Fatalf("expected ctor to return nil err, got (%v)", err)
			}

			if err := task.Exec(); err != nil {
				t.Fatalf("expected exec to return nil err, got (%v)", err)
			}

			fmt.Println(buffer.String())

			if !stringMatchesGoldenFile(t, buffer.String(), test.GoldenFile) {
				t.Fatalf("command output does not match golden file")
			}
		})
	}
}
