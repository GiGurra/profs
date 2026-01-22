# Use Cases

## Work vs Personal

The most common scenario: separate work and personal identities.

### Setup

```bash
# Start with work (your current config)
profs add ~/.gitconfig --profile work
profs add ~/.ssh

# Create personal profile
profs add-profile personal
profs set personal

# Configure personal identity
git config --global user.name "Personal Name"
git config --global user.email "personal@example.com"
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "personal@example.com"
```

### Daily Usage

```bash
# Start of work day
profs set work

# After hours
profs set personal
```

## Multiple Clients

Consultants working with multiple clients need separate credentials.

### Setup

```bash
profs add ~/.gitconfig --profile client-a
profs add ~/.ssh
profs add ~/.kube
profs add ~/.config/gcloud
profs add ~/.aws

profs add-profile client-b
profs add-profile client-c
```

### Per-Client Configuration

```bash
profs set client-a
# Configure Client A's credentials...

profs set client-b
# Configure Client B's credentials...
```

### Switching Context

```bash
# Morning: Client A standup
profs set client-a

# Afternoon: Client B deployment
profs set client-b
```

## Multi-Cloud

Different cloud provider configurations.

### Setup

```bash
profs add ~/.aws --profile aws-main
profs add ~/.config/gcloud
profs add ~/.azure

profs add-profile gcp-main
profs add-profile azure-main
```

### Cloud-Specific Work

```bash
profs set aws-main
aws s3 ls

profs set gcp-main
gcloud compute instances list
```

## Dev/Staging/Prod

Separate environments for deployment pipelines.

### Setup

```bash
profs add ~/.kube --profile dev
profs add ~/.config/gcloud

profs add-profile staging
profs add-profile prod
```

### Deployment Workflow

```bash
# Test in dev
profs set dev
kubectl apply -f deployment.yaml

# Promote to staging
profs set staging
kubectl apply -f deployment.yaml

# Production deploy
profs set prod
kubectl apply -f deployment.yaml
```

## GitHub Accounts

Multiple GitHub identities (personal, work, OSS maintainer).

### Setup

```bash
profs add ~/.ssh --profile personal
profs add ~/.config/gh
profs add ~/.gitconfig

profs add-profile work-github
profs add-profile oss-maintainer
```

### Per-Account SSH Keys

Each profile has its own SSH key:

```bash
profs set personal
cat ~/.ssh/id_ed25519.pub
# Add to personal GitHub account

profs set work-github
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519
cat ~/.ssh/id_ed25519.pub
# Add to work GitHub account
```

### SSH Config (Per Profile)

In `~/.ssh/config`:

```
# personal profile
Host github.com
    IdentityFile ~/.ssh/id_ed25519
    User git

# work-github profile
Host github.com
    IdentityFile ~/.ssh/id_ed25519
    User git
```

## Team Onboarding

Quickly set up new team members with standard profiles.

### Create Template

```bash
# Set up standard configs
profs add ~/.gitconfig --profile company
profs add ~/.ssh
profs add ~/.kube
profs add ~/.config/gcloud

# Document setup
cat > ~/.config/gigurra/profs/SETUP.md << 'EOF'
Company Profile Setup:
1. Copy SSH key from 1Password
2. Import kubeconfig from team drive
3. Run `gcloud auth login`
EOF
```

### New Team Member

```bash
# Install profs
go install github.com/GiGurra/profs@latest

# They set up their machine
profs add ~/.gitconfig --profile company
profs add ~/.ssh
profs add ~/.kube
profs add ~/.config/gcloud

# Follow company setup docs
cat ~/.config/gigurra/profs/SETUP.md
```

## CI/CD Profiles

Different credentials for CI environments.

### Setup

```bash
profs add ~/.kube --profile local
profs add ~/.config/gcloud

profs add-profile ci
```

### CI Configuration

The `ci` profile contains service account credentials, not personal credentials.

```bash
profs set ci
# Uses service account for deployments
```
