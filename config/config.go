package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

// Command holds data about a command in a task.
type Command struct {
	Exec string
}

// Task holds data about a set of tasks.
type Task struct {
	Local  []Command
	Remote []Command
}

// TaskMap is a dictionary of tasks.
type TaskMap map[string]Task

// Server holds data about a configured server.
type Server struct {
	Host  string
	User  string
	Tasks TaskMap
}

// ServerMap is a dictionary of servers.
type ServerMap map[string]Server

// FindTask searches for a task by name.
func (m ServerMap) FindTask(taskName string, serverName string) (Task, error) {
	var task Task
	found := false

	for server, s := range m {
		if serverName != "" && serverName != server {
			continue
		}

		for name, t := range s.Tasks {
			if taskName == name {
				if found {
					return Task{}, fmt.Errorf("error: task '%s' is ambiguous, specify server", taskName)
				}

				found = true
				task = t
			}
		}
	}

	if !found {
		return Task{}, fmt.Errorf("error: task '%s' not found, check config", taskName)
	}

	return task, nil
}

func (s *Server) String() string {
	return fmt.Sprintf("%s@%s", s.User, s.Host)
}

// Config holds top level data from the cosmo config file.
type Config struct {
	Servers ServerMap
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

		for taskName, task := range server.Tasks {
			if len(task.Local) == 0 && len(task.Remote) == 0 {
				return config, fmt.Errorf("task '%s' for server '%s' has no local or remote commands", taskName, serverName)
			}

			if len(task.Local) > 0 {
				for i, cmd := range task.Local {
					if cmd.Exec == "" {
						return config, fmt.Errorf("local command %d of task '%s' is missing 'exec' key", i+1, taskName)
					}
				}
			}

			if len(task.Remote) > 0 {
				for i, cmd := range task.Remote {
					if cmd.Exec == "" {
						return config, fmt.Errorf("remote command %d of task '%s' is missing 'exec' key", i+1, taskName)
					}
				}
			}
		}
	}

	return config, nil
}
