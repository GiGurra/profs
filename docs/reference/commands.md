# Command Reference

## Path Management

### profs add

Add a path to be managed.

```bash
profs add <path> [--profile <name>]
```

| Flag | Description |
|------|-------------|
| `--profile` | Profile name (required for first path) |

**Examples:**

```bash
# First path (creates profile)
profs add ~/.gitconfig --profile work

# Subsequent paths
profs add ~/.ssh
profs add ~/.kube
```

**Aliases:** `profs add-path`

### profs remove

Remove a path from management.

```bash
profs remove <path>
```

Restores the original path with current profile's content.

**Aliases:** `profs remove-path`

## Profile Management

### profs add-profile

Create a new profile.

```bash
profs add-profile <name> [--copy-existing]
```

| Flag | Description |
|------|-------------|
| `--copy-existing` | Copy current profile's content |

**Examples:**

```bash
# Empty profile
profs add-profile personal

# Copy current content
profs add-profile backup --copy-existing
```

### profs remove-profile

Delete a profile.

```bash
profs remove-profile <name>
```

!!! danger
    This permanently deletes all content in the profile's directories.

### profs set

Switch to a profile.

```bash
profs set <profile>
```

Updates all symlinks to point to the specified profile.

### profs list

List all available profiles.

```bash
profs list
```

**Aliases:** `profs list-profiles`

## Status Commands

### profs status

Show current profile and path status.

```bash
profs status
```

Output:
```
Profile: work
~/.gitconfig -> ~/.gitconfig.profs/work [ok]
~/.ssh       -> ~/.ssh.profs/work [ok]
```

### profs status-profile

Show only the current profile name.

```bash
profs status-profile
```

### profs status-full

Show detailed status with all profiles.

```bash
profs status-full
```

### profs status-config

Show raw JSON configuration.

```bash
profs status-config
```

## Diagnostics

### profs doctor

Check for configuration issues.

```bash
profs doctor
```

Finds:
- Broken symlinks
- Missing profile directories
- Inconsistent state across paths

## Maintenance

### profs reset

Remove all profs configuration.

```bash
profs reset
```

!!! warning
    This clears the profs config file. Files remain in `.profs` directories.

### profs migrate-config-dir

Migrate from legacy config location.

```bash
profs migrate-config-dir
```

Moves `~/.profs` to `~/.config/gigurra/profs`.

## Shell Completion

### profs completion

Generate shell completion scripts.

```bash
profs completion bash
profs completion zsh
profs completion fish
profs completion powershell
```

**Setup examples:**

```bash
# Bash
profs completion bash >> ~/.bashrc

# Zsh
profs completion zsh >> ~/.zshrc

# Fish
profs completion fish > ~/.config/fish/completions/profs.fish
```

## Global Flags

| Flag | Description |
|------|-------------|
| `-h, --help` | Help for any command |

## Help

```bash
profs --help
profs <command> --help
```
