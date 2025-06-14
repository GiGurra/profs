package internal

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"strings"
)

func StatusCmd(gc GlobalConfig) *cobra.Command {
	return boa.Cmd{
		Use:   "status",
		Short: "Show current configuration",
		RunFunc: func(cmd *cobra.Command, args []string) {
			profileOf := func(p Path) string {
				if p.ResolvedTgt != nil {
					return p.ResolvedTgt.Name
				} else {
					return ""
				}
			}

			grouped := lo.GroupBy(gc.Paths, profileOf)
			for profile, paths := range grouped {

				longestSrc := 0

				// for padding
				for _, p := range paths {
					srcLen := len(simplifyPath(p.SrcPath))
					if srcLen > longestSrc {
						longestSrc = srcLen
					}
				}

				fmt.Println("Profile: " + profile)
				for _, p := range paths {

					src := simplifyPath(p.SrcPath)
					if len(src) < longestSrc {
						src += strings.Repeat(" ", longestSrc-len(src))
					}

					infoStr := ""
					if p.TgtPath != nil {
						infoStr = infoStr + simplifyPath(*p.TgtPath) + " [" + string(p.Status) + "]"
					}
					fmt.Println(fmt.Sprintf("  %v -> %v", src, infoStr))
				}
			}
		},
	}.ToCobra()
}

func StatusProfileCmd(gc GlobalConfig) *cobra.Command {
	return boa.Cmd{
		Use:   "status-profile",
		Short: "Show current configuration",
		RunFunc: func(cmd *cobra.Command, args []string) {
			profileNames := gc.ActiveProfileNames()
			if len(gc.ActiveProfileNames()) == 0 {
				fmt.Println("No active profiles")
			} else if len(gc.ActiveProfileNames()) == 1 {
				fmt.Println(profileNames[0])
				if !gc.AllProfilesResolved() {
					fmt.Println("WARNING: Not all configured profile resolved!")
					fmt.Println(" -> Run 'profs status-all' to see full configuration")
				}
			} else {
				fmt.Println("WARNING: Multiple active profiles:")
				for _, p := range profileNames {
					fmt.Println(fmt.Sprintf("  %v", p))
				}
				fmt.Println(" -> Run 'profs show-all' to see full configuration")
			}
		},
	}.ToCobra()
}

func FullStatusCmd(gc GlobalConfig) *cobra.Command {
	return boa.Cmd{
		Use:   "status-full",
		Short: "Show full configuration and alternatives",
		RunFunc: func(cmd *cobra.Command, args []string) {
			fmt.Println(PrettyJson(gc))
		},
	}.ToCobra()
}
