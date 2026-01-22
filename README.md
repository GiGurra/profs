# profs

[![CI Status](https://github.com/GiGurra/profs/actions/workflows/ci.yml/badge.svg)](https://github.com/GiGurra/profs/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/GiGurra/profs)](https://goreportcard.com/report/github.com/GiGurra/profs)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://gigurra.github.io/profs/)

Manage multiple configuration profiles using symlinks. Switch between work/personal configs for `.gitconfig`, `.ssh`, `.kube`, and more.

## How It Works

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

## Installation

```bash
go install github.com/GiGurra/profs@latest
```

## Quick Start

```bash
# Add a path to manage (creates 'work' profile)
profs add ~/.gitconfig --profile work

# Create another profile
profs add-profile personal

# Switch profiles
profs set personal

# Check status
profs status
```

## Commands

| Command | Description |
|---------|-------------|
| `profs add <path>` | Add a path to manage |
| `profs add-profile <name>` | Create a new profile |
| `profs set <profile>` | Switch to a profile |
| `profs status` | Show current status |
| `profs list` | List all profiles |
| `profs doctor` | Check for issues |
| `profs remove <path>` | Stop managing a path |
| `profs remove-profile <name>` | Delete a profile |

## Example Usage

```bash
$ profs status
Profile: work
~/.gitconfig             -> ~/.gitconfig.profs/work [ok]
~/.ssh                   -> ~/.ssh.profs/work [ok]
~/.kube                  -> ~/.kube.profs/work [ok]
~/.config/gcloud         -> ~/.config/gcloud.profs/work [ok]
```

## Documentation

See the [full documentation](https://gigurra.github.io/profs/) for detailed guides.

## License

MIT
