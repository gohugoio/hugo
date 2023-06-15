// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testenv_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gohugoio/hugo/tpl/internal/go_templates/testenv"
)

func TestGoToolLocation(t *testing.T) {
	t.Skip("skipping test that requires go command")
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
