package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// Step holds data about a command in a task.
type Step struct {
	Exec      string `toml:"exec"`
	NoEcho    bool   `toml:"no_echo"`
	SkipError bool   `toml:"skip_error"`
}

// Task holds data about a set of tasks.
type Task struct {
	ServerName string `toml:"-"`
	NoEcho     bool   `toml:"no_echo"`
	SkipErrors bool   `toml:"skip_errors"`
	Local      []Step `toml:"local"`
	LocalRaw   string `toml:"local_raw"`
	Remote     []Step `toml:"remote"`
	RemoteRaw  string `toml:"remote_raw"`
}

// TaskMap is a dictionary of tasks.
type TaskMap map[string]Task

// Server holds data about a configured server.
type Server struct {
	Host  string  `toml:"host"`
	User  string  `toml:"user"`
	Tasks TaskMap `toml:"tasks"`
}

func (s *Server) String() string {
	return fmt.Sprintf("%s@%s", s.User, s.Host)
}

// ServerMap is a dictionary of servers.
type ServerMap map[string]Server

// FindTask searches for a task by name.
func (sm ServerMap) FindTask(taskName string, serverName string) (Task, Server, error) {
	var foundTask Task
	var foundServer Server

	found := false

	for server, s := range sm {
		if serverName != "" && serverName != server {
			continue
		}

		for name, t := range s.Tasks {
			if taskName == name {
				if found {
					return Task{}, Server{}, fmt.Errorf("error: task '%s' is ambiguous, specify server", taskName)
				}

				found = true
				foundTask = t
				foundTask.ServerName = server
				foundServer = s
			}
		}
	}

	if !found {
		return Task{}, Server{}, fmt.Errorf("error: task '%s' not found, check config", taskName)
	}

	return foundTask, foundServer, nil
}

// Config holds top level data from the cosmo config file.
type Config struct {
	Servers ServerMap
}

func validateTasks(serverName string, tasks TaskMap) error {
	for taskName, task := range tasks {
		if len(task.Local) == 0 && len(task.Remote) == 0 {
			return fmt.Errorf("task '%s' for server '%s' has no local or remote steps", taskName, serverName)
		}

		if err := validateSteps(taskName, task.Local); err != nil {
			return err
		}

		if err := validateSteps(taskName, task.Remote); err != nil {
			return err
		}
	}

	return nil
}

func validateSteps(taskName string, steps []Step) error {
	for i, step := range steps {
		if step.Exec == "" {
			return fmt.Errorf("step %d of task '%s' is missing 'exec' key", i+1, taskName)
		}
	}

	return nil
}

func processRawSteps(server *Server) {
	buildSteps := func(raw string) []Step {
		steps := make([]Step, 0)

		if raw == "" {
			return steps
		}

		for _, command := range strings.Split(raw, "\n") {
			command = strings.TrimSpace(command)

			if command == "" || command[0:1] == "#" {
				continue
			}

			steps = append(steps, Step{
				Exec: command,
			})
		}

		return steps
	}

	for taskName, task := range server.Tasks {
		localSteps := buildSteps(task.LocalRaw)

		if len(localSteps) > 0 {
			task.Local = localSteps
		}

		remoteSteps := buildSteps(task.RemoteRaw)

		if len(remoteSteps) > 0 {
			task.Remote = remoteSteps
		}

		server.Tasks[taskName] = task
	}
}

// Read reads and parses the cosmo config file.
func Read(path string) (Config, error) {
	config := Config{}

	if path == "" {
		return config, errors.New("missing config file path")
	}

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return config, err
	}

	for serverName, server := range config.Servers {
		if server.Host == "" {
			return config, fmt.Errorf("server '%s' is missing 'host' key", serverName)
		}

		if server.User == "" {
			return config, fmt.Errorf("server '%s' is missing 'user' key", serverName)
		}

		processRawSteps(&server)

		if err := validateTasks(serverName, server.Tasks); err != nil {
			return config, err
		}
	}

	return config, nil
}
