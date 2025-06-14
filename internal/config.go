package internal

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
)

type GlobalConfigRaw struct {
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

func LoadGlobalConfRaw() GlobalConfigRaw {

	bytes, err := os.ReadFile(GlobalConfigPath())
	if err != nil {
		panic(fmt.Sprintf("Failed to read global config file: %v", err))
	}

	gc := GlobalConfigRaw{}
	err = json.Unmarshal(bytes, &gc)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse global config from file: %v", err))
	}

	return gc
}

func SaveGlobalConfRaw(gcr GlobalConfigRaw) {
	jsBytes, err := json.MarshalIndent(gcr, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal global config to json: %v", err))
	}

	err = os.WriteFile(GlobalConfigPath(), jsBytes, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Failed to write global config to file: %v", err))
	}
}

func LoadGlobalConf() GlobalConfig {

	gcr := LoadGlobalConfRaw()

	// convert any paths starting with ~ to absolute paths with home dir
	for i, p := range gcr.Paths {
		if strings.HasPrefix(p, "~") {
			gcr.Paths[i] = filepath.Join(HomeDir(), p[1:])
		}
	}

	return GlobalConfig{
		Paths: lo.Map(gcr.Paths, func(p string, _ int) Path {
			symSrcPath := func() string {
				if strings.HasPrefix(p, "~") {
					return filepath.Join(HomeDir(), p[1:])
				} else {
					return p
				}
			}()
			detectedProfs := profsOnPath(symSrcPath + ".profs")
			var status Status
			var resolvedTgt *DetectedProfile = nil
			var target *string = nil

			if !fileOrDirExists(symSrcPath) {
				status = StatusErrorSrcNotFound
			} else if isSymlink(symSrcPath) {
				symLinKT := symlinkTarget(symSrcPath)
				target = &symLinKT
				if isRelativePath(symLinKT) {
					symLinKT = filepath.Join(filepath.Dir(symSrcPath), symLinKT)
				}

				if !fileOrDirExists(symLinKT) {
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

func LegacyConfigDirPath() string {
	return filepath.Join(HomeDir(), ".profs")
}

func ConfigDirPath() string {
	return filepath.Join(HomeDir(), "/.config/gigurra/profs")
}

func ConfigDir() string {
	path := ConfigDirPath()
	if fileOrDirExists(path) {
		return path
	}

	legacyPath := LegacyConfigDirPath()
	if fileOrDirExists(legacyPath) {
		return legacyPath
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("Failed to read or create config dir: %v", err))
	}
	return path
}

func GlobalConfigPath() string {
	filePath := filepath.Join(ConfigDir(), "global.json")
	if TestMode {
		filePath = filepath.Join(ConfigDir(), "global.test.json")
	}
	// check that file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create a blank config and save it
		blankConfig := GlobalConfigRaw{}
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
