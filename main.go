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
			boa.Wrap{
				Use:  "show",
				Long: "Show current configuration",
				Run: func(cmd *cobra.Command, args []string) {
					gc := LoadGlobalConf()
					fmt.Println(fmt.Sprintf("Global config: %s", PrettyJson(gc)))
				},
			}.ToCmd(),
		},
	}.ToApp()

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
			profilesPath := path + ".profs"
			return Path{
				Path:          path,
				DetectedProfs: profsOnPath(profilesPath),
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
		if f.IsDir() || f.Type().IsRegular() || isSymlink(f) {
			fullPath := filepath.Join(path, f.Name())
			items = append(items, DetectedProfile{
				Name: f.Name(),
				Path: fullPath,
			})
		}
	}

	return items
}

func isSymlink(f os.DirEntry) bool {
	fi, err := f.Info()
	if err != nil {
		panic(fmt.Sprintf("Failed to get file info: %v", err))
	}
	return fi.Mode()&os.ModeSymlink == os.ModeSymlink
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
}

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
