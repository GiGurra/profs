package main

import (
	"encoding/json"
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"path/filepath"
)

type Params struct {
	Profile boa.Required[int] `descr:"The profile to load"`
}

func main() {
	p := Params{}
	boa.Wrap{
		Use:         "profs",
		Long:        "Load user profile",
		Params:      &p,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {
			gc := LoadGlobalConf()
			slog.Info(fmt.Sprintf("Global config: %v", gc))
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

	gc := GlobalConfig{}
	err = json.Unmarshal(bytes, &gc)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse global config from file: %v", err))
	}

	return gc
}

type GlobalConfig struct {
	Paths []string `json:"paths"`
}
