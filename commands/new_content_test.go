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
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/viper"

	qt "github.com/frankban/quicktest"
)

// Issue #1133
func TestNewContentPathSectionWithForwardSlashes(t *testing.T) {
	c := qt.New(t)
	p, s := newContentPathSection(nil, "/post/new.md")
	c.Assert(p, qt.Equals, filepath.FromSlash("/post/new.md"))
	c.Assert(s, qt.Equals, "post")
}

func checkNewSiteInited(fs *hugofs.Fs, basepath string, t *testing.T) {
	c := qt.New(t)
	paths := []string{
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		filepath.Join(basepath, "archetypes"),
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
		filepath.Join(basepath, "config.toml"),
	}

	for _, path := range paths {
		_, err := fs.Source.Stat(path)
		c.Assert(err, qt.IsNil)
	}
}

func TestDoNewSite(t *testing.T) {
	c := qt.New(t)
	n := newNewSiteCmd()
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()

	c.Assert(n.doNewSite(fs, basepath, false), qt.IsNil)

	checkNewSiteInited(fs, basepath, t)
}

func TestDoNewSite_noerror_base_exists_but_empty(t *testing.T) {
	c := qt.New(t)
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()
	n := newNewSiteCmd()

	c.Assert(fs.Source.MkdirAll(basepath, 0777), qt.IsNil)

	c.Assert(n.doNewSite(fs, basepath, false), qt.IsNil)
}

func TestDoNewSite_error_base_exists(t *testing.T) {
	c := qt.New(t)
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()
	n := newNewSiteCmd()

	c.Assert(fs.Source.MkdirAll(basepath, 0777), qt.IsNil)
	_, err := fs.Source.Create(filepath.Join(basepath, "foo"))
	c.Assert(err, qt.IsNil)
	// Since the directory already exists and isn't empty, expect an error
	c.Assert(n.doNewSite(fs, basepath, false), qt.Not(qt.IsNil))

}

func TestDoNewSite_force_empty_dir(t *testing.T) {
	c := qt.New(t)
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()
	n := newNewSiteCmd()

	c.Assert(fs.Source.MkdirAll(basepath, 0777), qt.IsNil)
	c.Assert(n.doNewSite(fs, basepath, true), qt.IsNil)

	checkNewSiteInited(fs, basepath, t)
}

func TestDoNewSite_error_force_dir_inside_exists(t *testing.T) {
	c := qt.New(t)
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()
	n := newNewSiteCmd()

	contentPath := filepath.Join(basepath, "content")

	c.Assert(fs.Source.MkdirAll(contentPath, 0777), qt.IsNil)
	c.Assert(n.doNewSite(fs, basepath, true), qt.Not(qt.IsNil))
}

func TestDoNewSite_error_force_config_inside_exists(t *testing.T) {
	c := qt.New(t)
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()
	n := newNewSiteCmd()

	configPath := filepath.Join(basepath, "config.toml")
	c.Assert(fs.Source.MkdirAll(basepath, 0777), qt.IsNil)
	_, err := fs.Source.Create(configPath)
	c.Assert(err, qt.IsNil)

	c.Assert(n.doNewSite(fs, basepath, true), qt.Not(qt.IsNil))
}

func newTestCfg() (*viper.Viper, *hugofs.Fs) {

	v := viper.New()
	fs := hugofs.NewMem(v)

	v.SetFs(fs.Source)

	return v, fs

}
