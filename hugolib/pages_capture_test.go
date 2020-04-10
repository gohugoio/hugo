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

package hugolib

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"

	"github.com/gohugoio/hugo/common/loggers"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/afero"
)

func TestPagesCapture(t *testing.T) {

	cfg, hfs := newTestCfg()
	fs := hfs.Source

	c := qt.New(t)

	var writeFile = func(filename string) {
		c.Assert(afero.WriteFile(fs, filepath.FromSlash(filename), []byte(fmt.Sprintf("content-%s", filename)), 0755), qt.IsNil)
	}

	writeFile("_index.md")
	writeFile("logo.png")
	writeFile("root.md")
	writeFile("blog/index.md")
	writeFile("blog/hello.md")
	writeFile("blog/images/sunset.png")
	writeFile("pages/page1.md")
	writeFile("pages/page2.md")
	writeFile("pages/page.png")

	ps, err := helpers.NewPathSpec(hugofs.NewFrom(fs, cfg), cfg, loggers.NewErrorLogger())
	c.Assert(err, qt.IsNil)
	sourceSpec := source.NewSourceSpec(ps, fs)

	t.Run("Collect", func(t *testing.T) {
		c := qt.New(t)
		proc := &testPagesCollectorProcessor{}
		coll := newPagesCollector(sourceSpec, nil, loggers.NewErrorLogger(), nil, proc)
		c.Assert(coll.Collect(), qt.IsNil)
		c.Assert(len(proc.items), qt.Equals, 4)
	})

}

type testPagesCollectorProcessor struct {
	items   []interface{}
	waitErr error
}

func (proc *testPagesCollectorProcessor) Process(item interface{}) error {
	proc.items = append(proc.items, item)
	return nil
}
func (proc *testPagesCollectorProcessor) Start(ctx context.Context) context.Context {
	return ctx
}

func (proc *testPagesCollectorProcessor) Wait() error { return proc.waitErr }
