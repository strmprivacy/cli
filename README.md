# Stream Machine Command Line Interface
[![GitHub Actions](https://github.com/streammachineio/cli/workflows/Build/badge.svg)](https://github.com/streammachineio/cli/actions)
[![Latest Release](https://img.shields.io/github/v/release/streammachineio/cli)](https://github.com/streammachineio/cli/releases/latest)

This package contains a command line interface (CLI) for interacting with https://www.streammachine.io [Stream Machine].

## Installation

### Manually
Download the latest release for your platform from the [releases page](https://github.com/streammachineio/cli/releases/latest).
Put the binary somewhere on your path.

#### Shell Completion

In order to set up command completion, please follow the instructions below:
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

The CLI is also available through Homebrew. Install the formula as follows:
```
brew install streammachineio/cli/strm
```

Ensure to read the caveats section for setting up command completion.

Upgrades to the CLI can be done through `brew upgrade strm`.

### Other package managers

More package managers will be added in the future, so stay tuned.

## Configuration

The `strm` CLI can be configured using either the flags as specified by the help (as command line arguments), with environment variables, or with a configuration file, named strm.yaml, located in the Stream Machine [Configuration directory]. If a flag is not present, the default value is used.

Note: The ordering is the same as specified above, so arguments take precedence over environment variables, which take precedence over the configuration file, which takes precedence over the default values.

### Default configuration values

Below are the default values for all `strm' flags, in the YAML format used by the configuration file:

```yaml
save: true
event-auth-host: https://auth.strm.services
events-gateway: https://in.strm.services/event
api-auth-url: https://api.streammachine.io/v1
api-host: apis.streammachine.io:443
```

In normal circumstances, these defaults should work and there is no need to create this configuration file and override any Flags. It can be useful in special cases, for example if you'd like to use a mock endpoint for testing.

#### Save
The option `save` indicates whether the output of create commands is saved to files in your Stream Machine [Configuration directory]. This can be very useful, especially when managing many resources which have secrets. Be wary however that this information is stored in plain-text on your file system.

#### event-auth-host etc TODO


### Configuration directory
The Stream Machine CLI stores it's information in a configuration directory, by default located in:
`$HOME/.config/stream-machine/`. In this directory, the CLI looks for a file named: `strm.yaml`, which is used for setting global flags.

By default, this directory also contains the login information used by the `strm auth` commands, in a file named: `strm-creds-<api-auth-url>.json`. This file is generated and updated by the CLI, so there is no need for any manual editing.

In this directory you can also find all entities that have been `save`d (see the [Save] option).
These entities are saved in the following files: `<config-dir>/<Entity>/<name>.json`, where `Entity` is the Entity name, i.e. "Stream" or "Sink" and the `name` is the unique name of the created entity, i.e. "MyImportantStream" or "s3-sink".

## Getting help
If you encounter an error, or you'd like a new feature, please create an issue here[https://github.com/streammachineio/cli-wip/issues/new]. Please be thorough in your description, as it helps us to help you more quickly. At least include the version of the CLI, your OS. terminal and any custom Stream Machine flags that are present in your config or environment.

---
**IMPORTANT**

Don't provide the login configuration JSON file, as it includes sensitive information!

---

We’re also frequently checking our https://gitter.im/stream-machine/community[Gitter] channel, and others in the Stream Machine community may be able to help you as well.

Or email developer-support@streammachine.io[Developer Support], with the details of the issue you’re experiencing. A minimum working example (MWE) would help us in reproducing the issue, and could help in solving it sooner for you. If you have to option to include an MWE, please do so.



## More resources

See our [documentation](https://docs.streammachine.io) or [reach out to us](https://docs.streammachine.io/docs/0.3.4/contact/index.html).
