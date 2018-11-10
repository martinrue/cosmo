package commands_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
)

func TestCommandServersOutput(t *testing.T) {
	stringMatchesGoldenFile := func(t *testing.T, str string, filename string) bool {
		golden, err := ioutil.ReadFile(path.Join("testdata", "servers", filename))
		if err != nil {
			t.Fatal(err)
		}

		return str == string(golden)
	}

	tests := []struct {
		Name       string
		ConfigFile string
		GoldenFile string
	}{
		{"one server", "testdata/servers/one-server.toml", "one-server.golden"},
		{"many servers", "testdata/servers/many-servers.toml", "many-servers.golden"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg, err := config.Read(test.ConfigFile)
			if err != nil {
				t.Fatalf("expected config read err to be nil, got (%v)", err)
			}

			buffer := &bytes.Buffer{}
			task, err := commands.NewCommandServers(cfg, []string{}, buffer)
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
