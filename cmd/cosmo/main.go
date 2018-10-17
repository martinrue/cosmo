package main

import (
	"fmt"
	"os"

	"github.com/martinrue/cosmo/commands"
	"github.com/martinrue/cosmo/config"
)

const usage = `Cosmo

Usage: cosmo <command> [<args>]

Commands:
  server-ls  lists known servers
  server-df  shows disk usage info for a server
  version    displays the current cosmo version
`

func main() {
	usageAndExit := func(code int) {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(code)
	}

	if len(os.Args) == 1 {
		usageAndExit(1)
	}

	if os.Args[1] == "--help" {
		usageAndExit(0)
	}

	conf, err := config.Read("cosmo.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %s\n", err)
		os.Exit(1)
	}

	commands := map[string]commands.Command{
		"server-ls": commands.NewServerLs(conf),
		"server-df": commands.NewServerDf(conf),
		"version":   commands.NewVersion(conf),
	}

	command, ok := commands[os.Args[1]]
	if !ok {
		fmt.Fprintf(os.Stderr, "cosmo: '%s' is not a cosmo command. See 'cosmo --help'.\n", os.Args[1])
		os.Exit(1)
	}

	command.Exec()
}
