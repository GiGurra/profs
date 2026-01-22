# Getting Started

## Installation

### Using Go

```bash
go install github.com/GiGurra/profs@latest
```

### Using mise

```bash
mise use -g "go:github.com/GiGurra/profs@latest"
```

### Verify Installation

```bash
profs --help
```

## First-Time Setup

### 1. Add Your First Path

Pick a configuration file you want to manage:

```bash
profs add ~/.gitconfig --profile work
```

This:
- Creates `~/.gitconfig.profs/work` directory
- Moves your current `.gitconfig` into it
- Creates a symlink `~/.gitconfig -> ~/.gitconfig.profs/work`

### 2. Add More Paths

```bash
profs add ~/.ssh
profs add ~/.kube
profs add ~/.config/gcloud
```

Since you already have a profile (`work`), these are added to it automatically.

### 3. Create Another Profile

```bash
profs add-profile personal
```

This creates empty `personal` directories for all managed paths:
- `~/.gitconfig.profs/personal`
- `~/.ssh.profs/personal`
- `~/.kube.profs/personal`
- `~/.config/gcloud.profs/personal`

### 4. Populate the New Profile

You can either:

**Copy from existing profile:**
```bash
profs add-profile personal --copy-existing
```

**Or manually set up:**
```bash
# Switch to personal
profs set personal

# Edit configs directly (they're now pointing to personal profile)
vim ~/.gitconfig
ssh-keygen -f ~/.ssh/id_ed25519
```

### 5. Switch Between Profiles

```bash
profs set work
profs set personal
```

### 6. Check Status

```bash
profs status
```

Output:
```
Profile: work
~/.gitconfig     -> ~/.gitconfig.profs/work [ok]
~/.ssh           -> ~/.ssh.profs/work [ok]
~/.kube          -> ~/.kube.profs/work [ok]
~/.config/gcloud -> ~/.config/gcloud.profs/work [ok]
```

## Shell Completion

Enable tab completion for profile names:

**Bash:**
```bash
profs completion bash >> ~/.bashrc
```

**Zsh:**
```bash
profs completion zsh >> ~/.zshrc
```

**Fish:**
```bash
profs completion fish > ~/.config/fish/completions/profs.fish
```

## Troubleshooting

### Check for Issues

```bash
profs doctor
```

This finds:
- Broken symlinks
- Missing profile directories
- Inconsistent state

### View Full Status

```bash
profs status-full
```

Shows all profiles and their paths.

### Reset Everything

If things go wrong:

```bash
profs reset
```

!!! warning
    This removes all profs configuration. Your actual files remain in the `.profs` directories.
