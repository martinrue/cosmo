package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

// Server holds data about a configured server.
type Server struct {
	Host string
	User string
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

	for name, server := range config.Servers {
		if server.Host == "" {
			return nil, fmt.Errorf("server '%s' is missing host key", name)
		}

		if server.User == "" {
			return nil, fmt.Errorf("server '%s' is missing user key", name)
		}
	}

	return config, nil
}
