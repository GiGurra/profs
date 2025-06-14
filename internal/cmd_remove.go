package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

func RemoveCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Path string `positional:"true" description:"Path to add"`
		Yes  bool   `short:"y" name:"yes" description:"Skip confirmation prompt"`
	}

	rawConfig := LoadGlobalConfRaw()
	alternatives := rawConfig.Paths

	return boa.Cmd{
		Use:         "remove",
		Short:       "Removes a new directory from profs config",
		Long:        "Removes a new directory from profs config.\nNOTE: This does not remove symlinks or directories,\nit only removes the path from the profs configuration.",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		ValidArgs:   alternatives,
		RunFunc: func(cmd *cobra.Command, args []string) {

			if !lo.Contains(rawConfig.Paths, params.Path) {
				slog.Error("Path does not exist in configuration, aborting", "path", params.Path)
				os.Exit(1)
			}

			if !params.Yes {
				if !askForConfirmation("Are you sure you want to remove the path from profs configuration?") {
					fmt.Printf("Aborting removal of path %s\n", params.Path)
					os.Exit(1)
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
