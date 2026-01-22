# Working with Profiles

## Creating Profiles

### Initial Profile

Created when you add your first path:

```bash
profs add ~/.gitconfig --profile work
```

### Additional Profiles

```bash
profs add-profile personal
```

This creates empty directories for all managed paths.

### Copy Existing Content

```bash
profs add-profile client-a --copy-existing
```

Copies current profile's content into the new profile.

## Switching Profiles

```bash
profs set personal
```

This updates all symlinks atomically.

### Verify the Switch

```bash
profs status
```

```
Profile: personal
~/.gitconfig -> ~/.gitconfig.profs/personal [ok]
~/.ssh       -> ~/.ssh.profs/personal [ok]
```

## Listing Profiles

```bash
profs list
```

```
work
personal
client-a
```

## Removing Profiles

```bash
profs remove-profile client-a
```

!!! warning
    This deletes all content in that profile's directories. Make sure you've backed up anything important.

## Profile Discovery

profs discovers profiles by scanning `.profs` directories:

```
~/.gitconfig.profs/
├── work
├── personal
└── client-a

~/.ssh.profs/
├── work
└── personal    # client-a missing here!
```

If profiles are inconsistent across paths, `profs doctor` will warn you.

## Current Profile

### Check Current

```bash
profs status-profile
```

### How It's Determined

The current profile is detected by checking where symlinks point:

```bash
ls -la ~/.gitconfig
# ~/.gitconfig -> ~/.gitconfig.profs/work
```

If symlinks point to different profiles, you have an inconsistent state.

## Populating a New Profile

After creating an empty profile:

```bash
profs add-profile personal
profs set personal
```

Now edit files directly:

```bash
# Git config
vim ~/.gitconfig

# Generate new SSH key
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519

# Kubernetes
kubectl config use-context personal-cluster
```

The edits go to the personal profile directories.

## Best Practices

### Consistent Profiles

Keep the same profiles across all paths:

```bash
profs doctor  # Check for inconsistencies
```

### Meaningful Names

Use clear, descriptive names:

- `work` / `personal`
- `client-acme` / `client-globex`
- `prod` / `staging` / `dev`

### Document Your Profiles

Keep notes about what each profile contains:

```bash
cat ~/.config/gigurra/profs/README.md
# work: Corporate identity, VPN keys, prod clusters
# personal: GitHub identity, hobby projects
```

### Backup Before Deleting

Before `profs remove-profile`:

```bash
cp -r ~/.gitconfig.profs/old-profile ~/backup/
```
