package internal

import (
	"fmt"
	"os"

	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func SetCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Profile string `descr:"The profile to load" positional:"true"`
	}

	return boa.Cmd{
		Use:         "set",
		Short:       "Set current profile",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		ValidArgs:   gc.DetectedProfileNames(),
		InitFuncCtx: func(ctx *boa.HookContext, _ any, _ *cobra.Command) error {
			boa.GetParamT(ctx, &params.Profile).SetAlternatives(gc.DetectedProfileNames())
			return nil
		},
		RunFunc: func(cmd *cobra.Command, args []string) {
			if !lo.Contains(gc.DetectedProfileNames(), params.Profile) {
				fmt.Println("Available profiles:")
				for _, p := range gc.DetectedProfileNames() {
					fmt.Printf("  %v\n", p)
				}
				ExitWithMsg(1, fmt.Sprintf("Profile '%s' not found", params.Profile))
			}

			for _, p := range gc.Paths {
				fmt.Printf("Setting profile %v for path %s\n", params.Profile, p.SrcPath)
				if p.Status != StatusOk && p.Status != StatusErrorTgtNotFound && p.Status != StatusErrorSrcNotFound {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it has status %v\n", params.Profile, p.SrcPath, p.Status)
					continue
				}

				tgt, found := lo.Find(p.DetectedProfs, func(prof DetectedProfile) bool {
					return prof.Name == params.Profile
				})

				if !found {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it is not detected\n", params.Profile, p.SrcPath)
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
					ExitWithMsg(1, fmt.Sprintf("Failed to create symlink for profile %s: %v", params.Profile, err))
				}
			}
		},
	}.ToCobra()
}
