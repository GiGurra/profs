# profs

`profs` is a custom CLI directory profiles tool that allows you to manage and switch between different profiles for
specified directories.

The idea is to be able to quickly switch between for example work and personal profiles without necessarily having to
use different computer accounts.

For example, let's say you have different git profiles (email and other diff config). Then you would want to create
different versions of the `~/.gitconfig` file and switch between them. This could be achieved by creating a
`~/.gitconfig.profs/work` and `~/.gitconfig.profs/personal`, and then using `profs` to switch between them.

Let's say you wish to extend this to also have separate ssh configurations. Using `profs`, you would follow the same
pattern and create `~/.ssh.profs/work` and `~/.ssh.profs/personal` directories.

The default `~/.ssh` directory would be a symlink to the active profile directory. To switch all configurations, you
would simply run `profs set work` or `profs set personal`.

WARNING: This is a hack created in about 2 hours. You should probably expect bugs and other issues. There are currently
zero automated tests. Use at your own risk. :)

## Features

- Set and manage profiles for specified directories
- View current configuration and status
- Easy switching between profiles

## Installation

To install `profs`, you need to have Go installed on your system. Then, you can use the following command:

```
go install github.com/GiGurra/profs@latest
```

## Usage

### Set Profile

To set a profile:

```
profs set <profile_name>
```

This command will set the specified profile for all configured paths.

### View Status

To view the current status:

```
profs status
```

This will show the active profile(s) or indicate if no profiles are active.

### View Full Status

To view the full configuration and alternatives:

```
profs status-full
```

This command displays detailed information about all configured paths and detected profiles.

## Configuration

The tool uses a configuration file located at `~/.profs/global.json`. This file should contain a JSON object with a
`paths` array specifying the directories to manage.

Example configuration:

```json
{
  "paths": [
    "~/path/to/manage1",
    "~/path/to/manage2"
  ]
}
```

## How It Works

- The tool manages symlinks for the specified paths.
- Profiles are detected in a `.profs` companion directory of each managed path. For example:
    - `~/.ssh` has a companion directory `~/.ssh.profs`, which contains the profiles, e.g., `~/.ssh.profs/work`.
    - `~/.ssh` should be a symlink pointing to `~/.ssh.profs/work`.
- When setting a profile, the tool updates the symlinks to point to the correct profile directory.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Authors

- Johan Kjölhede
- Claude 3.5
