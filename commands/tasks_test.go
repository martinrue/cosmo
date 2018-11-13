package commands_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
)

func TestCommandTasksInvalidFlag(t *testing.T) {
	cfg, err := config.Read("testdata/tasks/tasks.toml")
	if err != nil {
		t.Fatalf("expected config read err to be nil, got (%v)", err)
	}

	buffer := &bytes.Buffer{}
	if _, err := commands.NewCommandTasks(cfg, []string{"--invalid"}, buffer); err != commands.ErrFlagParse {
		t.Fatalf("expected ctor to return ErrFlagParse err, got (%v)", err)
	}

	actual := buffer.String()
	expected := "Usage: cosmo tasks [--server=<name>]\n"

	if actual != expected {
		t.Fatalf("expected ctor to display error (%v), got (%v)", expected, actual)
	}
}

func TestCommandTasksMissingServer(t *testing.T) {
	cfg, err := config.Read("testdata/tasks/tasks.toml")
	if err != nil {
		t.Fatalf("expected config read err to be nil, got (%v)", err)
	}

	expected := "server 'missing' not found, check config"

	_, err = commands.NewCommandTasks(cfg, []string{"--server", "missing"}, ioutil.Discard)
	if err != nil && err.Error() != expected {
		t.Fatalf("expected ctor err (%v), got (%v)", expected, err)
	}
}

func TestCommandTasksOutput(t *testing.T) {
	stringMatchesGoldenFile := func(t *testing.T, str string, filename string) bool {
		golden, err := ioutil.ReadFile(path.Join("testdata", "tasks", filename))
		if err != nil {
			t.Fatal(err)
		}

		return str == string(golden)
	}

	tests := []struct {
		Name       string
		ConfigFile string
		Args       []string
		GoldenFile string
	}{
		{"all servers", "testdata/tasks/tasks.toml", []string{}, "all-servers.golden"},
		{"server 1", "testdata/tasks/tasks.toml", []string{"--server", "server-1"}, "server-1.golden"},
		{"server 2", "testdata/tasks/tasks.toml", []string{"--server", "server-2"}, "server-2.golden"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg, err := config.Read(test.ConfigFile)
			if err != nil {
				t.Fatalf("expected config read err to be nil, got (%v)", err)
			}

			buffer := &bytes.Buffer{}
			task, err := commands.NewCommandTasks(cfg, test.Args, buffer)
			if err != nil {
				t.Fatalf("expected ctor to return nil err, got (%v)", err)
			}

			if err := task.Exec(); err != nil {
				t.Fatalf("expected exec to return nil err, got (%v)", err)
			}

			if !stringMatchesGoldenFile(t, buffer.String(), test.GoldenFile) {
				t.Fatalf("command output does not match golden file")
			}
		})
	}
}
