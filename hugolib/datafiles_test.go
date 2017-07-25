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

package hugolib

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"io/ioutil"
	"log"
	"os"

	"github.com/gohugoio/hugo/deps"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/gohugoio/hugo/parser"
	"github.com/gohugoio/hugo/source"
	"github.com/stretchr/testify/require"
)

func TestDataDirJSON(t *testing.T) {
	t.Parallel()

	sources := []source.ByteSource{
		{Name: filepath.FromSlash("data/test/foo.json"), Content: []byte(`{ "bar": "foofoo"  }`)},
		{Name: filepath.FromSlash("data/test.json"), Content: []byte(`{ "hello": [ { "world": "foo" } ] }`)},
	}

	expected, err := parser.HandleJSONMetaData([]byte(`{ "test": { "hello": [{ "world": "foo"  }] , "foo": { "bar":"foofoo" } } }`))

	if err != nil {
		t.Fatalf("Error %s", err)
	}

	doTestDataDir(t, expected, sources)
}

func TestDataDirToml(t *testing.T) {
	t.Parallel()

	sources := []source.ByteSource{
		{Name: filepath.FromSlash("data/test/kung.toml"), Content: []byte("[foo]\nbar = 1")},
	}

	expected, err := parser.HandleTOMLMetaData([]byte("[test]\n[test.kung]\n[test.kung.foo]\nbar = 1"))

	if err != nil {
		t.Fatalf("Error %s", err)
	}

	doTestDataDir(t, expected, sources)
}

func TestDataDirYAMLWithOverridenValue(t *testing.T) {
	t.Parallel()

	sources := []source.ByteSource{
		// filepath.Walk walks the files in lexical order, '/' comes before '.'. Simulate this:
		{Name: filepath.FromSlash("data/a.yaml"), Content: []byte("a: 1")},
		{Name: filepath.FromSlash("data/test/v1.yaml"), Content: []byte("v1-2: 2")},
		{Name: filepath.FromSlash("data/test/v2.yaml"), Content: []byte("v2:\n- 2\n- 3")},
		{Name: filepath.FromSlash("data/test.yaml"), Content: []byte("v1: 1")},
	}

	expected := map[string]interface{}{"a": map[string]interface{}{"a": 1},
		"test": map[string]interface{}{"v1": map[string]interface{}{"v1-2": 2}, "v2": map[string]interface{}{"v2": []interface{}{2, 3}}}}

	doTestDataDir(t, expected, sources)
}

// issue 892
func TestDataDirMultipleSources(t *testing.T) {
	t.Parallel()

	sources := []source.ByteSource{
		{Name: filepath.FromSlash("data/test/first.toml"), Content: []byte("bar = 1")},
		{Name: filepath.FromSlash("themes/mytheme/data/test/first.toml"), Content: []byte("bar = 2")},
		{Name: filepath.FromSlash("data/test/second.toml"), Content: []byte("tender = 2")},
	}

	expected, _ := parser.HandleTOMLMetaData([]byte("[test.first]\nbar = 1\n[test.second]\ntender=2"))

	doTestDataDir(t, expected, sources,
		"theme", "mytheme")

}

func doTestDataDir(t *testing.T, expected interface{}, sources []source.ByteSource, configKeyValues ...interface{}) {
	var (
		cfg, fs = newTestCfg()
	)

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}

	var (
		logger  = jww.NewNotepad(jww.LevelError, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
		depsCfg = deps.DepsCfg{Fs: fs, Cfg: cfg, Logger: logger}
	)

	writeSource(t, fs, filepath.Join("content", "dummy.md"), "content")
	writeSourcesToSource(t, "", fs, sources...)

	expectBuildError := false

	if ok, shouldFail := expected.(bool); ok && shouldFail {
		expectBuildError = true
	}

	s := buildSingleSiteExpected(t, expectBuildError, depsCfg, BuildCfg{SkipRender: true})

	if !expectBuildError && !reflect.DeepEqual(expected, s.Data) {
		t.Errorf("Expected structure\n%#v got\n%#v", expected, s.Data)
	}
}

func TestDataFromShortcode(t *testing.T) {
	t.Parallel()

	var (
		cfg, fs = newTestCfg()
	)

	writeSource(t, fs, "data/hugo.toml", "slogan = \"Hugo Rocks!\"")
	writeSource(t, fs, "layouts/_default/single.html", `
* Slogan from template: {{  .Site.Data.hugo.slogan }}
* {{ .Content }}`)
	writeSource(t, fs, "layouts/shortcodes/d.html", `{{  .Page.Site.Data.hugo.slogan }}`)
	writeSource(t, fs, "content/c.md", `---
---
Slogan from shortcode: {{< d >}}
`)

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	content := readSource(t, fs, "public/c/index.html")
	require.True(t, strings.Contains(content, "Slogan from template: Hugo Rocks!"), content)
	require.True(t, strings.Contains(content, "Slogan from shortcode: Hugo Rocks!"), content)

}
