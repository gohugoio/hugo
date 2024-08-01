// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testenv_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/tpl/internal/go_templates/testenv"
)

func TestGoToolLocation(t *testing.T) {
	testenv.MustHaveGoBuild(t)

	var exeSuffix string
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}

	// Tests are defined to run within their package source directory,
	// and this package's source directory is $GOROOT/src/internal/testenv.
	// The 'go' command is installed at $GOROOT/bin/go, so if the environment
	// is correct then testenv.GoTool() should be identical to ../../../bin/go.

	relWant := "../../../bin/go" + exeSuffix
	absWant, err := filepath.Abs(relWant)
	if err != nil {
		t.Fatal(err)
	}

	wantInfo, err := os.Stat(absWant)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found go tool at %q (%q)", relWant, absWant)

	goTool, err := testenv.GoTool()
	if err != nil {
		t.Fatalf("testenv.GoTool(): %v", err)
	}
	t.Logf("testenv.GoTool() = %q", goTool)

	gotInfo, err := os.Stat(goTool)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(wantInfo, gotInfo) {
		t.Fatalf("%q is not the same file as %q", absWant, goTool)
	}
}

// Modified by Hugo (not needed)
func TestHasGoBuild(t *testing.T) {
}

func TestMustHaveExec(t *testing.T) {
	hasExec := false
	t.Run("MustHaveExec", func(t *testing.T) {
		testenv.MustHaveExec(t)
		t.Logf("MustHaveExec did not skip")
		hasExec = true
	})

	switch runtime.GOOS {
	case "js", "wasip1":
		if hasExec {
			// js and wasip1 lack an “exec” syscall.
			t.Errorf("expected MustHaveExec to skip on %v", runtime.GOOS)
		}
	case "ios":
		if b := testenv.Builder(); isCorelliumBuilder(b) && !hasExec {
			// Most ios environments can't exec, but the corellium builder can.
			t.Errorf("expected MustHaveExec not to skip on %v", b)
		}
	default:
		if b := testenv.Builder(); b != "" && !hasExec {
			t.Errorf("expected MustHaveExec not to skip on %v", b)
		}
	}
}

func TestCleanCmdEnvPWD(t *testing.T) {
	// Test that CleanCmdEnv sets PWD if cmd.Dir is set.
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("PWD is not used on %s", runtime.GOOS)
	}
	dir := t.TempDir()
	cmd := testenv.Command(t, testenv.GoToolPath(t), "help")
	cmd.Dir = dir
	cmd = testenv.CleanCmdEnv(cmd)

	for _, env := range cmd.Env {
		if strings.HasPrefix(env, "PWD=") {
			pwd := strings.TrimPrefix(env, "PWD=")
			if pwd != dir {
				t.Errorf("unexpected PWD: want %s, got %s", dir, pwd)
			}
			return
		}
	}
	t.Error("PWD not set in cmd.Env")
}

func isCorelliumBuilder(builderName string) bool {
	// Support both the old infra's builder names and the LUCI builder names.
	// The former's names are ad-hoc so we could maintain this invariant on
	// the builder side. The latter's names are structured, and "corellium" will
	// appear as a "host" suffix after the GOOS and GOARCH, which always begin
	// with an underscore.
	return strings.HasSuffix(builderName, "-corellium") || strings.Contains(builderName, "_corellium")
}

func isEmulatedBuilder(builderName string) bool {
	// Support both the old infra's builder names and the LUCI builder names.
	// The former's names are ad-hoc so we could maintain this invariant on
	// the builder side. The latter's names are structured, and the signifier
	// of emulation "emu" will appear as a "host" suffix after the GOOS and
	// GOARCH because it modifies the run environment in such a way that it
	// the target GOOS and GOARCH may not match the host. This suffix always
	// begins with an underscore.
	return strings.HasSuffix(builderName, "-emu") || strings.Contains(builderName, "_emu")
}
