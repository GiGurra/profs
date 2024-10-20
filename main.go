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

	boa.Wrap{
		Use:   "profs",
		Short: "Load user profile",
		SubCommands: []*cobra.Command{
			SetCmd(gc),
			ShowCmd(gc),
			ShowAllCmd(gc),
		},
	}.ToApp()

}

func ShowCmd(gc GlobalConfig) *cobra.Command {
	return boa.Wrap{
		Use:   "show",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			profileNames := gc.ActiveProfileNames()
			if len(gc.ActiveProfileNames()) == 0 {
				fmt.Println("No active profiles")
			} else if len(gc.ActiveProfileNames()) == 1 {
				fmt.Println(profileNames[0])
				if !gc.AllProfilesResolved() {
					fmt.Println("WARNING: Not all configured profile resolved!")
					fmt.Println(" -> Run 'profs show-all' to see full configuration")
				}
			} else {
				fmt.Println("WARNING: Multiple active profiles:")
				for _, p := range profileNames {
					fmt.Println(fmt.Sprintf("  %v", p))
				}
				fmt.Println(" -> Run 'profs show-all' to see full configuration")
			}
		},
	}.ToCmd()
}

func ShowAllCmd(gc GlobalConfig) *cobra.Command {
	return boa.Wrap{
		Use:   "show-all",
		Short: "Show full configuration and alternatives",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("Full config: %s", PrettyJson(gc)))
		},
	}.ToCmd()
}

func SetCmd(gc GlobalConfig) *cobra.Command {

	var params struct {
		Profile boa.Required[string] `descr:"The profile to load" positional:"true"`
	}

	params.Profile.SetAlternatives(gc.DetectedProfileNames())

	return boa.Wrap{
		Use:         "set",
		Short:       "Set current profile",
		Params:      &params,
		ParamEnrich: boa.ParamEnricherDefault,
		ValidArgs:   params.Profile.GetAlternatives(),
		Run: func(cmd *cobra.Command, args []string) {
			if !lo.Contains(gc.DetectedProfileNames(), params.Profile.Value()) {
				fmt.Println(fmt.Sprintf("Profile not found: %v", params.Profile.Value()))
				fmt.Println("Available profiles:")
				for _, p := range gc.DetectedProfileNames() {
					fmt.Println(fmt.Sprintf("  %v", p))
				}
				os.Exit(1)
			}

			for _, p := range gc.Paths {
				fmt.Printf("Setting profile %v for path %s\n", params.Profile.Value(), p.Path)
				if p.Status != StatusOk && p.Status != StatusErrorTgtNotFound {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it has status %v\n", params.Profile.Value(), p.Path, p.Status)
					continue
				}

				tgt, found := lo.Find(p.DetectedProfs, func(prof DetectedProfile) bool {
					return prof.Name == params.Profile.Value()
				})

				if !found {
					fmt.Printf("WARNING: Unable to set profile %v for path %s, because it is not detected\n", params.Profile.Value(), p.Path)
					continue
				}

				// remove existing symlink
				err := os.Remove(p.Path)
				if err != nil {
					panic(fmt.Sprintf("Failed to remove existing symlink: %v", err))
				}

				// create new symlink
				err = os.Symlink(tgt.Path, p.Path)
				if err != nil {
					panic(fmt.Sprintf("Failed to create symlink: %v", err))
				}
			}
		},
	}.ToCmd()
}

func HomeDir() string {
	hd, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get home dir: %v", err))
	}

	return hd
}

func ConfigDir() string {
	return filepath.Join(HomeDir(), ".profs")
}

func GlobalConfigPath() string {
	return filepath.Join(ConfigDir(), "global.json")
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
			path := func() string {
				if strings.HasPrefix(p, "~") {
					return filepath.Join(HomeDir(), p[1:])
				} else {
					return p
				}
			}()
			detectedProfs := profsOnPath(path + ".profs")
			status := StatusErrorSrcNotFound
			var resolvedTgt *DetectedProfile = nil
			var target *string = nil

			if !fileExists(path) {
				status = StatusErrorSrcNotFound
			} else if isSymlink(path) {
				symLinKT := symlinkTarget(path)
				target = &symLinKT
				if isRelativePath(symLinKT) {
					symLinKT = filepath.Join(filepath.Dir(path), symLinKT)
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
				Path:          path,
				DetectedProfs: detectedProfs,
				Tgt:           target,
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
			slog.Warn(fmt.Sprintf("Path does not exist: %v", path))
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
			slog.Warn(fmt.Sprintf("Path does not exist: %v", path))
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
	Path          string            `json:"path"`
	Status        Status            `json:"status"`
	Tgt           *string           `json:"tgt"`
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
