Here's a GitHub README for the `profs` tool based on the provided information:

# profs

`profs` is a custom CLI directory profiles tool that allows you to manage and switch between different profiles for
specified directories.

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

---

This README provides a basic overview of the `profs` tool, its usage, and configuration. You may want to expand on
certain sections, add examples, or include more detailed instructions depending on the complexity and specific use cases
of your tool.
