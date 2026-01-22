# Managing Paths

## Adding Paths

### First Path (Creates Initial Profile)

```bash
profs add ~/.gitconfig --profile work
```

The `--profile` flag is required for the first path to name your initial profile.

### Subsequent Paths

```bash
profs add ~/.ssh
profs add ~/.kube
profs add ~/.config/gcloud
```

No `--profile` needed - paths are added to the current profile.

## What Happens When You Add

1. A `.profs` companion directory is created:
   - `~/.ssh` → `~/.ssh.profs/`

2. The current content is moved into a profile subdirectory:
   - `~/.ssh/*` → `~/.ssh.profs/work/*`

3. A symlink replaces the original:
   - `~/.ssh` → `~/.ssh.profs/work`

## Path Types

### Files

```bash
profs add ~/.gitconfig
```

Structure:
```
~/.gitconfig -> ~/.gitconfig.profs/work
~/.gitconfig.profs/
├── work        # (file)
└── personal    # (file)
```

### Directories

```bash
profs add ~/.ssh
```

Structure:
```
~/.ssh -> ~/.ssh.profs/work
~/.ssh.profs/
├── work/
│   ├── id_ed25519
│   └── config
└── personal/
    ├── id_ed25519
    └── config
```

## Removing Paths

### Stop Managing a Path

```bash
profs remove ~/.kube
```

This:
- Removes the symlink
- Moves the current profile's content back to the original location
- Removes from profs configuration

The `.profs` directory with other profiles remains (for safety).

### Clean Up .profs Directory

Manually delete if you're sure:

```bash
rm -rf ~/.kube.profs
```

## Listing Managed Paths

### Current Status

```bash
profs status
```

```
Profile: work
~/.gitconfig     -> ~/.gitconfig.profs/work [ok]
~/.ssh           -> ~/.ssh.profs/work [ok]
```

### Raw Configuration

```bash
profs status-config
```

Shows the JSON configuration.

## Common Paths to Manage

| Path | Purpose |
|------|---------|
| `~/.gitconfig` | Git identity and settings |
| `~/.ssh` | SSH keys and config |
| `~/.kube` | Kubernetes contexts |
| `~/.config/gcloud` | Google Cloud credentials |
| `~/.aws` | AWS credentials |
| `~/.config/gh` | GitHub CLI auth |
| `~/.npmrc` | npm registry auth |

## Tips

### Add All at Once

Set up a new machine quickly:

```bash
profs add ~/.gitconfig --profile work
profs add ~/.ssh
profs add ~/.kube
profs add ~/.config/gcloud
profs add ~/.aws
```

### Check Before Switching

Always verify the switch will work:

```bash
profs doctor
```

Fix any issues before `profs set`.
