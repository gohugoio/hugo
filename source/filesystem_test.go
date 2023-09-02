// Copyright 2023 The Hugo Authors. All rights reserved.
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

package source_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/afero"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/source"
)

func TestEmptySourceFilesystem(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ss := newTestSourceSpec()
	src := ss.NewFilesystem("")
	files, err := src.Files()
	c.Assert(err, qt.IsNil)
	if len(files) != 0 {
		t.Errorf("new filesystem should contain 0 files.")
	}
}

func TestUnicodeNorm(t *testing.T) {
	t.Parallel()
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
		files, err := src.Files()
		c.Assert(err, qt.IsNil)
		f := files[0]
		if f.BaseFileName() != path.NFC {
			t.Fatalf("file %q name in NFD form should be normalized (%s)", f.BaseFileName(), path.NFC)
		}
	}
}

func newTestSourceSpec() *source.SourceSpec {
	v := config.New()
	afs := hugofs.NewBaseFileDecorator(afero.NewMemMapFs())
	conf := testconfig.GetTestConfig(afs, v)
	fs := hugofs.NewFrom(afs, conf.BaseConfig())
	ps, err := helpers.NewPathSpec(fs, conf, nil)
	if err != nil {
		panic(err)
	}
	return source.NewSourceSpec(ps, nil, fs.Source)
}
