package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/ssh"
)

// ServerUptime is a server-uptime command.
type ServerUptime struct {
	Config *config.Config
	Server string
	All    bool
}

func (cmd *ServerUptime) uptime(server config.Server) (string, error) {
	client := &ssh.Client{
		Host: server.String(),
	}

	response, err := client.Exec("uptime | awk '{print $3 \"\\n\" $5}'")
	if err != nil {
		return "", err
	}

	lines := strings.Split(response, "\n")
	days := strings.TrimSpace(lines[0])
	hours := strings.TrimSpace(strings.TrimSuffix(lines[1], ","))

	return fmt.Sprintf("%s days, %s hours", days, hours), nil
}

// Exec handles running of the sub command.
func (cmd *ServerUptime) Exec() {
	if !cmd.All {
		server, ok := cmd.Config.Servers[cmd.Server]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: server '%s' not found, check cosmo.toml\n", cmd.Server)
			return
		}

		response, err := cmd.uptime(server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}

		fmt.Println(response)
		return
	}

	for name, server := range cmd.Config.Servers {
		response, err := cmd.uptime(server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}

		fmt.Printf("[%s]\n%s\n\n", name, response)
	}
}

// NewServerUptime creates a new server-df command.
func NewServerUptime(config *config.Config) *ServerUptime {
	if len(os.Args) < 2 || os.Args[1] != "server-uptime" {
		return nil
	}

	flags := flag.NewFlagSet("server-uptime", flag.ExitOnError)
	server := flags.String("server", "", "")
	all := flags.Bool("all", false, "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo server-uptime --server=<id> [--all]")
	}

	flags.Parse(os.Args[2:])

	if *server == "" && !*all {
		flags.Usage()
		os.Exit(1)
	}

	return &ServerUptime{
		Config: config,
		Server: *server,
		All:    *all,
	}
}
