# Cosmo

An agentless tool for interacting with bare-metal servers.

## Config

Cosmo expects to find a `cosmo.toml` config file in the working directory.
Use the `--config=<path>` flag to specify a different location.

### Example cosmo.toml
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

Usage: cosmo [--version] [--help] [--config=<path>] <command> [<args>]

Commands:
  servers  lists servers
  disk     shows disk space info
  uptime   shows uptime info
  version  displays the current cosmo version
```
