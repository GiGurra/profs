# profs

Manage multiple configuration profiles using symlinks. Switch between work/personal configs instantly.

## The Problem

You have different configurations for different contexts:

- **Work**: Corporate git identity, VPN SSH keys, production Kubernetes clusters
- **Personal**: Personal git identity, hobby project SSH keys, home lab clusters
- **Client A vs Client B**: Different cloud credentials, different git identities

Manually swapping these files is error-prone and tedious.

## The Solution

profs manages symlinks to switch entire configuration sets at once:

```
~/.gitconfig -> ~/.gitconfig.profs/work
~/.ssh       -> ~/.ssh.profs/work
~/.kube      -> ~/.kube.profs/work
```

One command switches everything:

```bash
profs set personal
```

Now all symlinks point to the `personal` profile.

## Features

| Feature | Description |
|---------|-------------|
| Atomic switching | All paths switch together |
| Any path | Files or directories |
| Profile detection | Auto-discovers available profiles |
| Health checks | `profs doctor` finds issues |
| Shell completion | Tab completion for profiles |

## Quick Example

```bash
# Start managing your git config
profs add ~/.gitconfig --profile work

# Add more paths
profs add ~/.ssh
profs add ~/.kube

# Create a personal profile
profs add-profile personal

# Switch to personal
profs set personal

# Check status
profs status
# Profile: personal
# ~/.gitconfig -> ~/.gitconfig.profs/personal [ok]
# ~/.ssh       -> ~/.ssh.profs/personal [ok]
# ~/.kube      -> ~/.kube.profs/personal [ok]
```

## Installation

```bash
go install github.com/GiGurra/profs@latest
```

Or with [mise](https://github.com/jdx/mise):

```bash
mise use -g "go:github.com/GiGurra/profs@latest"
```

## Next Steps

- [Getting Started](guide/getting-started.md) - First-time setup
- [Managing Paths](guide/paths.md) - Adding and removing paths
- [Working with Profiles](guide/profiles.md) - Creating and switching profiles
- [Use Cases](guide/use-cases.md) - Common scenarios
