package internal

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func AddCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Path    string               `positional:"true" description:"Path to add"`
		Profile boa.Optional[string] `short:"p" name:"profile" description:"Profile to use (defaults to active profile)"`
	}

	return boa.Cmd{
		Use:         "add",
		Short:       "Adds a new directory to be managed by profs",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			profileName := params.Profile.GetOrElse("")
			if profileName == "" {
				activeProfiles := gc.ActiveProfileNames()
				switch len(activeProfiles) {
				case 0:
					println("No active profile found and no profile specified. Don't know how to add path.")
					os.Exit(1)
				case 1:
					profileName = activeProfiles[0]
				default:
					println("Multiple active profiles found, please specify which profile to use with --profile")
					os.Exit(1)
				}
			}

			// check if it under the current user's home directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				slog.Error("Failed to get user home directory", "error", err)
				os.Exit(1)
			}

			// check that the path is absolute
			path := params.Path
			if !filepath.IsAbs(path) {
				path, err = filepath.Abs(path)
				if err != nil {
					slog.Error("Failed to get absolute path", "path", path, "error", err)
					os.Exit(1)
				}
			}

			rawConfig := LoadGlobalConfRaw()
			if lo.ContainsBy(rawConfig.Paths, func(item string) bool {
				a, err := filepath.Abs(item)
				if err != nil {
					slog.Error("Failed to get absolute path", "path", item, "error", err)
					os.Exit(1)
				}
				b, err := filepath.Abs(path)
				if err != nil {
					slog.Error("Failed to get absolute path", "path", path, "error", err)
					os.Exit(1)
				}
				return a == b
			}) {
				slog.Error("Path already exists in configuration, aborting", "path", path)
				os.Exit(1)
			}

			// create a .profs directory if it doesn't exist
			profsDir := path + ".profs"
			err = os.MkdirAll(profsDir, 0755)
			if err != nil {
				slog.Error("Failed to create .profs directory", "path", profsDir, "error", err)
				os.Exit(1)
			}

			// Check if the path is already managed, i.e. is a symlink
			if isSymlink(path) {
				// Check fi it points to a profile in the .profs directory
				target, err := os.Readlink(path)
				if err != nil {
					slog.Error("Failed to read symlink", "link", path, "error", err)
					os.Exit(1)
				}

				parentOfTarget, err := filepath.Abs(filepath.Dir(target))
				if err != nil {
					slog.Error("Failed to get parent directory of symlink target", "target", target, "error", err)
					os.Exit(1)
				}
				if parentOfTarget == profsDir {
					slog.Warn("Path is already a symlink managed by profs (=is a symlink), skipping", "path", path)
				} else {
					slog.Error("Path is already a symlink, but doesn't look to be managed by profs, aborting",
						"path", path,
						"parentOfTarget", parentOfTarget,
						"expectedParent", profsDir,
					)
					os.Exit(1)
				}
			} else {

				// Move the existing directory to .profs, and rename it to the profile name
				newPath := profsDir + "/" + profileName
				if !fileOrDirExists(newPath) {

					// Check that the path exists
					if !fileOrDirExists(path) {
						slog.Warn("Path to add does not exist, creating it", "path", path)
						err := os.MkdirAll(path, 0755)
						if err != nil {
							slog.Error("Failed to create path", "path", path, "error", err)
							os.Exit(1)
						}
					}

					err = os.Rename(path, newPath)
					if err != nil {
						slog.Error("Failed to move existing directory to .profs", "oldPath", path, "newPath", newPath, "error", err)
						os.Exit(1)
					}
				} else {
					slog.Warn("Path already exists in .profs, skipping move", "path", newPath)
				}

				// Create a symlink from the new path to the original path
				err = os.Symlink(newPath, path)
				if err != nil {
					slog.Error("Failed to create symlink", "target", newPath, "link", path, "error", err)
					os.Exit(1)
				}

			}

			// Add the new path to the configuration file
			storedPath := path
			if strings.HasPrefix(storedPath, homeDir) {
				storedPath = strings.ReplaceAll(storedPath, homeDir, "~")
			}
			rawConfig.Paths = append(rawConfig.Paths, storedPath)
			SaveGlobalConfRaw(rawConfig)

		},
	}.ToCobra()
}
