// Copyright 2015 The Hugo Authors. All rights reserved.
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

package source

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/langs"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
)

func TestUnicodeNorm(t *testing.T) {
	if runtime.GOOS != "darwin" {
		// Normalization code is only for Mac OS, since it is not necessary for other OSes.
		return
	}

	c := qt.New(t)

	paths := []struct {
		NFC string
		NFD string
	}{
		{NFC: "å", NFD: "\x61\xcc\x8a"},
		{NFC: "é", NFD: "\x65\xcc\x81"},
	}

	ss := newTestSourceSpec()

	for i, path := range paths {
		base := fmt.Sprintf("base%d", i)
		c.Assert(afero.WriteFile(ss.Fs.Source, filepath.Join(base, path.NFD), []byte("some data"), 0777), qt.IsNil)
		src := ss.NewFilesystem(base)
		var found bool
		err := src.Walk(func(f File) error {
			found = true
			c.Assert(f.BaseFileName(), qt.Equals, path.NFC)
			return nil
		})
		c.Assert(err, qt.IsNil)
		c.Assert(found, qt.IsTrue)

	}
}

func newTestConfig() config.Provider {
	v := config.New()
	v.Set("contentDir", "content")
	v.Set("dataDir", "data")
	v.Set("i18nDir", "i18n")
	v.Set("layoutDir", "layouts")
	v.Set("archetypeDir", "archetypes")
	v.Set("resourceDir", "resources")
	v.Set("publishDir", "public")
	v.Set("assetDir", "assets")
	_, err := langs.LoadLanguageSettings(v, nil)
	if err != nil {
		panic(err)
	}
	mod, err := modules.CreateProjectModule(v)
	if err != nil {
		panic(err)
	}
	v.Set("allModules", modules.Modules{mod})

	return v
}

func newTestSourceSpec() *SourceSpec {
	v := newTestConfig()
	fs := hugofs.NewFrom(hugofs.NewBaseFileDecorator(afero.NewMemMapFs()), v)
	ps, err := helpers.NewPathSpec(fs, v, nil)
	if err != nil {
		panic(err)
	}
	return NewSourceSpec(ps, nil, fs.Source)
}
