package main

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/profs/internal"
	"os"
	"sync"
	"testing"
)

var testMutex = sync.Mutex{}

// Tests must be executed in a single-threaded manner,
// since we override global os.Args
func runTest(
	t *testing.T,
	args []string,
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

	os.Args = args
	var pan any = nil
	var err error = nil
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Test panicked: %v", r)
			}
		}()
		mainCmd().RunH(boa.ResultHandler{
			Panic:   func(a any) { pan = a },
			Failure: func(e error) { err = e },
		})
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

func TestHelp(t *testing.T) {
	runTest(t, []string{"profs", "--help"}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)
	})
}

func TestList(t *testing.T) {
	runTest(t, []string{"profs", "list"}, func(t *testing.T, pan any, err error) {
		checkNoFailures(t, pan, err)
	})
}
