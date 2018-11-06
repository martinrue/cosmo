package runner_test

import (
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/martinrue/cosmo/runner"
)

func createMockExec() (*string, *string, runner.Executor) {
	execName := ""
	execArgs := ""

	return &execName, &execArgs, func(command *exec.Cmd) error {
		execName = command.Path
		execArgs = strings.Join(command.Args, " ")
		return nil
	}
}

func TestLocalRunner(t *testing.T) {
	name, args, exec := createMockExec()

	local := &runner.Local{Exec: exec}
	_ = local.Run("<local script>")

	if path.Base(*name) != "bash" {
		t.Fatalf("expected runner to exec (bash), got (%v)", *name)
	}

	expected := "bash -c <local script>"

	if *args != expected {
		t.Fatalf("expected runner to use args (%v), got (%v)", expected, *args)
	}
}

func TestRemoteRunner(t *testing.T) {
	name, args, exec := createMockExec()

	remote := &runner.Remote{Exec: exec, Host: "user@host.domain"}
	_ = remote.Run("<remote script>")

	if path.Base(*name) != "ssh" {
		t.Fatalf("expected runner to exec (ssh), got (%v)", *name)
	}

	expected := "ssh user@host.domain bash -s <<COSMO <remote script> \nCOSMO"

	if *args != expected {
		t.Fatalf("expected runner to use args (%s), got (%v)", expected, *args)
	}
}

func TestExec(t *testing.T) {
	cmd := exec.Command("go", "version")

	if err := runner.Exec(cmd); err != nil {
		t.Fatalf("expected successful exit, got (%v)", err)
	}
}

func TestExecStderrPipeFailure(t *testing.T) {
	cmd := exec.Command("go", "version")
	cmd.Stderr = ioutil.Discard

	expected := "exec: Stderr already set"

	if err := runner.Exec(cmd); err != nil && err.Error() != expected {
		t.Fatalf("expected exec err (%v), got (%v)", expected, err)
	}
}

func TestExecStdoutPipeFailure(t *testing.T) {
	cmd := exec.Command("go", "version")
	cmd.Stdout = ioutil.Discard

	expected := "exec: Stdout already set"

	if err := runner.Exec(cmd); err != nil && err.Error() != expected {
		t.Fatalf("expected exec err (%v), got (%v)", expected, err)
	}
}
