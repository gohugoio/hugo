// Copyright 2024 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package releaser implements a set of utilities to help automate the
// Hugo release process.
package releaser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/hugo"
)

const commitPrefix = "releaser:"

// New initializes a ReleaseHandler.
func New(skipPush, try bool, step int) (*ReleaseHandler, error) {
	if step < 1 || step > 2 {
		return nil, fmt.Errorf("step must be 1 or 2")
	}

	prefix := "release-"
	branch, err := git("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, err
	}
	branch = strings.TrimSpace(branch)

	if !strings.HasPrefix(branch, prefix) {
		return nil, fmt.Errorf("branch %q is not a release branch", branch)
	}

	version := strings.TrimPrefix(branch, prefix)
	version = strings.TrimPrefix(version, "v")

	logf("Branch: %s|Version: v%s\n", branch, version)

	rh := &ReleaseHandler{branchVersion: version, skipPush: skipPush, try: try, step: step}

	if try {
		rh.git = func(args ...string) (string, error) {
			logln("git", strings.Join(args, " "))
			return "", nil
		}
	} else {
		rh.git = git
	}

	return rh, nil
}

// ReleaseHandler provides functionality to release a new version of Hugo.
// Test this locally without doing an actual release:
// go run -tags release main.go release --skip-publish --try -r 0.90.0
// Or a variation of the above -- the skip-publish flag makes sure that any changes are performed to the local Git only.
type ReleaseHandler struct {
	branchVersion string

	// 1 or 2.
	step int

	// No remote pushes.
	skipPush bool

	// Just simulate, no actual changes.
	try bool

	git func(args ...string) (string, error)
}

// Run creates a new release.
func (r *ReleaseHandler) Run() error {
	newVersion, finalVersion := r.calculateVersions()
	version := newVersion.String()
	tag := "v" + version
	mainVersion := newVersion
	mainVersion.PatchLevel = 0

	r.gitPull()

	defer r.gitPush()

	if r.step == 1 {
		if err := r.bumpVersions(newVersion); err != nil {
			return err
		}

		if _, err := r.git("commit", "-a", "-m", fmt.Sprintf("%s Bump versions for release of %s\n\n[ci skip]", commitPrefix, newVersion)); err != nil {
			return err
		}

		// The above commit will be the target for this release, so print it to the console in a env friendly way.
		sha, err := git("rev-parse", "HEAD")
		if err != nil {
			return err
		}

		// Hugoreleaser will do the actual release using these values.
		if err := r.replaceInFile("hugoreleaser.env",
			`HUGORELEASER_TAG=(\S*)`, "HUGORELEASER_TAG="+tag,
			`HUGORELEASER_COMMITISH=(\S*)`, "HUGORELEASER_COMMITISH="+sha,
		); err != nil {
			return err
		}
		logf("HUGORELEASER_TAG=%s\n", tag)
		logf("HUGORELEASER_COMMITISH=%s\n", sha)

		return nil
	}

	if err := r.bumpVersions(finalVersion); err != nil {
		return err
	}

	if _, err := r.git("commit", "-a", "-m", fmt.Sprintf("%s Prepare repository for %s\n\n[ci skip]", commitPrefix, finalVersion)); err != nil {
		return err
	}

	return nil
}

func (r *ReleaseHandler) bumpVersions(ver hugo.Version) error {
	toDev := ""

	if ver.Suffix != "" {
		toDev = ver.Suffix
	}

	if err := r.replaceInFile("common/hugo/version_current.go",
		`Minor:(\s*)(\d*),`, fmt.Sprintf(`Minor:${1}%d,`, ver.Minor),
		`PatchLevel:(\s*)(\d*),`, fmt.Sprintf(`PatchLevel:${1}%d,`, ver.PatchLevel),
		`Suffix:(\s*)".*",`, fmt.Sprintf(`Suffix:${1}"%s",`, toDev)); err != nil {
		return err
	}

	var minVersion string
	if ver.Suffix != "" {
		// People use the DEV version in daily use, and we cannot create new themes
		// with the next version before it is released.
		minVersion = ver.Prev().String()
	} else {
		minVersion = ver.String()
	}

	if err := r.replaceInFile("commands/new.go",
		`min_version = "(.*)"`, fmt.Sprintf(`min_version = "%s"`, minVersion)); err != nil {
		return err
	}

	return nil
}

func (r ReleaseHandler) calculateVersions() (hugo.Version, hugo.Version) {
	newVersion := hugo.MustParseVersion(r.branchVersion)
	finalVersion := newVersion.Next()
	finalVersion.PatchLevel = 0

	if newVersion.Suffix != "-test" {
		newVersion.Suffix = ""
	}

	finalVersion.Suffix = "-DEV"

	return newVersion, finalVersion
}

func (r *ReleaseHandler) gitPull() {
	if _, err := r.git("pull", "origin", "HEAD"); err != nil {
		log.Fatal("pull failed:", err)
	}
}

func (r *ReleaseHandler) gitPush() {
	if r.skipPush {
		return
	}
	if _, err := r.git("push", "origin", "HEAD"); err != nil {
		log.Fatal("push failed:", err)
	}
}

func (r *ReleaseHandler) replaceInFile(filename string, oldNew ...string) error {
	filename = filepath.FromSlash(filename)
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if r.try {
		logf("Replace in %q: %q\n", filename, oldNew)
		return nil
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	newContent := string(b)

	for i := 0; i < len(oldNew); i += 2 {
		re := regexp.MustCompile(oldNew[i])
		newContent = re.ReplaceAllString(newContent, oldNew[i+1])
	}

	return os.WriteFile(filename, []byte(newContent), fi.Mode())
}

func git(args ...string) (string, error) {
	cmd, _ := hexec.SafeCommand("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git failed: %q: %q (%q)", err, out, args)
	}
	return string(out), nil
}

func logf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func logln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
