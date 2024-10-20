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

type Params struct {
	Profile boa.Required[string] `descr:"The profile to load" positional:"true"`
}

func main() {
	p := Params{}
	boa.Wrap{
		Use:         "profs",
		Long:        "Load user profile",
		Params:      &p,
		ParamEnrich: boa.ParamEnricherDefault,
		ValidArgs:   LoadGlobalConf().ProfileNames(),
		Run: func(cmd *cobra.Command, args []string) {
			gc := LoadGlobalConf()
			fmt.Println(fmt.Sprintf("Global config: %s", PrettyJson(gc)))
			//res := cmder.New("ls", "-la").Run(context.Background())
			//if res.Err != nil {
			//	util.FailAndExit(fmt.Sprintf("Failed to run command: %v", res.Err))
			//}
			//
			//slog.Info(fmt.Sprintf("Result: %v", res.StdOut))
		},
	}.ToApp()
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
				Path:     path,
				Profiles: profsOnPath(profilesPath),
			}
		}),
	}
}

func profsOnPath(path string) []Profile {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn(fmt.Sprintf("Path does not exist: %v", path))
			return []Profile{}
		} else {
			panic(fmt.Sprintf("Failed to stat path: %v", err))
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		panic(fmt.Sprintf("Failed to read dir: %v", err))
	}

	var items []Profile
	for _, f := range files {
		if f.IsDir() || f.Type().IsRegular() || isSymlink(f) {
			fullPath := filepath.Join(path, f.Name())
			items = append(items, Profile{
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

func (g GlobalConfig) ProfileNames() []string {
	return lo.Uniq(lo.FlatMap(g.Paths, func(p Path, _ int) []string {
		return lo.Map(p.Profiles, func(prof Profile, _ int) string {
			return prof.Name
		})
	}))
}

type Path struct {
	Path     string
	Profiles []Profile
}

type Profile struct {
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
