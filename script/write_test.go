package script_test

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/script"
)

var steps = []config.Step{
	config.Step{Exec: "echo one"},
	config.Step{Exec: "@echo two"},
	config.Step{Exec: "!echo three"},
}

func TestWriteInvalidTemplate(t *testing.T) {
	bash := &script.Bash{
		Template: "{{.}",
	}

	expected := `template: script:1: unexpected "}" in operand`

	_, err := bash.Write(steps)
	if err != nil && err.Error() != expected {
		t.Fatalf("expected write err to be (%v), got (%v)", expected, err)
	}
}

func TestWriteValidTemplate(t *testing.T) {
	scriptMatchesGoldenFile := func(t *testing.T, script string, filename string) bool {
		golden, err := ioutil.ReadFile(path.Join("testdata", filename))
		if err != nil {
			t.Fatal(err)
		}

		return strings.TrimSpace(script) == strings.TrimSpace(string(golden))
	}

	tests := []struct {
		Name       string
		Template   string
		NoEcho     bool
		SkipErrors bool
		Verbose    bool
		Steps      []config.Step
		GoldenFile string
	}{
		{"no steps", script.BashTemplate, false, false, false, []config.Step{}, "empty.golden"},
		{"default settings", script.BashTemplate, false, false, false, steps, "default-settings.golden"},
		{"no echo", script.BashTemplate, true, false, false, steps, "no-echo.golden"},
		{"skip errors", script.BashTemplate, false, true, false, steps, "skip-errors.golden"},
		{"verbose", script.BashTemplate, true, true, true, steps, "verbose.golden"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			bash := &script.Bash{
				Template:   test.Template,
				NoEcho:     test.NoEcho,
				SkipErrors: test.SkipErrors,
				Verbose:    test.Verbose,
			}

			script, err := bash.Write(test.Steps)
			if err != nil {
				t.Fatalf("expected write err to be nil, got (%v)", err)
			}

			if !scriptMatchesGoldenFile(t, script, test.GoldenFile) {
				t.Fatalf("script output does not match golden file")
			}
		})
	}
}
