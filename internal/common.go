package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func PrettyJson[T any](t T) string {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal to json: %v", err))
	}

	return string(bytes)
}

func fileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
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

func HomeDir() string {
	hd, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get home dir: %v", err))
	}

	return hd
}

func simplifyPath(in string) string {
	homeDir := HomeDir()
	if strings.HasPrefix(in, homeDir) {
		return "~" + in[len(homeDir):]
	} else {
		return in
	}
}

func askForConfirmation(prompt string) bool {
	var response string
	fmt.Printf("%s (y/n): ", prompt)
	_, err := fmt.Scanln(&response)
	if err != nil {
		slog.Error("Failed to read input", "error", err)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return strings.HasPrefix(response, "y")
}

var paramEnricherDefault = boa.ParamEnricherCombine(
	boa.ParamEnricherName,
	boa.ParamEnricherShort,
	//ParamEnricherEnv,
	boa.ParamEnricherBool,
)
