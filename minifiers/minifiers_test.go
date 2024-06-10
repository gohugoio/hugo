// Copyright 2018 The Hugo Authors. All rights reserved.
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

package minifiers_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/minifiers"
	"github.com/gohugoio/hugo/output"
	"github.com/spf13/afero"
	"github.com/tdewolff/minify/v2/html"
)

func TestNew(t *testing.T) {
	c := qt.New(t)
	m, _ := minifiers.New(media.DefaultTypes, output.DefaultFormats, testconfig.GetTestConfig(afero.NewMemMapFs(), nil))

	var rawJS string
	var minJS string
	rawJS = " var  foo =1 ;   foo ++  ;  "
	minJS = "var foo=1;foo++"

	var rawJSON string
	var minJSON string
	rawJSON = "  { \"a\" : 123 , \"b\":2,  \"c\": 5 } "
	minJSON = "{\"a\":123,\"b\":2,\"c\":5}"

	for _, test := range []struct {
		tp                media.Type
		rawString         string
		expectedMinString string
	}{
		{media.Builtin.CSSType, " body { color: blue; }  ", "body{color:blue}"},
		{media.Builtin.RSSType, " <hello>  Hugo!   </hello>  ", "<hello>Hugo!</hello>"}, // RSS should be handled as XML
		{media.Builtin.JSONType, rawJSON, minJSON},
		{media.Builtin.JavascriptType, rawJS, minJS},
		// JS Regex minifiers
		{media.Type{Type: "application/ecmascript", MainType: "application", SubType: "ecmascript"}, rawJS, minJS},
		{media.Type{Type: "application/javascript", MainType: "application", SubType: "javascript"}, rawJS, minJS},
		{media.Type{Type: "application/x-javascript", MainType: "application", SubType: "x-javascript"}, rawJS, minJS},
		{media.Type{Type: "application/x-ecmascript", MainType: "application", SubType: "x-ecmascript"}, rawJS, minJS},
		{media.Type{Type: "text/ecmascript", MainType: "text", SubType: "ecmascript"}, rawJS, minJS},
		{media.Type{Type: "application/javascript", MainType: "text", SubType: "javascript"}, rawJS, minJS},
		// JSON Regex minifiers
		{media.Type{Type: "application/json", MainType: "application", SubType: "json"}, rawJSON, minJSON},
		{media.Type{Type: "application/x-json", MainType: "application", SubType: "x-json"}, rawJSON, minJSON},
		{media.Type{Type: "application/ld+json", MainType: "application", SubType: "ld+json"}, rawJSON, minJSON},
		{media.Type{Type: "application/json", MainType: "text", SubType: "json"}, rawJSON, minJSON},
	} {
		var b bytes.Buffer

		c.Assert(m.Minify(test.tp, &b, strings.NewReader(test.rawString)), qt.IsNil)
		c.Assert(b.String(), qt.Equals, test.expectedMinString)
	}
}

func TestConfigureMinify(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	v.Set("minify", map[string]any{
		"disablexml": true,
		"tdewolff": map[string]any{
			"html": map[string]any{
				"keepwhitespace": true,
			},
		},
	})
	m, _ := minifiers.New(media.DefaultTypes, output.DefaultFormats, testconfig.GetTestConfig(afero.NewMemMapFs(), v))

	for _, test := range []struct {
		tp                media.Type
		rawString         string
		expectedMinString string
		errorExpected     bool
	}{
		{media.Builtin.HTMLType, "<hello> Hugo! </hello>", "<hello> Hugo! </hello>", false},            // configured minifier
		{media.Builtin.CSSType, " body { color: blue; }  ", "body{color:blue}", false},                 // default minifier
		{media.Builtin.XMLType, " <hello>  Hugo!   </hello>  ", " <hello>  Hugo!   </hello>  ", false}, // disable Xml minification
	} {
		var b bytes.Buffer
		if !test.errorExpected {
			c.Assert(m.Minify(test.tp, &b, strings.NewReader(test.rawString)), qt.IsNil)
			c.Assert(b.String(), qt.Equals, test.expectedMinString)
		} else {
			err := m.Minify(test.tp, &b, strings.NewReader(test.rawString))
			c.Assert(err, qt.ErrorMatches, "minifier does not exist for mimetype")
		}
	}
}

func TestJSONRoundTrip(t *testing.T) {
	c := qt.New(t)
	m, _ := minifiers.New(media.DefaultTypes, output.DefaultFormats, testconfig.GetTestConfig(nil, nil))

	for _, test := range []string{`{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            }
        }
    }
}`} {

		var b bytes.Buffer
		m1 := make(map[string]any)
		m2 := make(map[string]any)
		c.Assert(json.Unmarshal([]byte(test), &m1), qt.IsNil)
		c.Assert(m.Minify(media.Builtin.JSONType, &b, strings.NewReader(test)), qt.IsNil)
		c.Assert(json.Unmarshal(b.Bytes(), &m2), qt.IsNil)
		c.Assert(m1, qt.DeepEquals, m2)
	}
}

func TestBugs(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	m, _ := minifiers.New(media.DefaultTypes, output.DefaultFormats, testconfig.GetTestConfig(nil, v))

	for _, test := range []struct {
		tp                media.Type
		rawString         string
		expectedMinString string
	}{
		// https://github.com/gohugoio/hugo/issues/5506
		{media.Builtin.CSSType, " body { color: rgba(000, 000, 000, 0.7); }", "body{color:rgba(0,0,0,.7)}"},
		// https://github.com/gohugoio/hugo/issues/8332
		{media.Builtin.HTMLType, "<i class='fas fa-tags fa-fw'></i> Tags", `<i class='fas fa-tags fa-fw'></i> Tags`},
	} {
		var b bytes.Buffer

		c.Assert(m.Minify(test.tp, &b, strings.NewReader(test.rawString)), qt.IsNil)
		c.Assert(b.String(), qt.Equals, test.expectedMinString)
	}
}

// Renamed to Precision in v2.7.0. Check that we support both.
func TestDecodeConfigDecimalIsNowPrecision(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	v.Set("minify", map[string]any{
		"disablexml": true,
		"tdewolff": map[string]any{
			"css": map[string]any{
				"decimal": 3,
			},
			"svg": map[string]any{
				"decimal": 3,
			},
		},
	})

	conf := testconfig.GetTestConfigs(nil, v).Base.Minify

	c.Assert(conf.Tdewolff.CSS.Precision, qt.Equals, 3)
}

// Issue 9456
func TestDecodeConfigKeepWhitespace(t *testing.T) {
	c := qt.New(t)
	v := config.New()
	v.Set("minify", map[string]any{
		"tdewolff": map[string]any{
			"html": map[string]any{
				"keepEndTags": false,
			},
		},
	})

	conf := testconfig.GetTestConfigs(nil, v).Base.Minify

	c.Assert(conf.Tdewolff.HTML, qt.DeepEquals,
		html.Minifier{
			KeepComments:        false,
			KeepSpecialComments: true,
			KeepDefaultAttrVals: true,
			KeepDocumentTags:    true,
			KeepEndTags:         false,
			KeepQuotes:          false,
			KeepWhitespace:      false,
		},
	)
}
