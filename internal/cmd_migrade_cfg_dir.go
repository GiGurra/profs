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
				fmt.Println("No legacy configuration found, nothing to migrate")
				os.Exit(1)
			}

			if fileOrDirExists(newPath) {
				fmt.Println("New configuration directory already exists, please remove it before migrating")
				os.Exit(1)
			}

			if !params.Yes {
				if !askForConfirmation(
					"Are you sure you want to migrate the configuration dir?\n" +
						"from legacy " + legacyPath + " -> " + newPath,
				) {
					fmt.Println("Aborting migration")
					os.Exit(0)
				}
			}

			fmt.Printf("Migrating configuration from %s to %s\n", legacyPath, newPath)
			err := os.Rename(legacyPath, newPath)
			if err != nil {
				fmt.Printf("Failed to migrate configuration: %v\n", err)
				os.Exit(1)
			}
		},
	}.ToCobra()
}
