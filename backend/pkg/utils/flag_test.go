package flag

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestRequire_AllFlagsPresent(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	name := flag.String("name", "test", "name flag")
	env := flag.String("env", "dev", "env flag")

	_ = flag.CommandLine.Parse([]string{"-name=foo", "-env=prod"})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	Require("name", "env")

	if got := *name; got != "foo" {
		t.Errorf("expected name=foo, got %q", got)
	}
	if got := *env; got != "prod" {
		t.Errorf("expected env=prod, got %q", got)
	}
	if buf.Len() > 0 {
		t.Errorf("unexpected log output: %s", buf.String())
	}
}

// For failure scenarios we fork the test in a subprocess
func TestRequire_MissingFlags(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		_ = flag.CommandLine.Parse([]string{})

		Require("missing")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRequire_MissingFlags")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	out, err := cmd.CombinedOutput()

	if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
		if !strings.Contains(string(out), `"missing" is not a valid flag`) {
			t.Errorf("expected log to mention missing flag, got %q", string(out))
		}
	} else {
		t.Fatalf("expected process to exit with error, got err=%v, out=%s", err, string(out))
	}
}
