package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/spf13/cobra"
	"os"
)

func ResetCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Yes bool `descr:"Skip confirmation prompt" flag:"yes"  default:"false"`
	}

	return boa.Cmd{
		Use:         "reset",
		Short:       "Resets all configuration to zero",
		Params:      &params,
		ParamEnrich: paramEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			if !params.Yes {
				if !askForConfirmation(
					"Are you sure you want to reset all configurations?\n" +
						"This will reset the configuration to zero, but all symlinks will remain intact",
				) {
					ExitWithMsg(0, "Aborting reset")
				}
			}

			configDir := ConfigDir()
			if _, err := os.Stat(configDir); os.IsNotExist(err) {
				fmt.Println("No configuration found, nothing to reset")
				return
			}

			// delete the config directory
			err := os.RemoveAll(configDir)
			if err != nil {
				ExitWithMsg(1, fmt.Sprintf("Failed to reset configuration: %v", err))
			}

			// re-initialize the config directory
			newConfigDir := ConfigDir()
			fmt.Printf("Re-initializing configuration directory: %s\n", newConfigDir)
			err = os.MkdirAll(newConfigDir, 0755)
			if err != nil {
				ExitWithMsg(1, fmt.Sprintf("Failed to create configuration directory: %v", err))
			}
		},
	}.ToCobra()
}
