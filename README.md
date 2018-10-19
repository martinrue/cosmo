# Cosmo

An agentless tool for interacting with bare-metal servers.

## Config

Cosmo requires a `cosmo.toml` config file located in the working directory.
The `cosmo.toml` file should specify a list of servers to interact with.

```toml
[servers]

  [servers.jerry]
  host = "jerry.com"
  user = "kel"

  [servers.elaine]
  host = "elaine.com"
  user = "susie"

  [servers.george]
  host = "george.com"
  user = "art"
```

## Commands

```
Cosmo

Usage: cosmo <command> [<args>]

Commands:
  servers  lists servers
  disk     shows disk space info
  uptime   shows uptime info
  version  displays the current cosmo version
```
