// Copyright 2019 The Hugo Authors. All rights reserved.
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

package htesting

import (
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/afero"
)

// IsTest reports whether we're running as a test.
var IsTest bool

func init() {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			IsTest = true
			break
		}
	}
}

// CreateTempDir creates a temp dir in the given filesystem and
// returns the dirname and a func that removes it when done.
func CreateTempDir(fs afero.Fs, prefix string) (string, func(), error) {
	tempDir, err := afero.TempDir(fs, "", prefix)
	if err != nil {
		return "", nil, err
	}

	_, isOsFs := fs.(*afero.OsFs)

	if isOsFs && runtime.GOOS == "darwin" && !strings.HasPrefix(tempDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		tempDir = "/private" + tempDir
	}
	return tempDir, func() { fs.RemoveAll(tempDir) }, nil
}

// BailOut panics with a stack trace after the given duration. Useful for
// hanging tests.
func BailOut(after time.Duration) {
	time.AfterFunc(after, func() {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, true)
		panic(string(buf))
	})
}

// Rnd is used only for testing.
var Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandBool() bool {
	return Rnd.Intn(2) != 0
}

// DiffStringSlices returns the difference between two string slices.
// Useful in tests.
// See:
// http://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings-in-golang
func DiffStringSlices(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}

// DiffStrings splits the strings into fields and runs it into DiffStringSlices.
// Useful for tests.
func DiffStrings(s1, s2 string) []string {
	return DiffStringSlices(strings.Fields(s1), strings.Fields(s2))
}

// IsCI reports whether we're running in a CI server.
func IsCI() bool {
	return (os.Getenv("CI") != "" || os.Getenv("CI_LOCAL") != "") && os.Getenv("CIRCLE_BRANCH") == ""
}

// IsGitHubAction reports whether we're running in a GitHub Action.
func IsGitHubAction() bool {
	return os.Getenv("GITHUB_ACTION") != ""
}

// SupportsAll reports whether the running system supports all Hugo features,
// e.g. Asciidoc, Pandoc etc.
func SupportsAll() bool {
	return IsGitHubAction() || os.Getenv("CI_LOCAL") != ""
}

// GoMinorVersion returns the minor version of the current Go version,
// e.g. 16 for Go 1.16.
func GoMinorVersion() int {
	return extractMinorVersionFromGoTag(runtime.Version())
}

// IsWindows reports whether this runs on Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

var goMinorVersionRe = regexp.MustCompile(`go1.(\d*)`)

func extractMinorVersionFromGoTag(tag string) int {
	// The tag may be on the form go1.17, go1.17.5 go1.17rc2 -- or just a commit hash.
	match := goMinorVersionRe.FindStringSubmatch(tag)

	if len(match) == 2 {
		i, err := strconv.Atoi(match[1])
		if err != nil {
			return -1
		}
		return i
	}

	// a commit hash, not useful.
	return -1
}

// NewPinnedRunner creates a new runner that will only Run tests matching the given regexp.
// This is added mostly to use in combination with https://marketplace.visualstudio.com/items?itemName=windmilleng.vscode-go-autotest
func NewPinnedRunner(t testing.TB, pinnedTestRe string) *PinnedRunner {
	if pinnedTestRe == "" {
		pinnedTestRe = ".*"
	}
	pinnedTestRe = strings.ReplaceAll(pinnedTestRe, "_", " ")
	re := regexp.MustCompile("(?i)" + pinnedTestRe)
	return &PinnedRunner{
		c:  qt.New(t),
		re: re,
	}
}

type PinnedRunner struct {
	c  *qt.C
	re *regexp.Regexp
}

func (r *PinnedRunner) Run(name string, f func(c *qt.C)) bool {
	if !r.re.MatchString(name) {
		if IsGitHubAction() {
			r.c.Fatal("found pinned test when running in CI")
		}
		return true
	}
	return r.c.Run(name, f)
}
