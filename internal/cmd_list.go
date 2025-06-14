package internal

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/spf13/cobra"
)

func ListProfilesCmd(cmdName string, gc GlobalConfig) *cobra.Command {

	var params struct {
	}

	return boa.Cmd{
		Use:         cmdName,
		Short:       "Lists all detected profiles",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			profileNames := gc.DetectedProfileNames()
			if len(profileNames) == 0 {
				println("No profiles detected")
				return
			} else {
				println("Detected profiles:")
				for _, profileName := range profileNames {
					println("  " + profileName)
				}
			}
		},
	}.ToCobra()
}
