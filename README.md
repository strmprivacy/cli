# Stream Machine Command Line Interface
[![GitHub Actions](https://github.com/streammachineio/cli/workflows/Build/badge.svg)](https://github.com/streammachineio/cli/actions)
[![Latest Release](https://img.shields.io/github/v/release/streammachineio/cli)](https://github.com/streammachineio/cli/releases/latest)

CLI for interacting with Stream Machine.

```
$ strm --help

Stream Machine CLI

Usage:
  strm [command]

Available Commands:
  auth        Authentication command
  completion  Generate completion script
  create      Create an entity
  delete      Delete an entity
  egress      Read from egress
  get         Get an entity
  help        Help about any command
  list        List entities
  sim         Simulate events
  version     Print CLI version

Flags:
      --api-auth-url string      Auth URL for user logins (default "https://api.streammachine.io/v1")
      --api-host string          API host name (default "apis.streammachine.io:443")
      --config-path string       config path (default is $HOME/.config/stream-machine/)
      --event-auth-host string   Security Token Service for events (default "auth.strm.services")
  -h, --help                     help for strm
      --token-file string        config file (default is $HOME/.config/stream-machine/strm-creds-<api-auth-host>.json)

Use "strm [command] --help" for more information about a command.
```

--- 

## Installation

### Manually
Download the latest release for your platform from the [releases page](https://github.com/streammachineio/cli/releases/latest).
Put the binary somewhere on your path.

#### Shell Completion

In order to set up auto completion, please follow the instructions below:
- for `bash` users \
  add the following line to your `.bash_profile` or `.bashrc`:
  `source <(strm completion bash)`
  or, to load completions for each session, execute once:
  - Linux users: `strm completion bash > /etc/bash_completion.d/strm`
  - macOS users: `strm completion bash > /usr/local/etc/bash_completion.d/strm`
- for `zsh` users \
  ensure that shell completion is enabled, then run (only needs to be done once):
  `/bin/zsh -c 'strm completion zsh > "${fpath[1]}/_strm"'`
- for fish users \
  `source "strm/path.fish.inc"`

### Homebrew

The CLI is available through Homebrew. Install the formula as follows:
```
brew install streammachineio/cli/strm
```

Ensure to read the caveats section for setting up auto complete. Upgrades to the CLI can be done through `brew upgrade strm`.

### Other package managers

More package managers will be added in the future, so stay tuned.

## Configuration

The `strm` CLI can be configured using either flags as specified by the help, or with environment variables, but also with a configuration file. The configuration file is located at `$HOME/.config/stream-machine/strm-creds-<api-auth-host>.json`. If it doesn't exist, you can create it and override properties there.

### Reference

Below is a reference of the configuration file for `strm`:

```yaml
save: true
event-auth-host: https://auth.strm.services
events-gateway: https://in.strm.services/event
api-auth-url: https://api.streammachine.io/v1
api-host: apis.streammachine.io:443
```

## Need help?

See our [documentation](https://docs.streammachine.io) or [reach out to us](https://docs.streammachine.io/docs/0.3.4/contact/index.html).
