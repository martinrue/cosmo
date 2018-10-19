package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
)

const usage = `Cosmo

Usage: cosmo [--version] [--help] [--config=<path>] <command> [<args>]

Commands:
  servers  lists servers
  disk     shows disk space info
  uptime   shows uptime info
  version  displays the current cosmo version
`

var (
	version = flag.Bool("version", false, "")
	help    = flag.Bool("help", false, "")
	conf    = flag.String("config", "", "")
)

func main() {
	usageAndExit := func(code int) {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(code)
	}

	flag.Usage = func() {
		usageAndExit(1)
	}

	flag.Parse()

	if *help {
		usageAndExit(0)
	}

	if *version {
		fmt.Println("v0.0.1")
		os.Exit(0)
	}

	configPath := "cosmo.toml"

	if *conf != "" {
		configPath = *conf
	}

	conf, err := config.Read(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %s\n", err)
		os.Exit(1)
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Missing <command>. See 'cosmo --help'.\n")
		os.Exit(1)
	}

	ctors := map[string]commands.Ctor{
		"servers": commands.NewCommandServers,
		"disk":    commands.NewCommandDisk,
		"uptime":  commands.NewCommandUptime,
	}

	ctor, ok := ctors[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "'%s' is not a cosmo command. See 'cosmo --help'.\n", args[0])
		os.Exit(1)
	}

	ctor(conf, args[1:]).Exec()
}
