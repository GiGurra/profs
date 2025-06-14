package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/spf13/cobra"
	"os"
)

func MigrateConfigDir(gc GlobalConfig) *cobra.Command {

	var params struct {
		Yes bool `descr:"Skip confirmation prompt" flag:"yes"  default:"false"`
	}

	legacyPath := LegacyConfigDirPath()
	newPath := ConfigDirPath()

	return boa.Cmd{
		Use:         "migrate-config-dir",
		Short:       "Migrate legacy " + legacyPath + " -> " + newPath,
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			if !fileOrDirExists(legacyPath) {
				ExitWithMsg(1, "No legacy configuration found at "+legacyPath+", nothing to migrate")
			}

			if fileOrDirExists(newPath) {
				ExitWithMsg(1, "New configuration directory already exists at "+newPath+", please remove it before migrating")
			}

			if !params.Yes {
				if !askForConfirmation(
					"Are you sure you want to migrate the configuration dir?\n" +
						"from legacy " + legacyPath + " -> " + newPath,
				) {
					ExitWithMsg(0, "Aborting migration")
				}
			}

			fmt.Printf("Migrating configuration from %s to %s\n", legacyPath, newPath)
			err := os.Rename(legacyPath, newPath)
			if err != nil {
				ExitWithMsg(1, fmt.Sprintf("Failed to migrate configuration: %v", err))
			}
		},
	}.ToCobra()
}
