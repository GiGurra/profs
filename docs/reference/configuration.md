# Configuration

## Config File Location

```
~/.config/gigurra/profs/global.json
```

### Legacy Location

Older versions used `~/.profs/`. Migrate with:

```bash
profs migrate-config-dir
```

## Config Format

```json
{
  "paths": [
    "~/.gitconfig",
    "~/.ssh",
    "~/.kube",
    "~/.config/gcloud"
  ]
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `paths` | `[]string` | List of managed paths |

## File Structure

After adding paths:

```
~/.config/gigurra/profs/
└── global.json

~/.gitconfig -> ~/.gitconfig.profs/work
~/.gitconfig.profs/
├── work
└── personal

~/.ssh -> ~/.ssh.profs/work
~/.ssh.profs/
├── work/
│   ├── config
│   ├── id_ed25519
│   └── id_ed25519.pub
└── personal/
    ├── config
    ├── id_ed25519
    └── id_ed25519.pub
```

## Profile Storage

Profiles are stored as subdirectories in `.profs` companion directories:

```
<original-path>.profs/
├── profile-a/
├── profile-b/
└── profile-c/
```

There's no central profile registry - profiles are discovered by scanning these directories.

## Current Profile Detection

The current profile is determined by following symlinks:

```bash
readlink ~/.gitconfig
# /home/user/.gitconfig.profs/work
```

If different paths point to different profiles, you have an inconsistent state.

## Manual Editing

You can edit `global.json` directly:

```bash
vim ~/.config/gigurra/profs/global.json
```

After editing, verify with:

```bash
profs doctor
```

## Backup

To backup your profs setup:

```bash
# Config
cp ~/.config/gigurra/profs/global.json ~/backup/

# All profiles
for path in $(cat ~/.config/gigurra/profs/global.json | jq -r '.paths[]'); do
  cp -r "${path}.profs" ~/backup/
done
```

## Environment

profs uses standard Go path expansion:

| Syntax | Expands To |
|--------|------------|
| `~` | Home directory |
| `$HOME` | Home directory |
| `$VAR` | Environment variable |

## Troubleshooting

### Config Not Found

If profs can't find config:

```bash
# Check location
ls -la ~/.config/gigurra/profs/

# Create manually if needed
mkdir -p ~/.config/gigurra/profs
echo '{"paths":[]}' > ~/.config/gigurra/profs/global.json
```

### Corrupted Config

Reset and start over:

```bash
profs reset
profs add ~/.gitconfig --profile work
# Re-add other paths...
```

### Symlink Issues

Check symlinks manually:

```bash
ls -la ~/.gitconfig
ls -la ~/.ssh
```

Fix broken symlinks:

```bash
# Remove broken link
rm ~/.gitconfig

# Re-create
ln -s ~/.gitconfig.profs/work ~/.gitconfig
```

Or let profs fix it:

```bash
profs set work  # Re-creates all symlinks
```
