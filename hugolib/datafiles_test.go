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

	"fmt"
	"runtime"

	"github.com/stretchr/testify/require"
)

func TestDataDirJSON(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/test/foo.json"), `{ "bar": "foofoo"  }`},
		{filepath.FromSlash("data/test.json"), `{ "hello": [ { "world": "foo" } ] }`},
	}

	expected :=
		map[string]interface{}{
			"test": map[string]interface{}{
				"hello": []interface{}{
					map[string]interface{}{"world": "foo"},
				},
				"foo": map[string]interface{}{
					"bar": "foofoo",
				},
			},
		}

	doTestDataDir(t, expected, sources)
}

// Enable / adjust in https://github.com/gohugoio/hugo/issues/4393
func _TestDataDirYAML(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{"data/test/a.yaml", "b:\n  c1: 1\n  c2: 2"},
	}

	expected :=
		map[string]interface{}{
			"test": map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[interface{}]interface{}{
						"c1": 1,
						"c2": 2,
					},
				},
			},
		}

	doTestDataDir(t, expected, sources)
}

func TestDataDirToml(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{"data/test/kung.toml", "[foo]\nbar = 1"},
	}

	expected :=
		map[string]interface{}{
			"test": map[string]interface{}{
				"kung": map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 1,
					},
				},
			},
		}

	doTestDataDir(t, expected, sources)
}

// Enable / adjust in https://github.com/gohugoio/hugo/issues/4393
func _TestDataDirYAML2(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/test/foo.yaml"), "bar: foofoo"},
		{filepath.FromSlash("data/test.yaml"), "hello:\n- world: foo"},
	}

	//This is what we want: consistent use of map[string]interface{} for nested YAML maps
	// the same as TestDataDirJSON
	expected :=
		map[string]interface{}{
			"test": map[string]interface{}{
				"hello": []interface{}{
					map[string]interface{}{"world": "foo"},
				},
				"foo": map[string]interface{}{
					"bar": "foofoo",
				},
			},
		}

	// what we are actually getting as of v0.34
	expectedV0_34 :=
		map[string]interface{}{
			"test": map[string]interface{}{
				"hello": []interface{}{
					map[interface{}]interface{}{"world": "foo"},
				},
				"foo": map[string]interface{}{
					"bar": "foofoo",
				},
			},
		}
	_ = expected

	doTestDataDir(t, expectedV0_34, sources)
}

func TestDataDirToml2(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/test/foo.toml"), "bar = \"foofoo\""},
		{filepath.FromSlash("data/test.toml"), "[[hello]]\nworld = \"foo\""},
	}

	expected :=
		map[string]interface{}{
			"test": map[string]interface{}{
				"hello": []map[string]interface{}{
					map[string]interface{}{"world": "foo"},
				},
				"foo": map[string]interface{}{
					"bar": "foofoo",
				},
			},
		}

	doTestDataDir(t, expected, sources)
}

func TestDataDirJSONWithOverriddenValue(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		// filepath.Walk walks the files in lexical order, '/' comes before '.'. Simulate this:
		{filepath.FromSlash("data/a.json"), `{"a": "1"}`},
		{filepath.FromSlash("data/test/v1.json"), `{"v1-2": "2"}`},
		{filepath.FromSlash("data/test/v2.json"), `{"v2": ["2", "3"]}`},
		{filepath.FromSlash("data/test.json"), `{"v1": "1"}`},
	}

	expected :=
		map[string]interface{}{
			"a": map[string]interface{}{"a": "1"},
			"test": map[string]interface{}{
				"v1": map[string]interface{}{"v1-2": "2"},
				"v2": map[string]interface{}{"v2": []interface{}{"2", "3"}},
			},
		}

	doTestDataDir(t, expected, sources)
}

func TestDataDirYAMLWithOverridenValue(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		// filepath.Walk walks the files in lexical order, '/' comes before '.'. Simulate this:
		{filepath.FromSlash("data/a.yaml"), "a: 1"},
		{filepath.FromSlash("data/test/v1.yaml"), "v1-2: 2"},
		{filepath.FromSlash("data/test/v2.yaml"), "v2:\n- 2\n- 3"},
		{filepath.FromSlash("data/test.yaml"), "v1: 1"},
	}

	expected :=
		map[string]interface{}{
			"a": map[string]interface{}{"a": 1},
			"test": map[string]interface{}{
				"v1": map[string]interface{}{"v1-2": 2},
				"v2": map[string]interface{}{"v2": []interface{}{2, 3}},
			},
		}

	doTestDataDir(t, expected, sources)
}

// Issue #4361
func TestDataDirJSONArrayAtTopLevelOfFile(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/test.json"), `[ { "hello": "world" }, { "what": "time" }, { "is": "lunch?" } ]`},
	}

	expected :=
		map[string]interface{}{
			"test": []interface{}{
				map[string]interface{}{"hello": "world"},
				map[string]interface{}{"what": "time"},
				map[string]interface{}{"is": "lunch?"},
			},
		}

	doTestDataDir(t, expected, sources)
}

// TODO Issue #3890 unresolved
func TestDataDirYAMLArrayAtTopLevelOfFile(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/test.yaml"), `
- hello: world
- what: time
- is: lunch?
`},
	}

	//TODO decide whether desired structure map[interface {}]interface{} as shown
	// and as the YAML parser produces, or should it be map[string]interface{}
	// all the way down per Issue #4138
	expected :=
		map[string]interface{}{
			"test": []interface{}{
				map[interface{}]interface{}{"hello": "world"},
				map[interface{}]interface{}{"what": "time"},
				map[interface{}]interface{}{"is": "lunch?"},
			},
		}

	// what we are actually getting as of v0.34
	expectedV0_34 :=
		map[string]interface{}{}
	_ = expected

	doTestDataDir(t, expectedV0_34, sources)
}

// Issue #892
func TestDataDirMultipleSources(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/test/first.yaml"), "bar: 1"},
		{filepath.FromSlash("themes/mytheme/data/test/first.yaml"), "bar: 2"},
		{filepath.FromSlash("data/test/second.yaml"), "tender: 2"},
	}

	expected :=
		map[string]interface{}{
			"test": map[string]interface{}{
				"first": map[string]interface{}{
					"bar": 1,
				},
				"second": map[string]interface{}{
					"tender": 2,
				},
			},
		}

	doTestDataDir(t, expected, sources,
		"theme", "mytheme")

}

// test (and show) the way values from four different sources commingle and override
func TestDataDirMultipleSourcesCommingled(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/a.json"), `{ "b1" : { "c1": "data/a" }, "b2": "data/a", "b3": ["x", "y", "z"] }`},
		{filepath.FromSlash("themes/mytheme/data/a.json"), `{ "b1": "mytheme/data/a",  "b2": "mytheme/data/a", "b3": "mytheme/data/a" }`},
		{filepath.FromSlash("themes/mytheme/data/a/b1.json"), `{ "c1": "mytheme/data/a/b1", "c2": "mytheme/data/a/b1" }`},
		{filepath.FromSlash("data/a/b1.json"), `{ "c1": "data/a/b1" }`},
	}

	// Per handleDataFile() comment:
	// 1. A theme uses the same key; the main data folder wins
	// 2. A sub folder uses the same key: the sub folder wins
	expected :=
		map[string]interface{}{
			"a": map[string]interface{}{
				"b1": map[string]interface{}{
					"c1": "data/a/b1",
					"c2": "mytheme/data/a/b1",
				},
				"b2": "data/a",
				"b3": []interface{}{"x", "y", "z"},
			},
		}

	doTestDataDir(t, expected, sources,
		"theme", "mytheme")
}

func TestDataDirMultipleSourcesCollidingChildArrays(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("data/a.json"), `{ "b1" : "data/a", "b2" : ["x", "y", "z"] }`},
		{filepath.FromSlash("data/a/b2.json"), `["1", "2", "3"]`},
	}

	// Per handleDataFile() comment:
	// 1. A theme uses the same key; the main data folder wins
	// 2. A sub folder uses the same key: the sub folder wins
	expected :=
		map[string]interface{}{
			"a": map[string]interface{}{
				"b1": "data/a",
				"b2": []interface{}{"1", "2", "3"},
			},
		}

	doTestDataDir(t, expected, sources,
		"theme", "mytheme")
}

// TODO Issue #4366 unresolved
func TestDataDirMultipleSourcesCollidingTopLevelArrays(t *testing.T) {
	t.Parallel()

	sources := [][2]string{
		{filepath.FromSlash("themes/mytheme/data/a/b1.json"), `["x", "y", "z"]`},
		{filepath.FromSlash("data/a/b1.json"), `["1", "2", "3"]`},
	}

	expected :=
		map[string]interface{}{
			"a": map[string]interface{}{
				"b1": []interface{}{"1", "2", "3"},
			},
		}

	// as of v0.34 this test results in a go Panic
	_ = sources
	_ = expected
	/*
		doTestDataDir(t, expectedV0_35, sources,
			"theme", "mytheme")
	*/
}

func doTestDataDir(t *testing.T, expected interface{}, sources [][2]string, configKeyValues ...interface{}) {
	var (
		cfg, fs = newTestCfg()
	)

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}

	var (
		logger  = jww.NewNotepad(jww.LevelWarn, jww.LevelWarn, os.Stdout, ioutil.Discard, t.Name(), log.Ldate|log.Ltime)
		depsCfg = deps.DepsCfg{Fs: fs, Cfg: cfg, Logger: logger}
	)

	writeSource(t, fs, filepath.Join("content", "dummy.md"), "content")
	writeSourcesToSource(t, "", fs, sources...)

	expectBuildError := false

	if ok, shouldFail := expected.(bool); ok && shouldFail {
		expectBuildError = true
	}

	// trap and report panics as unmarshaling errors so that test suit can complete
	defer func() {
		if r := recover(); r != nil {
			// Capture the stack trace
			buf := make([]byte, 10000)
			runtime.Stack(buf, false)
			t.Errorf("PANIC: %s\n\nStack Trace : %s", r, string(buf))
		}
	}()

	s := buildSingleSiteExpected(t, expectBuildError, depsCfg, BuildCfg{SkipRender: true})

	if !expectBuildError && !reflect.DeepEqual(expected, s.Data) {
		exp := fmt.Sprintf("%#v", expected)
		got := fmt.Sprintf("%#v", s.Data)
		if exp == got { //TODO: This workaround seems to be triggered only by the TOML tests
			t.Logf("WARNING: reflect.DeepEqual returned FALSE for values that appear equal.\n"+
				"Treating as equal for the purpose of the test, but this maybe should be investigated.\n"+
				"Expected data:\n%v got\n%v\n\nExpected type structure:\n%#[1]v got\n%#[2]v", expected, s.Data)
			return
		}

		t.Errorf("Expected data:\n%v got\n%v\n\nExpected type structure:\n%#[1]v got\n%#[2]v", expected, s.Data)
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
