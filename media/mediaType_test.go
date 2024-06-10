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

package media

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/paths"
)

func TestGetByType(t *testing.T) {
	c := qt.New(t)

	types := DefaultTypes

	mt, found := types.GetByType("text/HTML")
	c.Assert(found, qt.Equals, true)
	c.Assert(mt.SubType, qt.Equals, "html")

	_, found = types.GetByType("text/nono")
	c.Assert(found, qt.Equals, false)

	mt, found = types.GetByType("application/rss+xml")
	c.Assert(found, qt.Equals, true)
	c.Assert(mt.SubType, qt.Equals, "rss")

	mt, found = types.GetByType("application/rss")
	c.Assert(found, qt.Equals, true)
	c.Assert(mt.SubType, qt.Equals, "rss")
}

func TestGetByMainSubType(t *testing.T) {
	c := qt.New(t)
	f, found := DefaultTypes.GetByMainSubType("text", "plain")
	c.Assert(found, qt.Equals, true)
	c.Assert(f.SubType, qt.Equals, "plain")
	_, found = DefaultTypes.GetByMainSubType("foo", "plain")
	c.Assert(found, qt.Equals, false)
}

func TestBySuffix(t *testing.T) {
	c := qt.New(t)
	formats := DefaultTypes.BySuffix("xml")
	c.Assert(len(formats), qt.Equals, 2)
	c.Assert(formats[0].SubType, qt.Equals, "rss")
	c.Assert(formats[1].SubType, qt.Equals, "xml")
}

func TestGetFirstBySuffix(t *testing.T) {
	c := qt.New(t)

	types := make(Types, len(DefaultTypes))
	copy(types, DefaultTypes)

	// Issue #8406
	geoJSON := newMediaTypeWithMimeSuffix("application", "geo", "json", []string{"geojson", "gjson"})
	types = append(types, geoJSON)
	sort.Sort(types)

	check := func(suffix string, expectedType Type) {
		t, f, found := types.GetFirstBySuffix(suffix)
		c.Assert(found, qt.Equals, true)
		c.Assert(f, qt.Equals, SuffixInfo{
			Suffix:     suffix,
			FullSuffix: "." + suffix,
		})
		c.Assert(t, qt.Equals, expectedType)
	}

	check("js", Builtin.JavascriptType)
	check("json", Builtin.JSONType)
	check("geojson", geoJSON)
	check("gjson", geoJSON)
}

func TestFromTypeString(t *testing.T) {
	c := qt.New(t)
	f, err := FromString("text/html")
	c.Assert(err, qt.IsNil)
	c.Assert(f.Type, qt.Equals, Builtin.HTMLType.Type)

	f, err = FromString("application/custom")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, Type{Type: "application/custom", MainType: "application", SubType: "custom", mimeSuffix: ""})

	f, err = FromString("application/custom+sfx")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, Type{Type: "application/custom+sfx", MainType: "application", SubType: "custom", mimeSuffix: "sfx"})

	_, err = FromString("noslash")
	c.Assert(err, qt.Not(qt.IsNil))

	f, err = FromString("text/xml; charset=utf-8")
	c.Assert(err, qt.IsNil)

	c.Assert(f, qt.Equals, Type{Type: "text/xml", MainType: "text", SubType: "xml", mimeSuffix: ""})
}

func TestFromStringAndExt(t *testing.T) {
	c := qt.New(t)
	f, err := FromStringAndExt("text/html", "html", "htm")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, Builtin.HTMLType)
	f, err = FromStringAndExt("text/html", ".html", ".htm")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, Builtin.HTMLType)
}

// Add a test for the SVG case
// https://github.com/gohugoio/hugo/issues/4920
func TestFromExtensionMultipleSuffixes(t *testing.T) {
	c := qt.New(t)
	tp, si, found := DefaultTypes.GetBySuffix("svg")
	c.Assert(found, qt.Equals, true)
	c.Assert(tp.String(), qt.Equals, "image/svg+xml")
	c.Assert(si.Suffix, qt.Equals, "svg")
	c.Assert(si.FullSuffix, qt.Equals, ".svg")
	c.Assert(tp.FirstSuffix.Suffix, qt.Equals, si.Suffix)
	c.Assert(tp.FirstSuffix.FullSuffix, qt.Equals, si.FullSuffix)
	ftp, found := DefaultTypes.GetByType("image/svg+xml")
	c.Assert(found, qt.Equals, true)
	c.Assert(ftp.String(), qt.Equals, "image/svg+xml")
	c.Assert(found, qt.Equals, true)
}

func TestFromContent(t *testing.T) {
	c := qt.New(t)

	files, err := filepath.Glob("./testdata/resource.*")
	c.Assert(err, qt.IsNil)

	for _, filename := range files {
		name := filepath.Base(filename)
		c.Run(name, func(c *qt.C) {
			content, err := os.ReadFile(filename)
			c.Assert(err, qt.IsNil)
			ext := strings.TrimPrefix(paths.Ext(filename), ".")
			var exts []string
			if ext == "jpg" {
				exts = append(exts, "foo", "bar", "jpg")
			} else {
				exts = []string{ext}
			}
			expected, _, found := DefaultTypes.GetFirstBySuffix(ext)
			c.Assert(found, qt.IsTrue)
			got := FromContent(DefaultTypes, exts, content)
			c.Assert(got, qt.Equals, expected)
		})
	}
}

func TestFromContentFakes(t *testing.T) {
	c := qt.New(t)

	files, err := filepath.Glob("./testdata/fake.*")
	c.Assert(err, qt.IsNil)

	for _, filename := range files {
		name := filepath.Base(filename)
		c.Run(name, func(c *qt.C) {
			content, err := os.ReadFile(filename)
			c.Assert(err, qt.IsNil)
			ext := strings.TrimPrefix(paths.Ext(filename), ".")
			got := FromContent(DefaultTypes, []string{ext}, content)
			c.Assert(got, qt.Equals, zero)
		})
	}
}

func TestToJSON(t *testing.T) {
	c := qt.New(t)
	b, err := json.Marshal(Builtin.MPEGType)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, `{"mainType":"video","subType":"mpeg","delimiter":".","type":"video/mpeg","string":"video/mpeg","suffixes":["mpg","mpeg"]}`)
}

func BenchmarkTypeOps(b *testing.B) {
	mt := Builtin.MPEGType
	mts := DefaultTypes
	for i := 0; i < b.N; i++ {
		ff := mt.FirstSuffix
		_ = ff.FullSuffix
		_ = mt.IsZero()
		c, err := mt.MarshalJSON()
		if c == nil || err != nil {
			b.Fatal("failed")
		}
		_ = mt.String()
		_ = ff.Suffix
		_ = mt.Suffixes
		_ = mt.Type
		_ = mts.BySuffix("xml")
		_, _ = mts.GetByMainSubType("application", "xml")
		_, _, _ = mts.GetBySuffix("xml")
		_, _ = mts.GetByType("application")
		_, _, _ = mts.GetFirstBySuffix("xml")

	}
}

func TestIsContentFile(t *testing.T) {
	c := qt.New(t)

	c.Assert(DefaultContentTypes.IsContentFile(filepath.FromSlash("my/file.md")), qt.Equals, true)
	c.Assert(DefaultContentTypes.IsContentFile(filepath.FromSlash("my/file.ad")), qt.Equals, true)
	c.Assert(DefaultContentTypes.IsContentFile(filepath.FromSlash("textfile.txt")), qt.Equals, false)
}
