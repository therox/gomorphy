package gomorphy

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/therox/gomorphy/internal/dict"
)

// dictBinPath is the path to a temporary .bin built in TestMain from the
// mini dictionary in testdata/sample.xml. Used only by Init tests that need
// a real on-disk file to exercise dict.LoadFile.
var dictBinPath string

// childEnv is the environment variable by which a child process learns it
// was launched as a helper for Init/GOMORPHY_DICT tests rather than as an
// ordinary `go test` run.
const childEnv = "GOMORPHY_TEST_CHILD"

// TestMain builds a mini dictionary from testdata/sample.xml, saves it to a
// temporary .bin (for Init tests), and installs the global singleton
// pointer via the private useDict — so the bulk of the tests run without
// file I/O.
//
// Before "claiming" the singleton, the childEnv flag is checked: if the
// process was launched as a child of the subprocess tests, the matching
// mode is executed and os.Exit is called without entering the regular test
// runner.
func TestMain(m *testing.M) {
	if mode := os.Getenv(childEnv); mode != "" {
		runChildMode(mode)
		return // unreachable, runChildMode always os.Exit
	}

	xmlPath := filepath.Join("testdata", "sample.xml")
	d, err := dict.BuildSmall(xmlPath)
	if err != nil {
		panic("gomorphy tests: failed to build mini dictionary: " + err.Error())
	}

	// Save the mini dictionary to .bin: needed for child processes that
	// exercise the public Init/LoadFile path.
	tmp, err := os.CreateTemp("", "gomorphy-test-*.bin")
	if err != nil {
		panic("gomorphy tests: failed to create temp file: " + err.Error())
	}
	if err := d.Encode(tmp); err != nil {
		panic("gomorphy tests: failed to save mini dictionary: " + err.Error())
	}
	if err := tmp.Close(); err != nil {
		panic("gomorphy tests: failed to close temp file: " + err.Error())
	}
	dictBinPath = tmp.Name()
	defer os.Remove(dictBinPath)

	useDict(d)
	os.Exit(m.Run())
}

// runChildMode runs in the child process. Each mode terminates via os.Exit
// with a unique non-zero code on failure, or 0 on success — the parent test
// maps these into assertions.
func runChildMode(mode string) {
	path := os.Getenv("GOMORPHY_TEST_PATH")

	switch mode {
	case "init-ok":
		// Direct Init with an existing .bin: expect success and a working API.
		if err := Init(path); err != nil {
			fmt.Fprintln(os.Stderr, "Init returned error:", err)
			os.Exit(11)
		}
		got, err := Decline("аппетит", Genitive, Singular)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Decline after Init:", err)
			os.Exit(12)
		}
		if got != "аппетита" {
			fmt.Fprintf(os.Stderr, "Decline returned %q, expected \"аппетита\"\n", got)
			os.Exit(13)
		}
		os.Exit(0)

	case "init-twice":
		// A second Init must return the "already initialized" error.
		if err := Init(path); err != nil {
			fmt.Fprintln(os.Stderr, "first Init failed:", err)
			os.Exit(21)
		}
		if err := Init(path); err == nil {
			fmt.Fprintln(os.Stderr, "second Init did not return an error")
			os.Exit(22)
		}
		os.Exit(0)

	case "init-bad":
		// Init with a non-existent file — expect an error.
		if err := Init("/nonexistent/gomorphy-dict.bin"); err == nil {
			fmt.Fprintln(os.Stderr, "Init with a bad path did not return an error")
			os.Exit(31)
		}
		os.Exit(0)

	case "env-ok":
		// GOMORPHY_DICT is set externally; we do not call Init — the autoload
		// must trigger on the first API access.
		got, err := Decline("аппетит", Genitive, Singular)
		if err != nil {
			fmt.Fprintln(os.Stderr, "autoload did not work:", err)
			os.Exit(41)
		}
		if got != "аппетита" {
			fmt.Fprintf(os.Stderr, "Decline returned %q, expected \"аппетита\"\n", got)
			os.Exit(42)
		}
		os.Exit(0)

	case "env-missing":
		// Neither Init nor GOMORPHY_DICT — must yield a clear error.
		os.Unsetenv("GOMORPHY_DICT")
		_, err := Decline("аппетит", Genitive, Singular)
		if err == nil {
			fmt.Fprintln(os.Stderr, "expected \"dictionary not initialized\" error")
			os.Exit(51)
		}
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "unknown helper-process mode: %q\n", mode)
		os.Exit(99)
	}
}

// runChild starts the current test binary as a child process with the
// supplied env variables and the given mode. Returns the exit code together
// with combined stderr+stdout.
func runChild(t *testing.T, mode string, env ...string) (int, string) {
	t.Helper()
	// -test.run=^$ — do not run any ordinary test; we only want to go
	// through TestMain → runChildMode.
	cmd := exec.Command(os.Args[0], "-test.run=^$")
	cmd.Env = append(os.Environ(), childEnv+"="+mode)
	cmd.Env = append(cmd.Env, env...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	if err == nil {
		return 0, buf.String()
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), buf.String()
	}
	t.Fatalf("starting child process (%s): %v", mode, err)
	return -1, ""
}

func TestInitFromValidPath(t *testing.T) {
	ec, out := runChild(t, "init-ok", "GOMORPHY_TEST_PATH="+dictBinPath)
	if ec != 0 {
		t.Fatalf("Init from valid .bin: exit=%d, output=%s", ec, out)
	}
}

func TestInitTwiceFails(t *testing.T) {
	ec, out := runChild(t, "init-twice", "GOMORPHY_TEST_PATH="+dictBinPath)
	if ec != 0 {
		t.Fatalf("repeated Init must return an error: exit=%d, output=%s", ec, out)
	}
}

func TestInitBadPath(t *testing.T) {
	ec, out := runChild(t, "init-bad")
	if ec != 0 {
		t.Fatalf("Init with a bad path must return an error: exit=%d, output=%s", ec, out)
	}
}

func TestAutoloadFromEnv(t *testing.T) {
	ec, out := runChild(t, "env-ok", "GOMORPHY_DICT="+dictBinPath)
	if ec != 0 {
		t.Fatalf("autoload via GOMORPHY_DICT: exit=%d, output=%s", ec, out)
	}
}

func TestNoInitNoEnvErrors(t *testing.T) {
	ec, out := runChild(t, "env-missing")
	if ec != 0 {
		t.Fatalf("expected an error when neither Init nor GOMORPHY_DICT is set: exit=%d, output=%s", ec, out)
	}
}
