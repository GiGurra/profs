package main

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/profs/internal"
	"github.com/spf13/cobra"
)

func main() {
	mainCmd().Run()
}

// For testability
func mainCmd() *boa.Cmd {

	gc := internal.LoadGlobalConf()

	return &boa.Cmd{
		Use:   "profs",
		Short: "Manage user profiles",
		SubCmds: []*cobra.Command{
			internal.MigrateConfigDir(gc),
			internal.AddCmd("add", gc),
			internal.AddCmd("add-path", gc),
			internal.AddProfileCmd(gc),
			internal.DoctorCmd(gc),
			internal.RemoveCmd("remove", gc),
			internal.RemoveCmd("remove-path", gc),
			internal.RemoveProfileCmd(gc),
			internal.ListProfilesCmd("list", gc),
			internal.ListProfilesCmd("list-profiles", gc),
			internal.ResetCmd(gc),
			internal.SetCmd(gc),
			internal.StatusRawCmd(gc),
			internal.StatusCmd(gc),
			internal.StatusProfileCmd(gc),
			internal.FullStatusCmd(gc),
		},
	}
}
