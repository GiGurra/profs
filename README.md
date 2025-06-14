# profs

[![CI Status](https://github.com/GiGurra/profs/actions/workflows/ci.yml/badge.svg)](https://github.com/GiGurra/profs/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/GiGurra/profs)](https://goreportcard.com/report/github.com/GiGurra/profs)

`profs` helps you manage multiple configuration profiles for different contexts (work vs personal, different clients,
etc.) by organizing them into `.profs` companion directories and using symlinks to switch between them.

For example:

- `~/.gitconfig` → `~/.gitconfig.profs/work` or `~/.gitconfig.profs/personal`
- `~/.ssh` → `~/.ssh.profs/client1` or `~/.ssh.profs/client2`

## Installation

```bash
go install github.com/GiGurra/profs@latest
```

or if you prefer a tool manager like [mise](https://github.com/jdx/mise), you can use:

```bash
mise use -g "go:github.com/GiGurra/profs@latest"
```

## Quick Start

1. Add a directory to be managed:
   ```bash
   profs add ~/.gitconfig --profile work
   ```
   The `--profile work` is only required the first time you add a path.
   It creates a new profile named "work" and sets up the symlink. For subsequent additions,
   you can omit the profile flag, and it will then use the current profile.

2. Add another profile:
   ```bash
   profs add-profile personal
   ```
   By default, this creates new empty directories for the new profile at all paths managed by `profs`.
   You can change this behavior by using the `--copy-existing` flag to copy the current profile's 
   contents into the new profile.

3. Switch between profiles:
   ```bash
   profs set work
   profs set personal
   ```

4. Check current status:
   ```bash
   > profs status
   Profile: personal
   ~/.config/gigurra/ai     -> ~/.config/gigurra/ai.profs/personal [ok]
   ~/.config/gh             -> ~/.config/gh.profs/personal [ok]
   ~/.ssh                   -> ~/.ssh.profs/personal [ok]
   ~/.kube                  -> ~/.kube.profs/personal [ok]
   ~/.gitconfig             -> ~/.gitconfig.profs/personal [ok]
   ~/.config/github-copilot -> ~/.config/github-copilot.profs/personal [ok]
   ~/.config/gcloud         -> ~/.config/gcloud.profs/personal [ok]
   ```

## How it works

- Each managed path gets a companion `.profs` directory containing profile subdirectories
- The original path becomes a symlink pointing to the active profile
- Switching profiles updates all symlinks simultaneously

Example structure:

```
~/.gitconfig -> ~/.gitconfig.profs/work
~/.gitconfig.profs/
├── work
└── personal
~/.ssh -> ~/.ssh.profs/work
~/.ssh.profs/
├── work/
└── personal/
```

## Commands

- `profs add <path>` - Add a directory/file to be managed
- `profs add-profile <name>` - Create a new profile for all managed paths
- `profs set <profile>` - Switch to a specific profile
- `profs list` - Show all available profiles
- `profs status` - Show current profile status
- `profs doctor` - Check for configuration issues
- `profs remove <path>` - Remove a path from management
- `profs remove-profile <name>` - Delete a profile

## Configuration

Configuration is stored at `~/.config/gigurra/profs/global.json` and contains the list of managed paths.
Example:
```bash
> cat ~/.config/gigurra/profs/global.json
{
  "paths": [
    "~/.config/gigurra/ai",
    "~/.config/gh",
    "~/.ssh",
    "~/.kube",
    "~/.gitconfig",
    "~/.config/github-copilot",
    "~/.config/gcloud"
  ]
}
```

## Shell Completion

Generate completion scripts:

```bash
profs completion bash|zsh|fish|powershell
```

## Full command reference

run `profs --help` to see all available commands and options:

```bash
> profs --help
Manage user profiles

Usage:
  profs [command]

Available Commands:
  add                Adds a new directory to be managed by profs
  add-path           Adds a new directory to be managed by profs
  add-profile        Adds a new profile to be managed by profs
  completion         Generate the autocompletion script for the specified shell
  doctor             Show inconsistencies in the current configuration
  help               Help about any command
  list               Lists all detected profiles
  list-profiles      Lists all detected profiles
  migrate-config-dir Migrate legacy /Users/johkjo/.profs -> /Users/johkjo/.config/gigurra/profs
  remove             Removes a directory from profs config
  remove-path        Removes a directory from profs config
  remove-profile     Removes an existing profile managed by profs
  reset              Resets all configuration to zero
  set                Set current profile
  status             Show current status
  status-config      Show current raw configuration
  status-full        Show full status and alternatives
  status-profile     Show current profile status

Flags:
  -h, --help   help for profs

Use "profs [command] --help" for more information about a command.
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Pull requests welcome! The project includes CI/CD with automated testing and releases.