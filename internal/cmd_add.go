package internal

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
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

			rawConfig := LoadGlobalConfRaw()
			if lo.Contains(rawConfig.Paths, params.Path) {
				slog.Error("Path already exists in configuration, aborting", "path", params.Path)
				os.Exit(1)
			}

			// Check that the path exists
			if !fileOrDirExists(params.Path) {
				slog.Warn("Path to add does not exist, creating it", "path", params.Path)
				err := os.MkdirAll(params.Path, 0755)
				if err != nil {
					slog.Error("Failed to create path", "path", params.Path, "error", err)
					os.Exit(1)
				}
			}

			// create a .profs directory if it doesn't exist
			profsDir := params.Path + ".profs"
			err := os.MkdirAll(profsDir, 0755)
			if err != nil {
				slog.Error("Failed to create .profs directory", "path", profsDir, "error", err)
				os.Exit(1)
			}

			// Check if the path is already managed, i.e. is a symlink
			if isSymlink(params.Path) {
				// Check fi it points to a profile in the .profs directory
				target, err := os.Readlink(params.Path)
				if err != nil {
					slog.Error("Failed to read symlink", "link", params.Path, "error", err)
					os.Exit(1)
				}

				parentOfTarget := filepath.Dir(target)
				if parentOfTarget == profsDir {
					slog.Warn("Path is already a symlink managed by profs (=is a symlink), skipping", "path", params.Path)
				} else {
					slog.Error("Path is already a symlink, but doesn't look to be managed by profs, aborting", "path", params.Path)
					os.Exit(1)
				}
			} else {

				// Move the existing directory to .profs, and rename it to the profile name
				newPath := profsDir + "/" + profileName
				err = os.Rename(params.Path, newPath)
				if err != nil {
					slog.Error("Failed to move existing directory to .profs", "oldPath", params.Path, "newPath", newPath, "error", err)
					os.Exit(1)
				}

				// Create a symlink from the new path to the original path
				err = os.Symlink(newPath, params.Path)
				if err != nil {
					slog.Error("Failed to create symlink", "target", newPath, "link", params.Path, "error", err)
					os.Exit(1)
				}

			}

			// Add the new path to the configuration file
			rawConfig.Paths = append(rawConfig.Paths, params.Path)
			SaveGlobalConfRaw(rawConfig)

		},
	}.ToCobra()
}
