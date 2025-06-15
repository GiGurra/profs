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

func RemoveProfileCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Name string `positional:"true" description:"Name of the profile to add"`
		Yes  bool   `name:"yes" description:"Skip confirmation dialogue" default:"false"`
	}

	existingProfileNames := gc.DetectedProfileNames()

	return boa.Cmd{
		Use:         "remove-profile",
		Short:       "Removes an existing profile managed by profs",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		ValidArgs:   existingProfileNames,
		RunFunc: func(cmd *cobra.Command, args []string) {

			// Current profiles must not collide with the name of the profile to add
			if !lo.ContainsBy(existingProfileNames, func(item string) bool {
				return strings.ToLower(item) == strings.ToLower(params.Name)
			}) {
				ExitWithMsg(1, fmt.Sprintf("Profile name '%s' does not exist", params.Name))
			}

			// Cannot remove an active profile
			if lo.Contains(gc.ActiveProfileNames(), params.Name) {
				ExitWithMsg(1, fmt.Sprintf("Cannot remove active profile '%s'. Please deactivate it first.", params.Name))
			}

			// Ask for confirmation if not skipped
			if !params.Yes {
				if !askForConfirmation(fmt.Sprintf("Are you sure you want to remove the profile '%s'?", params.Name)) {
					ExitWithMsg(0, "Aborting removal of profile "+params.Name)
				}
			}

			// Remove the profile from each path
			for _, path := range gc.Paths {
				profsDir, err := path.ProfsDir()
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to get profs directory for path '%s': %v", path.SrcPath, err))
				}

				profileDir := filepath.Join(profsDir, params.Name)

				if !fileOrDirExists(profileDir) {
					fmt.Printf("Profile '%s' does not exist in path '%s', skipping\n", params.Name, path.SrcPath)
					continue
				}

				err = os.RemoveAll(profileDir)
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to remove profile directory '%s' in path '%s': %v", profileDir, path.SrcPath, err))
				}

				fmt.Printf("Removed profile '%s' from path '%s'\n", params.Name, path.SrcPath)
			}
		},
	}.ToCobra()
}
