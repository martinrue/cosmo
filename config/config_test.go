package config_test

import (
	"testing"

	"github.com/martinrue/cosmo/config"
)

func TestConfigReadNoPath(t *testing.T) {
	_, err := config.Read("")

	expected := "no config path"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadInvalidData(t *testing.T) {
	_, err := config.Read("./testdata/invalid.toml")
	if err == nil {
		t.Fatal("expected read err")
	}
}

func TestConfigReadEmpty(t *testing.T) {
	_, err := config.Read("./testdata/empty.toml")

	expected := "no servers"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadNoHost(t *testing.T) {
	_, err := config.Read("./testdata/no-host.toml")

	expected := "server 'no-host' is missing 'host' key"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadNoUser(t *testing.T) {
	_, err := config.Read("./testdata/no-user.toml")

	expected := "server 'no-user' is missing 'user' key"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadNoTasks(t *testing.T) {
	_, err := config.Read("./testdata/no-tasks.toml")

	expected := "no tasks for server 'no-tasks'"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadNoSteps(t *testing.T) {
	_, err := config.Read("./testdata/no-steps.toml")

	expected := "task 'task' for server 'no-steps' has no local or remote steps"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadInvalidLocalSteps(t *testing.T) {
	_, err := config.Read("./testdata/invalid-local-steps.toml")

	expected := "step 2 of task 'task' is missing 'exec' key"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadInvalidRemoteSteps(t *testing.T) {
	_, err := config.Read("./testdata/invalid-remote-steps.toml")

	expected := "step 3 of task 'task' is missing 'exec' key"

	if err != nil && err.Error() != expected {
		t.Fatalf("expected read err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigReadProcessesRawSteps(t *testing.T) {
	cfg, err := config.Read("./testdata/raw-steps.toml")
	if err != nil {
		t.Fatalf("expected read err to be nil, got (%v)", err)
	}

	server, ok := cfg.Servers["raw-steps"]
	if !ok {
		t.Fatalf("expected server map to 'raw-steps' server")
	}

	task, ok := server.Tasks["deploy"]
	if !ok {
		t.Fatalf("expected 'raw-steps' server to have 'deploy' task")
	}

	if len(task.Local) != 1 {
		t.Fatalf("expected 'deploy' task to have 1 local step")
	}

	if task.Local[0].Exec != `echo "Hello"` {
		t.Fatalf("unexpected first local step of 'deploy' task")
	}

	if len(task.Remote) != 2 {
		t.Fatalf("expected 'deploy' task to have 2 remote steps")
	}

	if task.Remote[0].Exec != `ls /etc` || task.Remote[1].Exec != `ls /usr` {
		t.Fatalf("unexpected remote steps of 'deploy' task")
	}
}

func TestConfigFindTaskDetectsAmbiguousTask(t *testing.T) {
	cfg, err := config.Read("./testdata/ambiguous-task.toml")
	if err != nil {
		t.Fatalf("expected read err to be nil, got (%v)", err)
	}

	expected := "task 'task-1' is ambiguous, specify server"

	_, _, err = cfg.Servers.FindTask("task-1", "")
	if err != nil && err.Error() != expected {
		t.Fatalf("expected find task err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigFindTaskIgnoresAmbiguousTaskWhenServerSpecified(t *testing.T) {
	cfg, err := config.Read("./testdata/ambiguous-task.toml")
	if err != nil {
		t.Fatalf("expected read err to be nil, got (%v)", err)
	}

	task, server, err := cfg.Servers.FindTask("task-1", "server-1")
	if err != nil {
		t.Fatalf("expected find task err to be nil, got (%v)", err)
	}

	if task.ServerName != "server-1" {
		t.Fatalf("expected find task to return task from 'server-1', got (%v)", task.ServerName)
	}

	if task.Local[0].Exec != `echo "from server-1"` {
		t.Fatalf("expected find task to return correct task")
	}

	if server.Host != "host-1" {
		t.Fatalf("expected find task to return correct server")
	}
}

func TestConfigFindTaskFailsIfTaskNotFound(t *testing.T) {
	cfg, err := config.Read("./testdata/valid.toml")
	if err != nil {
		t.Fatalf("expected read err to be nil, got (%v)", err)
	}

	expected := "task 'missing' not found, check config"

	_, _, err = cfg.Servers.FindTask("missing", "")
	if err != nil && err.Error() != expected {
		t.Fatalf("expected find task err to be (%v), got (%v)", expected, err)
	}
}

func TestConfigServerStringerReturnsUserHostCombiation(t *testing.T) {
	cfg, err := config.Read("./testdata/server-stringer.toml")
	if err != nil {
		t.Fatalf("expected read err to be nil, got (%v)", err)
	}

	_, server, err := cfg.Servers.FindTask("task", "")
	if err != nil {
		t.Fatalf("expected find task err to be nil, got (%v)", err)
	}

	if server.String() != "cosmo@kramerica.industries" {
		t.Fatalf("expected server string to return 'cosmo@kramerica.industries', got (%s)", server.String())
	}
}
