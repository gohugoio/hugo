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
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDefaultTypes(t *testing.T) {
	c := qt.New(t)
	for _, test := range []struct {
		tp               Type
		expectedMainType string
		expectedSubType  string
		expectedSuffix   string
		expectedType     string
		expectedString   string
	}{
		{CalendarType, "text", "calendar", "ics", "text/calendar", "text/calendar"},
		{CSSType, "text", "css", "css", "text/css", "text/css"},
		{SCSSType, "text", "x-scss", "scss", "text/x-scss", "text/x-scss"},
		{CSVType, "text", "csv", "csv", "text/csv", "text/csv"},
		{HTMLType, "text", "html", "html", "text/html", "text/html"},
		{JavascriptType, "application", "javascript", "js", "application/javascript", "application/javascript"},
		{TypeScriptType, "application", "typescript", "ts", "application/typescript", "application/typescript"},
		{TSXType, "text", "tsx", "tsx", "text/tsx", "text/tsx"},
		{JSXType, "text", "jsx", "jsx", "text/jsx", "text/jsx"},
		{JSONType, "application", "json", "json", "application/json", "application/json"},
		{RSSType, "application", "rss", "xml", "application/rss+xml", "application/rss+xml"},
		{SVGType, "image", "svg", "svg", "image/svg+xml", "image/svg+xml"},
		{TextType, "text", "plain", "txt", "text/plain", "text/plain"},
		{XMLType, "application", "xml", "xml", "application/xml", "application/xml"},
		{TOMLType, "application", "toml", "toml", "application/toml", "application/toml"},
		{YAMLType, "application", "yaml", "yaml", "application/yaml", "application/yaml"},
	} {
		c.Assert(test.tp.MainType, qt.Equals, test.expectedMainType)
		c.Assert(test.tp.SubType, qt.Equals, test.expectedSubType)

		c.Assert(test.tp.Type(), qt.Equals, test.expectedType)
		c.Assert(test.tp.String(), qt.Equals, test.expectedString)

	}

	c.Assert(len(DefaultTypes), qt.Equals, 27)
}

func TestGetByType(t *testing.T) {
	c := qt.New(t)

	types := Types{HTMLType, RSSType}

	mt, found := types.GetByType("text/HTML")
	c.Assert(found, qt.Equals, true)
	c.Assert(HTMLType, qt.Equals, mt)

	_, found = types.GetByType("text/nono")
	c.Assert(found, qt.Equals, false)

	mt, found = types.GetByType("application/rss+xml")
	c.Assert(found, qt.Equals, true)
	c.Assert(RSSType, qt.Equals, mt)

	mt, found = types.GetByType("application/rss")
	c.Assert(found, qt.Equals, true)
	c.Assert(RSSType, qt.Equals, mt)
}

func TestGetByMainSubType(t *testing.T) {
	c := qt.New(t)
	f, found := DefaultTypes.GetByMainSubType("text", "plain")
	c.Assert(found, qt.Equals, true)
	c.Assert(f, qt.Equals, TextType)
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
	_, f, found := DefaultTypes.GetFirstBySuffix("xml")
	c.Assert(found, qt.Equals, true)
	c.Assert(f, qt.Equals, SuffixInfo{
		Suffix:     "xml",
		FullSuffix: ".xml"})
}

func TestFromTypeString(t *testing.T) {
	c := qt.New(t)
	f, err := fromString("text/html")
	c.Assert(err, qt.IsNil)
	c.Assert(f.Type(), qt.Equals, HTMLType.Type())

	f, err = fromString("application/custom")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, Type{MainType: "application", SubType: "custom", mimeSuffix: ""})

	f, err = fromString("application/custom+sfx")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, Type{MainType: "application", SubType: "custom", mimeSuffix: "sfx"})

	_, err = fromString("noslash")
	c.Assert(err, qt.Not(qt.IsNil))

	f, err = fromString("text/xml; charset=utf-8")
	c.Assert(err, qt.IsNil)

	c.Assert(f, qt.Equals, Type{MainType: "text", SubType: "xml", mimeSuffix: ""})

}

func TestFromStringAndExt(t *testing.T) {
	c := qt.New(t)
	f, err := FromStringAndExt("text/html", "html")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, HTMLType)
	f, err = FromStringAndExt("text/html", ".html")
	c.Assert(err, qt.IsNil)
	c.Assert(f, qt.Equals, HTMLType)
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

func TestDecodeTypes(t *testing.T) {
	c := qt.New(t)

	tests := []struct {
		name        string
		maps        []map[string]interface{}
		shouldError bool
		assert      func(t *testing.T, name string, tt Types)
	}{
		{
			"Redefine JSON",
			[]map[string]interface{}{
				{
					"application/json": map[string]interface{}{
						"suffixes": []string{"jasn"},
					},
				},
			},
			false,
			func(t *testing.T, name string, tt Types) {
				c.Assert(len(tt), qt.Equals, len(DefaultTypes))
				json, si, found := tt.GetBySuffix("jasn")
				c.Assert(found, qt.Equals, true)
				c.Assert(json.String(), qt.Equals, "application/json")
				c.Assert(si.FullSuffix, qt.Equals, ".jasn")
			},
		},
		{
			"MIME suffix in key, multiple file suffixes, custom delimiter",
			[]map[string]interface{}{
				{
					"application/hugo+hg": map[string]interface{}{
						"suffixes":  []string{"hg1", "hG2"},
						"Delimiter": "_",
					},
				},
			},
			false,
			func(t *testing.T, name string, tt Types) {
				c.Assert(len(tt), qt.Equals, len(DefaultTypes)+1)
				hg, si, found := tt.GetBySuffix("hg2")
				c.Assert(found, qt.Equals, true)
				c.Assert(hg.mimeSuffix, qt.Equals, "hg")
				c.Assert(hg.FirstSuffix.Suffix, qt.Equals, "hg1")
				c.Assert(hg.FirstSuffix.FullSuffix, qt.Equals, "_hg1")
				c.Assert(si.Suffix, qt.Equals, "hg2")
				c.Assert(si.FullSuffix, qt.Equals, "_hg2")
				c.Assert(hg.String(), qt.Equals, "application/hugo+hg")

				_, found = tt.GetByType("application/hugo+hg")
				c.Assert(found, qt.Equals, true)

			},
		},
		{
			"Add custom media type",
			[]map[string]interface{}{
				{
					"text/hugo+hgo": map[string]interface{}{
						"Suffixes": []string{"hgo2"},
					},
				},
			},
			false,
			func(t *testing.T, name string, tp Types) {
				c.Assert(len(tp), qt.Equals, len(DefaultTypes)+1)
				// Make sure we have not broken the default config.

				_, _, found := tp.GetBySuffix("json")
				c.Assert(found, qt.Equals, true)

				hugo, _, found := tp.GetBySuffix("hgo2")
				c.Assert(found, qt.Equals, true)
				c.Assert(hugo.String(), qt.Equals, "text/hugo+hgo")
			},
		},
	}

	for _, test := range tests {
		result, err := DecodeTypes(test.maps...)
		if test.shouldError {
			c.Assert(err, qt.Not(qt.IsNil))
		} else {
			c.Assert(err, qt.IsNil)
			test.assert(t, test.name, result)
		}
	}
}

func TestToJSON(t *testing.T) {
	c := qt.New(t)
	b, err := json.Marshal(MPEGType)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, `{"mainType":"video","subType":"mpeg","delimiter":".","firstSuffix":{"suffix":"mpg","fullSuffix":".mpg"},"type":"video/mpeg","string":"video/mpeg","suffixes":["mpg","mpeg"]}`)
}

func BenchmarkTypeOps(b *testing.B) {
	mt := MPEGType
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
		_ = mt.Type()
		_ = mts.BySuffix("xml")
		_, _ = mts.GetByMainSubType("application", "xml")
		_, _, _ = mts.GetBySuffix("xml")
		_, _ = mts.GetByType("application")
		_, _, _ = mts.GetFirstBySuffix("xml")

	}
}
