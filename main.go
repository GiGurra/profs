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
		Use:  "profs",
		Long: "Load user profile",
		SubCommands: []*cobra.Command{
			SetCmd(gc),
			ShowCmd(gc),
			ShowAllCmd(gc),
		},
	}.ToApp()

}

func ShowCmd(gc GlobalConfig) *cobra.Command {
	return boa.Wrap{
		Use:  "show",
		Long: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			gc := LoadGlobalConf()
			fmt.Println(fmt.Sprintf("Global config: %s", PrettyJson(gc)))
		},
	}.ToCmd()
}

func ShowAllCmd(gc GlobalConfig) *cobra.Command {
	return boa.Wrap{
		Use:  "show-all",
		Long: "Show full configuration and alternatives",
		Run: func(cmd *cobra.Command, args []string) {
			gc := LoadGlobalConf()
			fmt.Println(fmt.Sprintf("Global config: %s", PrettyJson(gc)))
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
		Long:        "Set current profile",
		Params:      &params,
		ParamEnrich: boa.ParamEnricherDefault,
		ValidArgs:   params.Profile.GetAlternatives(),
		Run: func(cmd *cobra.Command, args []string) {
			gc := LoadGlobalConf()
			fmt.Println(fmt.Sprintf("Global config: %s", PrettyJson(gc)))
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
			status := StatusErrorTgtNotFound
			var target *DetectedProfile = nil

			if isSymlink(path) {
				symLinKT := symlinkTarget(path)
				_, i, found := lo.FindIndexOf(detectedProfs, func(prof DetectedProfile) bool {
					return pathsAreEqual(prof.Path, symLinKT)
				})

				if found {
					status = StatusOk
					target = &detectedProfs[i]
				}

			}

			return Path{
				Path:          path,
				DetectedProfs: detectedProfs,
				Target:        target,
				Status:        status,
			}
		}),
	}
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

type Path struct {
	Path          string
	DetectedProfs []DetectedProfile
	Target        *DetectedProfile
	Status        Status
}

type Status string

const (
	StatusOk               Status = "ok"
	StatusErrorTgtNotFound Status = "error_tgt_not_found"
	StatusErrorTgtNotProf  Status = "error_tgt_not_prof"
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
