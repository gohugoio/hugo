// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testenv_test

import (
	"github.com/gohugoio/hugo/tpl/internal/go_templates/testenv"
	//"internal/platform"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
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

func TestHasGoBuild(t *testing.T) {
	// Removed by Hugo.
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
		if b := testenv.Builder(); strings.HasSuffix(b, "-corellium") && !hasExec {
			// Most ios environments can't exec, but the corellium builder can.
			t.Errorf("expected MustHaveExec not to skip on %v", b)
		}
	default:
		if b := testenv.Builder(); b != "" && !hasExec {
			t.Errorf("expected MustHaveExec not to skip on %v", b)
		}
	}
}
