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
	"os"
	"runtime"
	"testing"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/viper"
)

func TestEmptySourceFilesystem(t *testing.T) {
	ss := newTestSourceSpec()
	src := ss.NewFilesystem("Empty")
	if len(src.Files()) != 0 {
		t.Errorf("new filesystem should contain 0 files.")
	}
}

type TestPath struct {
	filename string
	logical  string
	content  string
	section  string
	dir      string
}

func TestUnicodeNorm(t *testing.T) {
	if runtime.GOOS != "darwin" {
		// Normalization code is only for Mac OS, since it is not necessary for other OSes.
		return
	}

	paths := []struct {
		NFC string
		NFD string
	}{
		{NFC: "å", NFD: "\x61\xcc\x8a"},
		{NFC: "é", NFD: "\x65\xcc\x81"},
	}

	ss := newTestSourceSpec()
	var fi os.FileInfo

	for _, path := range paths {
		src := ss.NewFilesystem("base")
		_ = src.add(path.NFD, fi)
		f := src.Files()[0]
		if f.BaseFileName() != path.NFC {
			t.Fatalf("file name in NFD form should be normalized (%s)", path.NFC)
		}
	}

}

func newTestConfig() *viper.Viper {
	v := viper.New()
	v.Set("contentDir", "content")
	v.Set("dataDir", "data")
	v.Set("i18nDir", "i18n")
	v.Set("layoutDir", "layouts")
	v.Set("archetypeDir", "archetypes")
	v.Set("resourceDir", "resources")
	v.Set("publishDir", "public")
	v.Set("assetDir", "assets")
	return v
}

func newTestSourceSpec() *SourceSpec {
	v := newTestConfig()
	fs := hugofs.NewMem(v)
	ps, err := helpers.NewPathSpec(fs, v)
	if err != nil {
		panic(err)
	}
	return NewSourceSpec(ps, fs.Source)
}
