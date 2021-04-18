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

package modules

import (
	"os"
	"regexp"
	"testing"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/spf13/viper"

	"github.com/gohugoio/hugo/config"

	qt "github.com/frankban/quicktest"
)

func TestConfigHugoVersionIsValid(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		in     HugoVersion
		expect bool
	}{
		{HugoVersion{Min: "0.33.0"}, true},
		{HugoVersion{Min: "0.56.0-DEV"}, true},
		{HugoVersion{Min: "0.33.0", Max: "0.55.0"}, false},
		{HugoVersion{Min: "0.33.0", Max: "0.99.0"}, true},
	} {
		c.Assert(test.in.IsValid(), qt.Equals, test.expect)
	}
}

func TestDecodeConfig(t *testing.T) {
	c := qt.New(t)

	c.Run("Basic", func(c *qt.C) {
		tomlConfig := `
[module]

[module.hugoVersion]
min = "0.54.2"
max = "0.99.0"
extended = true

[[module.mounts]]
source="src/project/blog"
target="content/blog"
lang="en"
[[module.imports]]
path="github.com/bep/mycomponent"
[[module.imports.mounts]]
source="scss"
target="assets/bootstrap/scss"
[[module.imports.mounts]]
source="src/markdown/blog"
target="content/blog"
lang="en"
`
		cfg, err := config.FromConfigString(tomlConfig, "toml")
		c.Assert(err, qt.IsNil)

		mcfg, err := DecodeConfig(cfg)
		c.Assert(err, qt.IsNil)

		v056 := hugo.VersionString("0.56.0")

		hv := mcfg.HugoVersion

		c.Assert(v056.Compare(hv.Min), qt.Equals, -1)
		c.Assert(v056.Compare(hv.Max), qt.Equals, 1)
		c.Assert(hv.Extended, qt.Equals, true)

		if hugo.IsExtended {
			c.Assert(hv.IsValid(), qt.Equals, true)
		}

		c.Assert(len(mcfg.Mounts), qt.Equals, 1)
		c.Assert(len(mcfg.Imports), qt.Equals, 1)
		imp := mcfg.Imports[0]
		imp.Path = "github.com/bep/mycomponent"
		c.Assert(imp.Mounts[1].Source, qt.Equals, "src/markdown/blog")
		c.Assert(imp.Mounts[1].Target, qt.Equals, "content/blog")
		c.Assert(imp.Mounts[1].Lang, qt.Equals, "en")
	})

	c.Run("Replacements", func(c *qt.C) {
		for _, tomlConfig := range []string{`
[module]
replacements="a->b,github.com/bep/mycomponent->c"
[[module.imports]]
path="github.com/bep/mycomponent"
`, `
[module]
replacements=["a->b","github.com/bep/mycomponent->c"]
[[module.imports]]
path="github.com/bep/mycomponent"
`} {

			cfg, err := config.FromConfigString(tomlConfig, "toml")
			c.Assert(err, qt.IsNil)

			mcfg, err := DecodeConfig(cfg)
			c.Assert(err, qt.IsNil)
			c.Assert(mcfg.Replacements, qt.DeepEquals, []string{"a->b", "github.com/bep/mycomponent->c"})
			c.Assert(mcfg.replacementsMap, qt.DeepEquals, map[string]string{
				"a":                          "b",
				"github.com/bep/mycomponent": "c",
			})

			c.Assert(mcfg.Imports[0].Path, qt.Equals, "c")

		}
	})
}

func TestDecodeConfigBothOldAndNewProvided(t *testing.T) {
	c := qt.New(t)
	tomlConfig := `

theme = ["b", "c"]

[module]
[[module.imports]]
path="a"

`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	modCfg, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)
	c.Assert(len(modCfg.Imports), qt.Equals, 3)
	c.Assert(modCfg.Imports[0].Path, qt.Equals, "a")
}

// Test old style theme import.
func TestDecodeConfigTheme(t *testing.T) {
	c := qt.New(t)
	tomlConfig := `

theme = ["a", "b"]
`
	cfg, err := config.FromConfigString(tomlConfig, "toml")
	c.Assert(err, qt.IsNil)

	mcfg, err := DecodeConfig(cfg)
	c.Assert(err, qt.IsNil)

	c.Assert(len(mcfg.Imports), qt.Equals, 2)
	c.Assert(mcfg.Imports[0].Path, qt.Equals, "a")
	c.Assert(mcfg.Imports[1].Path, qt.Equals, "b")
}

type fileInfoMock struct {
	os.FileInfo
	name string
	size int64
}

func (mock fileInfoMock) Name() string {
	return mock.name
}

func (mock fileInfoMock) Size() int64 {
	return mock.size
}

func TestMountIsStaticfile(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		name     string
		rule     staticFileRule
		source   string
		target   string
		fi       fileInfoMock
		expected bool
	}{
		{
			"default:static",
			staticFileRule{},
			"static",
			files.ComponentFolderStatic,
			fileInfoMock{name: "/hugo.png"},
			true,
		},
		{
			"default:content:md",
			staticFileRule{},
			"content",
			files.ComponentFolderContent,
			fileInfoMock{name: "/foo/post.md"},
			false,
		},
		{
			"default:content:png",
			staticFileRule{},
			"content",
			files.ComponentFolderContent,
			fileInfoMock{name: "/foo/image.png"},
			true,
		},
		{
			"default:other dirs",
			staticFileRule{},
			"layouts",
			files.ComponentFolderLayouts,
			fileInfoMock{name: "/foo/image.png"},
			false,
		},
		{
			"ignore by regexp",
			staticFileRule{
				ignores: []*regexp.Regexp{
					regexp.MustCompile("\\.md$"),
					regexp.MustCompile("^content/non-static/.*$"),
				},
			},
			"content",
			files.ComponentFolderContent,
			fileInfoMock{name: "/non-static/hugo.mkdwn"},
			false,
		},
		{
			"not ignore by regexp",
			staticFileRule{
				ignores: []*regexp.Regexp{
					regexp.MustCompile("\\.md$"),
					regexp.MustCompile("^content/non-static/.*$"),
				},
			},
			"content",
			files.ComponentFolderContent,
			fileInfoMock{name: "/hugo.png"},
			true,
		},
		{
			"include by file size",
			staticFileRule{
				size: 500,
			},
			"content",
			files.ComponentFolderContent,
			fileInfoMock{name: "/foo/post.md", size: 500 * 1024},
			true,
		},
		{
			"exclude by file size",
			staticFileRule{
				size: 500,
			},
			"content",
			files.ComponentFolderContent,
			fileInfoMock{name: "/foo/post.md", size: 499 * 1024},
			false,
		},
		{
			"ignore rule is prior to size rule",
			staticFileRule{
				ignores: []*regexp.Regexp{
					regexp.MustCompile("\\.md$"),
				},
				size: 100,
			},
			"content",
			files.ComponentFolderContent,
			fileInfoMock{name: "/foo/post.md", size: 500 * 1024},
			false,
		},
	}

	for _, test := range tests {
		m := Mount{
			Source:         test.source,
			Target:         test.target,
			staticFileRule: test.rule,
		}
		got := m.IsStaticFile(test.fi)
		c.Assert(got, qt.Equals, test.expected)
	}
}

func TestGetStaticFileRule(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		name        string
		config      map[string]interface{}
		dirKey      string
		shouldError bool
		assert      func(got staticFileRule)
	}{
		{
			"default",
			map[string]interface{}{},
			"content",
			false,
			func(got staticFileRule) {
				c.Assert(got.size, qt.Equals, float64(0))
				c.Assert(len(got.ignores), qt.Equals, 0)
			},
		},
		{
			"normal",
			map[string]interface{}{
				"content": map[string]interface{}{
					"size": 1.5,
					"ignore": []string{
						".*",
					},
				},
			},
			"content",
			false,
			func(got staticFileRule) {
				c.Assert(got.size, qt.Equals, float64(1.5))
				c.Assert(len(got.ignores), qt.Equals, 1)
				c.Assert(got.ignores[0].MatchString("any"), qt.IsTrue)
			},
		},
		{
			"diffrent dirKey",
			map[string]interface{}{
				"content": map[string]interface{}{
					"size":   0,
					"ignore": []string{},
				},
			},
			"static",
			false,
			func(got staticFileRule) {
				c.Assert(got.size, qt.Equals, float64(0))
				c.Assert(len(got.ignores), qt.Equals, 0)
			},
		},
		{
			"decode error",
			map[string]interface{}{
				"content": map[string]interface{}{
					"size": "foo",
				},
			},
			"content",
			true,
			func(got staticFileRule) {},
		},
		{
			"regex compile error",
			map[string]interface{}{
				"content": map[string]interface{}{
					"ignore": []string{
						"(",
					},
				},
			},
			"content",
			true,
			func(got staticFileRule) {},
		},
	}

	for _, test := range tests {
		v := viper.New()
		v.Set("staticFileRules", test.config)
		got, err := getStaticFileRule(v, test.dirKey)
		if test.shouldError {
			c.Assert(err, qt.Not(qt.IsNil))
		} else {
			test.assert(got)
		}
	}
}
