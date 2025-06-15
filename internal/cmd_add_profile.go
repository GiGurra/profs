package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func AddProfileCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Name         string `positional:"true" description:"Name of the profile to add"`
		CopyExisting bool   `name:"copy-existing" description:"If set, copy existing profiles to the new profile instead of creating an empty one" default:"false"`
	}

	return boa.Cmd{
		Use:         "add-profile",
		Short:       "Adds a new profile to be managed by profs",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			// Current profiles must not collide with the name of the profile to add
			if lo.ContainsBy(gc.DetectedProfileNames(), func(item string) bool {
				return strings.ToLower(item) == strings.ToLower(params.Name)
			}) {
				ExitWithMsg(1, fmt.Sprintf("Profile name '%s' already exists, please choose a different name", params.Name))
			}

			// Go through all the paths and add the profile to each of them
			// First check they are all valid
			for _, path := range gc.Paths {
				if path.Status != StatusOk {
					ExitWithMsg(1, fmt.Sprintf("Cannot add profile '%s' to path '%s' because it is not valid: %s, aborting", params.Name, path.SrcPath, path.Status))
				}
			}

			// Then add the profile to each path
			for _, path := range gc.Paths {
				// create a new empty dir in the companion profs directory
				profsDir, err := path.ProfsDir()
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to get profs directory for path '%s': %v", path.SrcPath, err))
				}

				newProfilePath := filepath.Join(profsDir, params.Name)
				if params.CopyExisting {
					currentProfiles := gc.ActiveProfileNames()
					if len(currentProfiles) == 0 {
						ExitWithMsg(1, fmt.Sprintf("No active profiles found to copy from for path '%s'", path.SrcPath))
					} else if len(currentProfiles) > 1 {
						ExitWithMsg(1, fmt.Sprintf("Multiple active profiles found (%s) for path '%s', please specify which one to copy using --profile", strings.Join(currentProfiles, ", "), path.SrcPath))
					}
					currentProfile := currentProfiles[0]
					currentProfilePath := filepath.Join(profsDir, currentProfile)

					// if it's a file, copy it, if its a dir, use os.CopyFS
					if stat, err := os.Stat(currentProfilePath); err != nil {
						ExitWithMsg(1, fmt.Sprintf("Failed to stat current profile '%s' for path '%s': %v", currentProfile, path.SrcPath, err))
					} else if stat.IsDir() {
						err = os.CopyFS(newProfilePath, os.DirFS(currentProfilePath))
						if err != nil {
							ExitWithMsg(1, fmt.Sprintf("Failed to copy existing profile '%s' to new profile '%s' for path '%s': %v", currentProfile, params.Name, path.SrcPath, err))
						}
					} else {
						fileBytes, err := os.ReadFile(currentProfilePath)
						if err != nil {
							ExitWithMsg(1, fmt.Sprintf("Failed to read current profile file '%s' for path '%s': %v", currentProfilePath, path.SrcPath, err))
						}
						err = os.WriteFile(newProfilePath, fileBytes, 0644)
						if err != nil {
							ExitWithMsg(1, fmt.Sprintf("Failed to write new profile file '%s' for path '%s': %v", newProfilePath, path.SrcPath, err))
						}
					}
				} else {
					err = os.MkdirAll(newProfilePath, 0755)
					if err != nil {
						ExitWithMsg(1, fmt.Sprintf("Failed to create profile directory '%s' for path '%s': %v", newProfilePath, path.SrcPath, err))
					}
				}
			}

		},
	}.ToCobra()
}
