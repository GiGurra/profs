package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func DoctorCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		RepairMode string `descr:"Repair mode, currently just 'all', to be expanded in the future" default:""`
	}

	return boa.Cmd{
		Use:         "doctor",
		Short:       "Show inconsistencies in the current configuration and optional repair them",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			allProfileNames := gc.DetectedProfileNames()
			activeProfileNames := gc.ActiveProfileNames()

			if len(allProfileNames) == 0 {
				ExitWithMsg(1, "No profiles detected, nothing to do. Add at least one profile first")
			}

			if len(activeProfileNames) == 0 {
				// TODO: Implement fix/repair operation
				ExitWithMsg(1, "No active profiles detected, nothing to do. Activate at least one profile first")
			}

			if len(activeProfileNames) != 1 {
				// TODO: Implement fix/repair operation
				ExitWithMsg(1, fmt.Sprintf("Multiple active profiles detected (%d), expected only one active profile. Please deactivate all but one profile", len(activeProfileNames)))
			}

			activeProfileName := activeProfileNames[0]

			for _, path := range gc.Paths {

				fmt.Printf("Checking path: %s\n", path.SrcPath)

				if !isSymlink(path.SrcPath) {
					// TODO: Implement fix/repair operation
					fmt.Printf(" * WARN * Source path %s is not a symlink, expected a symlink to a profile directory.\n", path.SrcPath)
					continue
				}

				activeProfile, err := path.ActiveProfile()
				if err != nil {
					// TODO: Implement fix/repair operation
					fmt.Printf(" * WARN * Error getting active profile for path %s: %v\n", path.SrcPath, err)
				} else if activeProfile != activeProfileName {
					// TODO: Implement fix/repair operation
					fmt.Printf(" * WARN * Active profile for path %s is '%s', expected '%s'.\n", path.SrcPath, activeProfile, activeProfileName)
				}

				profilesForPath := lo.Map(path.DetectedProfs, func(item DetectedProfile, _ int) string {
					return item.Name
				})
				missingProfiles, _ := lo.Difference(allProfileNames, profilesForPath)
				if len(missingProfiles) > 0 {
					// TODO: Implement fix/repair operation
					for _, missingProfile := range missingProfiles {
						fmt.Printf(" * WARN * Profile '%s' is missing for path %s\n", missingProfile, path.SrcPath)
					}
				}
			}

		},
	}.ToCobra()
}
