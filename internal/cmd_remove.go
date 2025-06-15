package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func RemoveCmd(cmdName string, gc GlobalConfig) *cobra.Command {

	var params struct {
		Path string `positional:"true" description:"Path to add"`
		Yes  bool   `short:"y" name:"yes" description:"Skip confirmation prompt"`
	}

	rawConfig := LoadGlobalConfRaw()
	alternatives := rawConfig.Paths

	return boa.Cmd{
		Use:         cmdName,
		Short:       "Removes a directory from profs config",
		Long:        "Removes a directory from profs config.\nNOTE: This does not remove symlinks or directories,\nit only removes the path from the profs configuration.",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		ValidArgs:   alternatives,
		RunFunc: func(cmd *cobra.Command, args []string) {

			if !lo.Contains(rawConfig.Paths, params.Path) {
				ExitWithMsg(1, fmt.Sprintf("Path '%s' does not exist in profs configuration", params.Path))
			}

			if !params.Yes {
				if !askForConfirmation("Are you sure you want to remove the path from profs configuration?") {
					ExitWithMsg(0, "Aborting removal of path "+params.Path)
				}
			}

			// Remove the path from the configuration
			rawConfig.Paths = lo.Filter(rawConfig.Paths, func(p string, _ int) bool {
				return p != params.Path
			})

			SaveGlobalConfRaw(rawConfig)

		},
	}.ToCobra()
}
