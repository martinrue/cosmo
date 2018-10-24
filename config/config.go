package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

// Step holds data about a step in a task.
type Step struct {
	Exec string
}

// Task holds data about a set of tasks.
type Task struct {
	Steps []Step
}

// Server holds data about a configured server.
type Server struct {
	Host  string
	User  string
	Tasks map[string]Task
}

func (s *Server) String() string {
	return fmt.Sprintf("%s@%s", s.User, s.Host)
}

// Config holds top level data from the cosmo config file.
type Config struct {
	Servers map[string]Server
}

// Read reads and parses the cosmo config file.
func Read(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("missing config file path")
	}

	config := &Config{}

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	for serverName, server := range config.Servers {
		if server.Host == "" {
			return nil, fmt.Errorf("server '%s' is missing 'host' key", serverName)
		}

		if server.User == "" {
			return nil, fmt.Errorf("server '%s' is missing 'user' key", serverName)
		}

		for taskName, task := range server.Tasks {
			if len(task.Steps) == 0 {
				return nil, fmt.Errorf("task '%s' for server '%s' has no steps", taskName, serverName)
			}

			for i, step := range task.Steps {
				if step.Exec == "" {
					return nil, fmt.Errorf("step %d of task '%s' is missing 'exec' key", i+1, taskName)
				}
			}
		}
	}

	return config, nil
}
