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

func TestAddProfile(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd1 := mkDir(testDir, "dir1")
	runTest(t, [][]string{
		{"profs", "add", dirToAdd1, "--profile", "test-profile-1"},
		{"profs", "add-profile", "test-profile-2"},
		{"profs", "add-profile", "test-profile-3", "--copy-existing"},
	}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)

		// Verify that the profiles were added
		conf := internal.LoadGlobalConf()
		if len(conf.DetectedProfileNames()) != 3 {
			t.Fatalf("Expected 2 profiles, got: %d", len(conf.DetectedProfileNames()))
		}

		expectedProfiles := []string{"test-profile-1", "test-profile-2", "test-profile-3"}
		if diff := cmp.Diff(conf.DetectedProfileNames(), expectedProfiles); diff != "" {
			t.Fatalf("Profiles mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestRemoveProfile(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd1 := mkDir(testDir, "dir1")
	runTest(t, [][]string{
		{"profs", "add", dirToAdd1, "--profile", "test-profile-1"},
		{"profs", "add-profile", "test-profile-2"},
		{"profs", "add-profile", "test-profile-3"},
		{"profs", "remove-profile", "test-profile-2", "-y"},
	}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)

		// Verify that the profiles were added
		conf := internal.LoadGlobalConf()
		if len(conf.DetectedProfileNames()) != 2 {
			t.Fatalf("Expected 2 profiles, got: %d", len(conf.DetectedProfileNames()))
		}

		expectedProfiles := []string{"test-profile-1", "test-profile-3"}
		if diff := cmp.Diff(conf.DetectedProfileNames(), expectedProfiles); diff != "" {
			t.Fatalf("Profiles mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestList2(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd1 := mkDir(testDir, "dir1")
	dirToAdd2 := mkDir(testDir, "dir2")
	runTest(t, [][]string{
		{"profs", "add", dirToAdd1, "--profile", "test"},
		{"profs", "add", dirToAdd2},
		{"profs", "list"},
	}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)

		// TODO: Capture stdout to verify output
	})
}

func TestList2b(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd1 := mkDir(testDir, "dir1")
	dirToAdd2 := mkDir(testDir, "dir2")
	runTest(t, [][]string{
		{"profs", "add", dirToAdd1, "--profile", "test"},
		{"profs", "add", dirToAdd2},
		{"profs", "list-profiles"},
	}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)

		// TODO: Capture stdout to verify output
	})
}

func TestCantAddWithoutSpecifyingProfile(t *testing.T) {
	testDir := mkTempDir()
	defer func() { deleteDirAndContents(testDir) }()

	dirToAdd1 := mkDir(testDir, "dir1")
	runTest(t, [][]string{
		{"profs", "add", dirToAdd1},
	}, func(t *testing.T, pan any, err error) {
		expectPanic(t, pan, err, "No active profile found and no profile specified. Don't know how to add path.")
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
		expectPanic(t, pan, err, "Profile 'testx' not found")
	})
}

func expectPanic(t *testing.T, pan any, err error, expectedMsg string) {
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if pan == nil {
		t.Fatal("Expected a panic, got none")
	}

	errMsg := fmt.Sprintf("%v", pan)
	if !strings.Contains(errMsg, expectedMsg) {
		t.Fatalf("Expected panic message to contain '%s', got: %s", expectedMsg, errMsg)
	}
}

func expectError(t *testing.T, pan any, err error, expectedMsg string) {
	if pan != nil {
		t.Fatalf("Expected no panic, got: %v", pan)
	}
	if err == nil {
		t.Fatal("Expected an error, got none")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, expectedMsg) {
		t.Fatalf("Expected error message to contain '%s', got: %s", expectedMsg, errMsg)
	}
}

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
		t.Fatalf("Expected no error, got error: %v", err)
	}
	if pan != nil {
		t.Fatalf("Expected no panic, got panic: %v", pan)
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
