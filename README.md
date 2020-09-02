你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# Cosmo

[![CircleCI](https://circleci.com/gh/martinrue/cosmo.svg?style=svg)](https://circleci.com/gh/martinrue/cosmo)

An agentless tool for interacting with bare-metal servers.

## Config

Cosmo expects to find a `cosmo.toml` config file in the working directory.
Use the `--config=<path>` flag to specify a different location.

### Example cosmo.toml
```toml
[servers.jerry]
host = 'jerry.com'
user = 'kel'

  [servers.jerry.tasks.kill]
    local = [
      { exec = 'pkill -9 newman' },
    ]

  [servers.jerry.tasks.date]
    no_echo = true

    local = [
      { exec = 'echo "local time:"' },
      { exec = 'date' },
    ]

    remote = [
      { exec = 'echo "remote time:"' },
      { exec = 'date' },
    ]

[servers.elaine]
host = 'elaine.com'
user = 'susie'

  [servers.elaine.tasks.ls]
    no_echo = true
    skip_errors = true

    local_raw = '''
      echo "local directories:"
      ls /var
      ls /usr
      ls /etc
    '''

    remote_raw = '''
      echo "remote directories:"
      ls /var
      ls /usr
      ls /etc
    '''

[servers.george]
host = 'george.com'
user = 'art'
```

## Config Options

### Tasks

```
[servers.<server-name>.tasks.<task-name>]
no_echo = true
```

Disables echoing for all local and remote steps. Only `stdout` and `stderr` of each step will be ouput.

A step may individually specify `no_echo` to enable it on a step-by-step basis.

```
[servers.<server-name>.tasks.<task-name>]
skip_errors = true
```

Disables "break on error" behaviour for all steps. If a step fails, the next step will still be executed.

A step may individually specify `skip_error` to enable it on a step-by-step basis.

### Steps

Local steps (which are run on the local machine) should be specified under the `local` key of a task:

```
[servers.<server-name>.tasks.<task-name>]
  local = [
    # steps
  ]
```

Remote steps (which are run on the remote server) should be specified under the `remote` key:

```
[servers.<server-name>.tasks.<task-name>]
  remote = [
    # steps
  ]
```

Each step should be an object containing at least an `exec` field:

```
{ exec = 'echo "Hello World"' },
```

A step may optionally specify the boolean fields `no_echo` and `skip_error` to override default behaviour.

Local and remote steps can alternatively be specified as line-delimited strings, using the `local_raw` and `remote_raw` keys respectively:

```
[servers.<server-name>.tasks.<task-name>]
  local_raw = '''
    echo "local"
    echo "commands"
    echo "here"
  '''

  remote_raw = '''
    echo "remote"
    echo "commands"
    echo "here"
  '''
```

### Commands

Each `exec` field of a step should specify a shell command to execute.

If a command begins with the `@` character, `stdout` will be suppressed.

If a command begins with the `!` character, both `stdout` and `stderr` will be suppressed, silencing all output of the command.

## CLI

```
Cosmo

Usage: cosmo [--version] [--help] [--config=<path>] <command> [<args>]

Commands:
  run      runs a task
  servers  lists servers
  steps    lists the steps of a task
  tasks    lists tasks
```
