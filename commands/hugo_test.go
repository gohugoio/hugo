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

package commands

import (
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bep/clock"
	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
	"golang.org/x/tools/txtar"
)

// Issue #5662
func TestHugoWithContentDirOverride(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	files := `
-- config.toml --
baseURL = "https://example.org"
title = "Hugo Commands"
-- mycontent/p1.md --
---
title: "P1"
---
-- layouts/_default/single.html --
Page: {{ .Title }}|

`
	s := newTestHugoCmdBuilder(c, files, []string{"-c", "mycontent"}).Build()
	s.AssertFileContent("public/p1/index.html", `Page: P1|`)

}

// Issue #9794
func TestHugoStaticFilesMultipleStaticAndManyFolders(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	files := `
-- config.toml --
baseURL = "https://example.org"
theme = "mytheme"
-- layouts/index.html --
Home.

`
	const (
		numDirs     = 33
		numFilesMax = 12
	)

	r := rand.New(rand.NewSource(32))

	for i := 0; i < numDirs; i++ {
		for j := 0; j < r.Intn(numFilesMax); j++ {
			if j%3 == 0 {
				files += fmt.Sprintf("-- themes/mytheme/static/d%d/f%d.txt --\nHellot%d-%d\n", i, j, i, j)
				files += fmt.Sprintf("-- themes/mytheme/static/d%d/ft%d.txt --\nHellot%d-%d\n", i, j, i, j)
			}
			files += fmt.Sprintf("-- static/d%d/f%d.txt --\nHello%d-%d\n", i, j, i, j)
		}
	}

	r = rand.New(rand.NewSource(32))

	s := newTestHugoCmdBuilder(c, files, []string{"-c", "mycontent"}).Build()
	for i := 0; i < numDirs; i++ {
		for j := 0; j < r.Intn(numFilesMax); j++ {
			if j%3 == 0 {
				if j%3 == 0 {
					s.AssertFileContent(fmt.Sprintf("public/d%d/ft%d.txt", i, j), fmt.Sprintf("Hellot%d-%d", i, j))
				}
				s.AssertFileContent(fmt.Sprintf("public/d%d/f%d.txt", i, j), fmt.Sprintf("Hello%d-%d", i, j))
			}
		}
	}

}

// Issue #8787
func TestHugoListCommandsWithClockFlag(t *testing.T) {
	t.Cleanup(func() { htime.Clock = clock.System() })

	c := qt.New(t)

	files := `
-- config.toml --
baseURL = "https://example.org"
title = "Hugo Commands"
timeZone = "UTC"
-- content/past.md --
---
title: "Past"
date: 2000-11-06
---
-- content/future.md --
---
title: "Future"
date: 2200-11-06
---
-- layouts/_default/single.html --
Page: {{ .Title }}|

`
	s := newTestHugoCmdBuilder(c, files, []string{"list", "future"})
	s.captureOut = true
	s.Build()
	p := filepath.Join("content", "future.md")
	s.AssertStdout(p + ",2200-11-06T00:00:00Z")

	s = newTestHugoCmdBuilder(c, files, []string{"list", "future", "--clock", "2300-11-06"}).Build()
	s.AssertStdout("")
}

type testHugoCmdBuilder struct {
	*qt.C

	fs    afero.Fs
	dir   string
	files string
	args  []string

	captureOut bool
	out        string
}

func newTestHugoCmdBuilder(c *qt.C, files string, args []string) *testHugoCmdBuilder {
	s := &testHugoCmdBuilder{C: c, files: files, args: args}
	s.dir = s.TempDir()
	s.fs = afero.NewBasePathFs(hugofs.Os, s.dir)

	return s
}

func (s *testHugoCmdBuilder) Build() *testHugoCmdBuilder {
	data := txtar.Parse([]byte(s.files))

	for _, f := range data.Files {
		filename := filepath.Clean(f.Name)
		data := bytes.TrimSuffix(f.Data, []byte("\n"))
		s.Assert(s.fs.MkdirAll(filepath.Dir(filename), 0777), qt.IsNil)
		s.Assert(afero.WriteFile(s.fs, filename, data, 0666), qt.IsNil)
	}

	hugoCmd := newCommandsBuilder().addAll().build()
	cmd := hugoCmd.getCommand()
	args := append(s.args, "-s="+s.dir, "--quiet")
	cmd.SetArgs(args)

	if s.captureOut {
		out, err := captureStdout(func() error {
			_, err := cmd.ExecuteC()
			return err
		})
		s.Assert(err, qt.IsNil)
		s.out = out
	} else {
		_, err := cmd.ExecuteC()
		s.Assert(err, qt.IsNil)
	}

	return s
}

func (s *testHugoCmdBuilder) AssertFileContent(filename string, matches ...string) {
	s.Helper()
	data, err := afero.ReadFile(s.fs, filename)
	s.Assert(err, qt.IsNil)
	content := strings.TrimSpace(string(data))
	for _, m := range matches {
		lines := strings.Split(m, "\n")
		for _, match := range lines {
			match = strings.TrimSpace(match)
			if match == "" || strings.HasPrefix(match, "#") {
				continue
			}
			s.Assert(content, qt.Contains, match, qt.Commentf(m))
		}
	}
}

func (s *testHugoCmdBuilder) AssertStdout(match string) {
	s.Helper()
	content := strings.TrimSpace(s.out)
	s.Assert(content, qt.Contains, strings.TrimSpace(match))
}
