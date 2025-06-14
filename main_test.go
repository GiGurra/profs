package main

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/profs/internal"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var testMutex = sync.Mutex{}

// Tests must be executed in a single-threaded manner,
// since we override global os.Args
func runTest(
	t *testing.T,
	argSet [][]string,
	verifier func(t *testing.T, pan any, err error),
) {
	testMutex.Lock()
	defer testMutex.Unlock()
	orgOsArgs := os.Args
	defer func() { os.Args = orgOsArgs }()
	internal.TestMode = true
	defer func() { internal.TestMode = false }()

	// delete config file/state after each test
	// NOTE: TestMode must be set to true, else we will
	// be deleting the config file in the real environment
	defer func() {
		configPath := internal.GlobalConfigPath()
		if _, err := os.Stat(configPath); err == nil {
			if err := os.Remove(configPath); err != nil {
				t.Fatalf("Failed to remove config file: %v", err)
			}
		}
	}()

	var pan any = nil
	var err error = nil
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Test panicked: %v", r)
			}
		}()
		for _, args := range argSet {
			os.Args = args
			mainCmd().RunH(boa.ResultHandler{
				Panic:   func(a any) { pan = a },
				Failure: func(e error) { err = e },
			})
		}
		verifier(t, pan, err)
	}()
}

func checkNoFailures(t *testing.T, pan any, err error) {
	if err != nil {
		t.Fatalf("Expected no error, got panic: %v", err)
	}
	if pan != nil {
		t.Fatalf("Expected no panic, got error: %v", pan)
	}
}

func mkTempDir() string {
	tempDir, err := os.MkdirTemp("", "profs-test")
	if err != nil {
		panic("Failed to create temp dir: " + err.Error())
	}
	return tempDir
}

func mkDir(pathParts ...string) string {
	fullPath := filepath.Join(pathParts...)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		panic("Failed to create dir: " + err.Error())
	}
	return fullPath
}

func deleteDirAndContents(dir string) {
	if dir == "" {
		return
	}
	err := os.RemoveAll(dir)
	if err != nil {
		panic("Failed to delete temp dir: " + err.Error())
	}
}

func TestHelp(t *testing.T) {
	runTest(t, [][]string{{"profs", "--help"}}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)
	})
}

func TestList(t *testing.T) {
	runTest(t, [][]string{{"profs", "list"}}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)
		// configs should be empty
		conf := internal.LoadGlobalConf()
		if conf.Paths != nil && len(conf.Paths) > 0 {
			t.Fatalf("Expected no paths, got: %v", conf.Paths)
		}
	})
}

func TestAdd1(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd := mkDir(testDir, "dir1")
	runTest(t, [][]string{{"profs", "add", dirToAdd, "--profile", "test"}}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)

		// Verify that the path was added
		conf := internal.LoadGlobalConf()
		if conf.Paths == nil || len(conf.Paths) != 1 {
			t.Fatalf("Expected 1 path, got: %v", conf.Paths)
		}

		if conf.Paths[0].SrcPath != dirToAdd {
			t.Fatalf("Expected path to be '%s', got: %s", dirToAdd, conf.Paths[0].SrcPath)
		}
	})
}

func TestAdd2(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd1 := mkDir(testDir, "dir1")
	dirToAdd2 := mkDir(testDir, "dir2")
	runTest(t, [][]string{
		{"profs", "add", dirToAdd1, "--profile", "test"},
		{"profs", "add", dirToAdd2},
		{"profs", "set", "test"},
	}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)

		// Verify that the path was added
		conf := internal.LoadGlobalConf()
		if len(conf.Paths) != 2 {
			t.Fatalf("Expected 2 paths, got: %v", conf.Paths)
		}

		expected := []string{dirToAdd1, dirToAdd2}
		if diff := cmp.Diff(lo.Map(conf.Paths, func(item internal.Path, _ int) string {
			return item.SrcPath
		}), expected); diff != "" {
			t.Fatalf("Paths mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestSetNonExistingProfile(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd1 := mkDir(testDir, "dir1")
	runTest(t, [][]string{
		{"profs", "add", dirToAdd1, "--profile", "test"},
		{"profs", "set", "testx"},
	}, func(t *testing.T, pan any, err error) {
		if pan == nil {
			t.Fatal("Expected a panic when setting a non-existing profile, got none")
		}
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		expMsg := "Profile 'testx' not found"
		errMsg := fmt.Sprintf("%v", pan)

		if !strings.Contains(errMsg, expMsg) {
			t.Fatalf("Expected panic message to contain '%s', got: %s", expMsg, errMsg)
		}
	})
}
