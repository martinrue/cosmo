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

// CommandDisk displays disk space info.
type CommandDisk struct {
	Config *config.Config
	Server string
	All    bool
	Raw    bool
}

func (cmd *CommandDisk) remoteCommand() string {
	if cmd.Raw {
		return "df -hl /"
	}

	return "df -hl / | awk 'FNR == 2 {print $1 \"\\n\" $2 \"\\n\" $3 \"\\n\" $4 \"\\n\" $5}'"
}

func (cmd *CommandDisk) df(name string, server config.Server) (string, error) {
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

// Exec runs the subcommand.
func (cmd *CommandDisk) Exec() {
	table := &drawing.Table{}

	if !cmd.All {
		server, ok := cmd.Config.Servers[cmd.Server]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: server '%s' not found, check cosmo.toml\n", cmd.Server)
			return
		}

		response, err := cmd.df(cmd.Server, server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}

		table.AddRows(fmt.Sprintf("%s:", cmd.Server), response)
		fmt.Println(table)
		return
	}

	for name, server := range cmd.Config.Servers {
		response, err := cmd.df(name, server)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}

		table.AddRows(fmt.Sprintf("%s:", name), response)
	}

	fmt.Println(table)
}

// NewCommandDisk creates a new 'disk' subcommand.
func NewCommandDisk(config *config.Config) *CommandDisk {
	if len(os.Args) < 2 || os.Args[1] != "disk" {
		return nil
	}

	flags := flag.NewFlagSet("disk", flag.ExitOnError)
	server := flags.String("server", "", "")
	all := flags.Bool("all", false, "")
	raw := flags.Bool("raw", false, "")

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cosmo disk --server=<id> [--raw] [--all]")
	}

	flags.Parse(os.Args[2:])

	if *server == "" && !*all {
		flags.Usage()
		os.Exit(1)
	}

	return &CommandDisk{
		Config: config,
		Server: *server,
		All:    *all,
		Raw:    *raw,
	}
}
