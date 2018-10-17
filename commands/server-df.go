package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/martinrue/cosmo/config"
	"github.com/martinrue/cosmo/drawing"
	"github.com/martinrue/cosmo/ssh"
)

type diskUsage struct {
	Volume      string
	Size        string
	Used        string
	Available   string
	UsedPercent float64
}

func (du *diskUsage) render() string {
	progress := drawing.ProgressBar(du.UsedPercent, 30)

	return fmt.Sprintf("%s (%s)\nused: %s free: %s\n%s %0.0f%% used",
		du.Volume,
		du.Size,
		du.Used,
		du.Available,
		progress,
		du.UsedPercent)
}

func newDiskUsage(usage string) (*diskUsage, error) {
	du := &diskUsage{}

	lines := strings.Split(usage, "\n")

	du.Volume = lines[0]
	du.Size = lines[1]
	du.Used = lines[2]
	du.Available = lines[3]

	percent, err := strconv.ParseFloat(strings.TrimSuffix(lines[4], "%"), 64)
	if err != nil {
		return nil, errors.New("cannot parse disk usage string from server")
	}

	du.UsedPercent = percent

	return du, nil
}

// ServerDf is a server-df command.
type ServerDf struct {
	Config *config.Config
	Server string
	All    bool
	Raw    bool
}

func (cmd *ServerDf) remoteCommand() string {
	if cmd.Raw {
		return "df -hl /"
	}

	return "df -hl / | awk 'FNR == 2 {print $1 \"\\n\" $2 \"\\n\" $3 \"\\n\" $4 \"\\n\" $5}'"
}

func (cmd *ServerDf) df(server config.Server) (string, error) {
	client := &ssh.Client{
		Host: server.String(),
	}

	response, err := client.Exec(cmd.remoteCommand())
	if err != nil {
		return "", err
	}

	if cmd.Raw {
		return response, nil
	}

	usage, err := newDiskUsage(response)
	if err != nil {
		return "", err
	}

	return usage.render(), nil
}

// Exec handles running of the sub command.
func (cmd *ServerDf) Exec() {
	if !cmd.All {
		server, ok := cmd.Config.Servers[cmd.Server]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: server '%s' not found, check cosmo.toml\n", cmd.Server)
			return
		}

		response, err := cmd.df(server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}

		fmt.Println(response)
		return
	}

	for name, server := range cmd.Config.Servers {
		response, err := cmd.df(server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}

		fmt.Printf("[%s]\n%s\n\n", name, response)
	}
}

// NewServerDf creates a new server-df command.
func NewServerDf(config *config.Config) *ServerDf {
	if len(os.Args) < 2 || os.Args[1] != "server-df" {
		return nil
	}

	flags := flag.NewFlagSet("server-df", flag.ExitOnError)
	server := flags.String("server", "", "")
	all := flags.Bool("all", false, "")
	raw := flags.Bool("raw", false, "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo server-df --server=<id> [--raw] [--all]")
	}

	flags.Parse(os.Args[2:])

	if *server == "" && !*all {
		flags.Usage()
		os.Exit(1)
	}

	return &ServerDf{
		Config: config,
		Server: *server,
		All:    *all,
		Raw:    *raw,
	}
}
