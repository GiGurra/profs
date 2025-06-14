package main

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/profs/internal"
	"github.com/spf13/cobra"
)

func main() {

	gc := internal.LoadGlobalConf()

	boa.Cmd{
		Use:   "profs",
		Short: "Load user profile",
		SubCmds: []*cobra.Command{
			internal.MigrateConfigDir(gc),
			internal.AddCmd(gc),
			internal.ListProfilesCmd("list", gc),
			internal.ListProfilesCmd("list-profiles", gc),
			internal.ResetCmd(gc),
			internal.SetCmd(gc),
			internal.StatusRawCmd(gc),
			internal.StatusCmd(gc),
			internal.StatusProfileCmd(gc),
			internal.FullStatusCmd(gc),
		},
	}.Run()
}
