package main

import (
	"encoding/json"
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	gc := LoadGlobalConf()

	boa.Cmd{
		Use:   "profs",
		Short: "Load user profile",
		SubCmds: []*cobra.Command{
			SetCmd(gc),
			StatusCmd(gc),
			StatusProfileCmd(gc),
			FullStatusCmd(gc),
		},
	}.Run()
}

func simplifyPath(in string) string {
	homeDir := HomeDir()
	if strings.HasPrefix(in, homeDir) {
		return "~" + in[len(homeDir):]
	} else {
		return in
	}
}

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

func SetCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Profile boa.Required[string] `descr:"The profile to load" positional:"true"`
	}

	params.Profile.SetAlternatives(gc.DetectedProfileNames())

	return boa.Cmd{
		Use:         "set",
		Short:       "Set current profile",
		Params:      &params,
		ParamEnrich: boa.ParamEnricherDefault,
		ValidArgs:   params.Profile.GetAlternatives(),
		RunFunc: func(cmd *cobra.Command, args []string) {
			if !lo.Contains(gc.DetectedProfileNames(), params.Profile.Value()) {
				fmt.Println(fmt.Sprintf("Profile not found: %v", params.Profile.Value()))
				fmt.Println("Available profiles:")
				for _, p := range gc.DetectedProfileNames() {
					fmt.Println(fmt.Sprintf("  %v", p))
				}
				os.Exit(1)
			}

			for _, p := range gc.Paths {
				fmt.Printf("Setting profile %v for path %s\n", params.Profile.Value(), p.SrcPath)
				if p.Status != StatusOk && p.Status != StatusErrorTgtNotFound && p.Status != StatusErrorSrcNotFound {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it has status %v\n", params.Profile.Value(), p.SrcPath, p.Status)
					continue
				}

				tgt, found := lo.Find(p.DetectedProfs, func(prof DetectedProfile) bool {
					return prof.Name == params.Profile.Value()
				})

				if !found {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it is not detected\n", params.Profile.Value(), p.SrcPath)
					continue
				}

				if fileExists(p.SrcPath) && !isSymlink(p.SrcPath) {
					panic(fmt.Sprintf("SrcPath is not a symlink: %v", p.SrcPath))
				}

				// remove existing symlink
				if fileExists(p.SrcPath) {
					err := os.Remove(p.SrcPath)
					if err != nil {
						panic(fmt.Sprintf("Failed to remove existing symlink: %v", err))
					}
				}

				// create new symlink
				err := os.Symlink(tgt.Path, p.SrcPath)
				if err != nil {
					panic(fmt.Sprintf("Failed to create symlink: %v", err))
				}
			}
		},
	}.ToCobra()
}

func HomeDir() string {
	hd, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get home dir: %v", err))
	}

	return hd
}

func ConfigDir() string {
	path := filepath.Join(HomeDir(), ".profs")
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Failed to read or create config dir: %v", err))
	}
	return path
}

func GlobalConfigPath() string {
	filePath := filepath.Join(ConfigDir(), "global.json")
	// check that file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create a blank config and save it
		blankConfig := GlobalConfigStored{}
		jsBytes, err := json.MarshalIndent(blankConfig, "", "  ")
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal blank config to json: %v", err))
		}
		err = os.WriteFile(filePath, jsBytes, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Failed to write blank config to file: %v", err))
		}
	}
	return filePath
}

func LoadGlobalConf() GlobalConfig {

	bytes, err := os.ReadFile(GlobalConfigPath())
	if err != nil {
		panic(fmt.Sprintf("Failed to read global config file: %v", err))
	}

	gc := GlobalConfigStored{}
	err = json.Unmarshal(bytes, &gc)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse global config from file: %v", err))
	}

	// convert any paths starting with ~ to absolute paths with home dir
	for i, p := range gc.Paths {
		if strings.HasPrefix(p, "~") {
			gc.Paths[i] = filepath.Join(HomeDir(), p[1:])
		}
	}

	return GlobalConfig{
		Paths: lo.Map(gc.Paths, func(p string, _ int) Path {
			symSrcPath := func() string {
				if strings.HasPrefix(p, "~") {
					return filepath.Join(HomeDir(), p[1:])
				} else {
					return p
				}
			}()
			detectedProfs := profsOnPath(symSrcPath + ".profs")
			status := StatusErrorSrcNotFound
			var resolvedTgt *DetectedProfile = nil
			var target *string = nil

			if !fileExists(symSrcPath) {
				status = StatusErrorSrcNotFound
			} else if isSymlink(symSrcPath) {
				symLinKT := symlinkTarget(symSrcPath)
				target = &symLinKT
				if isRelativePath(symLinKT) {
					symLinKT = filepath.Join(filepath.Dir(symSrcPath), symLinKT)
				}

				if !fileExists(symLinKT) {
					status = StatusErrorTgtNotFound
				} else {

					_, i, found := lo.FindIndexOf(detectedProfs, func(prof DetectedProfile) bool {
						return pathsAreEqual(prof.Path, symLinKT)
					})

					if found {
						status = StatusOk
						resolvedTgt = &detectedProfs[i]
					} else {
						status = StatusErrorTgtUnresolvable
					}
				}
			} else {
				status = StatusErrorSrcNotSymlink
			}

			return Path{
				SrcPath:       symSrcPath,
				DetectedProfs: detectedProfs,
				TgtPath:       target,
				ResolvedTgt:   resolvedTgt,
				Status:        status,
			}
		}),
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(fmt.Sprintf("Failed to stat path: %v", err))
		}
	}

	return true
}

func isRelativePath(path string) bool {
	return !filepath.IsAbs(path)
}

func profsOnPath(path string) []DetectedProfile {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn(fmt.Sprintf("SrcPath does not exist: %v", path))
			return []DetectedProfile{}
		} else {
			panic(fmt.Sprintf("Failed to stat path: %v", err))
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read dir: %v", err))
	}

	var items []DetectedProfile
	for _, f := range files {
		if f.IsDir() || f.Type().IsRegular() || isSymlinkE(f) {
			fullPath := filepath.Join(path, f.Name())
			items = append(items, DetectedProfile{
				Name: f.Name(),
				Path: fullPath,
			})
		}
	}

	return items
}

func pathsAreEqual(p1, p2 string) bool {
	return filepath.Clean(p1) == filepath.Clean(p2)
}

func isSymlink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn(fmt.Sprintf("SrcPath does not exist: %v", path))
			return false
		}
		panic(fmt.Sprintf("Failed to get file info: %v", err))
	}

	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
}

func isSymlinkE(f os.DirEntry) bool {
	fi, err := f.Info()
	if err != nil {
		panic(fmt.Sprintf("Failed to get file info: %v", err))
	}
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
}

func symlinkTarget(path string) string {

	fi, err := os.Lstat(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to get file info: %v", err))
	}

	if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
		panic(fmt.Sprintf("Not a symlink: %v", path))
	}

	target, err := os.Readlink(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read symlink target: %v", err))
	}

	return target
}

type GlobalConfigStored struct {
	Paths []string `json:"paths"`
}

type GlobalConfig struct {
	Paths []Path
}

func (g GlobalConfig) DetectedProfileNames() []string {
	return lo.Uniq(lo.FlatMap(g.Paths, func(p Path, _ int) []string {
		return lo.Map(p.DetectedProfs, func(prof DetectedProfile, _ int) string {
			return prof.Name
		})
	}))
}

func (g GlobalConfig) ActiveProfileNames() []string {
	return lo.Uniq(lo.FlatMap(g.Paths, func(p Path, _ int) []string {
		if p.ResolvedTgt != nil {
			return []string{p.ResolvedTgt.Name}
		} else {
			return []string{}
		}
	}))
}

func (g GlobalConfig) AllProfilesResolved() bool {
	return lo.EveryBy(g.Paths, func(p Path) bool {
		return p.Status == StatusOk
	})
}

type Path struct {
	SrcPath       string            `json:"srcPath"`
	Status        Status            `json:"status"`
	TgtPath       *string           `json:"tgtPath"`
	ResolvedTgt   *DetectedProfile  `json:"resolvedTgt"`
	DetectedProfs []DetectedProfile `json:"detectedProfs"`
}

type Status string

const (
	StatusOk                   Status = "ok"
	StatusErrorSrcNotFound     Status = "error_src_not_found"
	StatusErrorTgtNotFound     Status = "error_tgt_not_found"
	StatusErrorSrcNotSymlink   Status = "error_tgt_not_prof"
	StatusErrorTgtUnresolvable Status = "error_tgt_not_resolvable"
)

type DetectedProfile struct {
	Name string
	Path string
}

func PrettyJson[T any](t T) string {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal to json: %v", err))
	}

	return string(bytes)
}
