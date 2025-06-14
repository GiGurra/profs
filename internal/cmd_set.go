package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
)

func SetCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Profile boa.Required[string] `descr:"The profile to load" positional:"true"`
	}

	params.Profile.SetAlternatives(gc.DetectedProfileNames())

	return boa.Cmd{
		Use:         "set",
		Short:       "Set current profile",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		ValidArgs:   params.Profile.GetAlternatives(),
		RunFunc: func(cmd *cobra.Command, args []string) {
			if !lo.Contains(gc.DetectedProfileNames(), params.Profile.Value()) {
				fmt.Println("Available profiles:")
				for _, p := range gc.DetectedProfileNames() {
					fmt.Printf("  %v\n", p)
				}
				ExitWithMsg(1, fmt.Sprintf("Profile '%s' not found", params.Profile.Value()))
			}

			for _, p := range gc.Paths {
				fmt.Printf("Setting profile %v for path %s\n", params.Profile.Value(), p.SrcPath)
				if p.Status != StatusOk && p.Status != StatusErrorTgtNotFound && p.Status != StatusErrorSrcNotFound {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it has status %v\n", params.Profile.Value(), p.SrcPath, p.Status)
					continue
				}

				tgt, found := lo.Find(p.DetectedProfs, func(prof DetectedProfile) bool {
					return prof.Name == params.Profile.Value()
				})

				if !found {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it is not detected\n", params.Profile.Value(), p.SrcPath)
					continue
				}

				if fileOrDirExists(p.SrcPath) && !isSymlink(p.SrcPath) {
					ExitWithMsg(1, fmt.Sprintf("Path %s already exists and is not a symlink, please remove it before setting the profile", p.SrcPath))
				}

				// remove existing symlink
				if fileOrDirExists(p.SrcPath) {
					err := os.Remove(p.SrcPath)
					if err != nil {
						ExitWithMsg(1, fmt.Sprintf("Failed to remove existing symlink: %v", err))
					}
				}

				// create new symlink
				err := os.Symlink(tgt.Path, p.SrcPath)
				if err != nil {
					ExitWithMsg(1, fmt.Sprintf("Failed to create symlink for profile %s: %v", params.Profile.Value(), err))
				}
			}
		},
	}.ToCobra()
}
