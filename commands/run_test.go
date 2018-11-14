package commands_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/runner"
	"github.com/martinrue/cosmo/script"
)

type CaptureRunner struct {
	Script string
}

func (r *CaptureRunner) Run(script string, writer io.Writer) error {
	r.Script = script
	return nil
}

type ErrorRunner struct {
	Error error
}

func (r *ErrorRunner) Run(script string, writer io.Writer) error {
	return r.Error
}

type ErrorScriptWriter struct {
	Error error
}

func (w *ErrorScriptWriter) Write(steps []config.Step) (string, error) {
	return "", w.Error
}

func readConfig(t *testing.T) config.Config {
	cfg, err := config.Read("testdata/run/run.toml")
	if err != nil {
		t.Fatalf("expected config read err to be nil, got (%v)", err)
	}

	return cfg
}

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
			_, err = commands.NewCommandRun(cfg, nil, nil, nil, test.Args, buffer)

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

func TestCommandRun(t *testing.T) {
	stringMatchesGoldenFile := func(t *testing.T, str string, filename string) bool {
		golden, err := ioutil.ReadFile(path.Join("testdata", "run", filename))
		if err != nil {
			t.Fatal(err)
		}

		return strings.TrimSpace(str) == strings.TrimSpace(string(golden))
	}

	cfg := readConfig(t)

	local := &CaptureRunner{}
	remote := &CaptureRunner{}
	writer := &script.Bash{}

	cmd, err := commands.NewCommandRun(cfg, local, remote, writer, []string{"task-1", "--server", "server-1"}, ioutil.Discard)
	if err != nil {
		t.Fatalf("expected ctor to return err")
	}

	if writer.Template != script.BashTemplate {
		t.Fatalf("expected ctor to set script writer template")
	}

	if writer.SkipErrors != true {
		t.Fatalf("expected ctor to set script writer skip errors flag")
	}

	if err := cmd.Exec(); err != nil {
		t.Fatalf("expected ctor to return nil err, got (%v)", err)
	}

	if !stringMatchesGoldenFile(t, local.Script, "local-script.golden") {
		t.Fatalf("local script did not match golden file, got (%v)", local.Script)
	}

	if !stringMatchesGoldenFile(t, remote.Script, "remote-script.golden") {
		t.Fatalf("remote script did not match golden file, got (%v)", remote.Script)
	}
}

func TestCommandRunRemoteRunner(t *testing.T) {
	cfg := readConfig(t)

	local := &CaptureRunner{}
	remote := &runner.Remote{}
	writer := &script.Bash{}

	_, err := commands.NewCommandRun(cfg, local, remote, writer, []string{"task-1", "--server", "server-1"}, ioutil.Discard)
	if err != nil {
		t.Fatalf("expected ctor to return err")
	}

	if remote.Host != "user-1@host-1" {
		t.Fatalf("expected remote runner host to be 'user-1@host-1', got (%v)", remote.Host)
	}
}

func TestCommandRunLocalRunnerError(t *testing.T) {
	cfg := readConfig(t)

	local := &CaptureRunner{}
	remote := &ErrorRunner{Error: errors.New("remote runner error")}
	writer := &script.Bash{}

	cmd, err := commands.NewCommandRun(cfg, local, remote, writer, []string{"task-1", "--server", "server-1"}, ioutil.Discard)
	if err != nil {
		t.Fatalf("expected ctor to return err")
	}

	expected := fmt.Sprintf("script run failed: %s", remote.Error)

	if err := cmd.Exec(); err.Error() != expected {
		t.Fatalf("expected exec to return err (%v), got (%v)", expected, err)
	}
}

func TestCommandRunRemoteRunnerError(t *testing.T) {
	cfg := readConfig(t)

	local := &ErrorRunner{Error: errors.New("local runner error")}
	remote := &CaptureRunner{}
	writer := &script.Bash{}

	cmd, err := commands.NewCommandRun(cfg, local, remote, writer, []string{"task-1", "--server", "server-1"}, ioutil.Discard)
	if err != nil {
		t.Fatalf("expected ctor to return err")
	}

	expected := fmt.Sprintf("script run failed: %s", local.Error)

	if err := cmd.Exec(); err.Error() != expected {
		t.Fatalf("expected exec to return err (%v), got (%v)", expected, err)
	}
}

func TestCommandRunScriptWriterError(t *testing.T) {
	cfg := readConfig(t)

	local := &CaptureRunner{}
	remote := &CaptureRunner{}
	writer := &ErrorScriptWriter{Error: errors.New("script writer error")}

	cmd, err := commands.NewCommandRun(cfg, local, remote, writer, []string{"task-1", "--server", "server-1"}, ioutil.Discard)
	if err != nil {
		t.Fatalf("expected ctor to return err")
	}

	expected := fmt.Sprintf("failed to write bash script: %s", writer.Error)

	if err := cmd.Exec(); err.Error() != expected {
		t.Fatalf("expected exec to return err (%v), got (%v)", expected, err)
	}
}

func TestCommandRunNoSteps(t *testing.T) {
	cfg := readConfig(t)

	cfg.Servers["server-2"].Tasks["task-1"] = config.Task{
		Local: []config.Step{},
	}

	local := &CaptureRunner{}
	remote := &CaptureRunner{}
	writer := &script.Bash{}

	cmd, err := commands.NewCommandRun(cfg, local, remote, writer, []string{"task-1", "--server", "server-2"}, ioutil.Discard)
	if err != nil {
		t.Fatalf("expected ctor to return err")
	}

	if err := cmd.Exec(); err != nil {
		t.Fatalf("expected exec to return nil err, got (%v)", err)
	}

	if local.Script != "" {
		t.Fatalf("expected local script to be empty, got (%v)", local.Script)
	}

	if remote.Script != "" {
		t.Fatalf("expected remote script to be empty, got (%v)", remote.Script)
	}
}
