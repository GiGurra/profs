package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func AddCmd(cmdName string, gc GlobalConfig) *cobra.Command {

	var params struct {
		Path    string               `positional:"true" description:"Path to add"`
		Profile boa.Optional[string] `short:"p" name:"profile" description:"Profile to use (defaults to active profile)"`
	}

	return boa.Cmd{
		Use:         cmdName,
		Short:       "Adds a new directory to be managed by profs",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			profileName := params.Profile.GetOrElse("")
			if profileName == "" {
				activeProfiles := gc.ActiveProfileNames()
				switch len(activeProfiles) {
				case 0:
					ExitWithMsg(1, "No active profile found and no profile specified. Don't know how to add path.")
				case 1:
					profileName = activeProfiles[0]
				default:
					ExitWithMsg(1, "Multiple active profiles found, please specify which profile to use with --profile")
				}
			}

			// check if it under the current user's home directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				ExitWithMsg(1, fmt.Sprintf("Failed to get user home directory, error: %v", err))
			}

			// check that the path is absolute
			path := params.Path
			if !filepath.IsAbs(path) {
				path, err = filepath.Abs(path)
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to get absolute path for '%s', error: %v", path, err))
				}
			}

			rawConfig := LoadGlobalConfRaw()
			if lo.ContainsBy(rawConfig.Paths, func(item string) bool {

				a, errA := filepath.Abs(item)
				if errA != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to get absolute path for '%s', error: %v", item, errA))
				}

				b, errB := filepath.Abs(path)
				if errB != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to get absolute path for '%s', error: %v", path, errB))
				}

				return a == b
			}) {
				ExitWithMsg(1, fmt.Sprintf("Path '%s' already exists in configuration, aborting", path))
			}

			// create a .profs directory if it doesn't exist
			profsDir := path + ".profs"
			err = os.MkdirAll(profsDir, 0755)
			if err != nil {
				ExitWithMsg(1, fmt.Sprintf("Failed to create .profs directory at '%s', error: %v", profsDir, err))
			}

			// Check if the path is already managed, i.e. is a symlink
			if isSymlink(path) {
				// Check fi it points to a profile in the .profs directory
				target, err := os.Readlink(path)
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to read symlink at '%s', error: %v", path, err))
				}

				parentOfTarget, err := filepath.Abs(filepath.Dir(target))
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to get parent directory of symlink target '%s', error: %v", target, err))
				}
				if parentOfTarget == profsDir {
					slog.Warn("Path is already a symlink managed by profs (=is a symlink), skipping", "path", path)
				} else {
					ExitWithMsg(1, fmt.Sprintf(
						"Path is already a symlink, but does not point to a profile in .profs, aborting\n"+
							"  Path: %s\n"+
							"  Parent of target: %s\n"+
							"  Expected parent: %s",
						path,
						parentOfTarget,
						profsDir,
					))
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
							ExitWithMsg(1, fmt.Sprintf("Failed to create path '%s', error: %v", path, err))
						}
					}

					err = os.Rename(path, newPath)
					if err != nil {
						ExitWithMsg(1, fmt.Sprintf("Failed to move existing directory to .profs, error: %v", err))
					}
				} else {
					slog.Warn("Path already exists in .profs, skipping move", "path", newPath)
				}

				// Create a symlink from the new path to the original path
				err = os.Symlink(newPath, path)
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to create symlink from '%s' to '%s', error: %v", newPath, path, err))
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
