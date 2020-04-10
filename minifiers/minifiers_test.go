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

package minifiers

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/media"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/output"
	"github.com/spf13/viper"
)

func TestNew(t *testing.T) {
	c := qt.New(t)
	v := viper.New()
	m, _ := New(media.DefaultTypes, output.DefaultFormats, v)

	var rawJS string
	var minJS string
	rawJS = " var  foo =1 ;   foo ++  ;  "
	minJS = "var foo=1;foo++;"

	var rawJSON string
	var minJSON string
	rawJSON = "  { \"a\" : 123 , \"b\":2,  \"c\": 5 } "
	minJSON = "{\"a\":123,\"b\":2,\"c\":5}"

	for _, test := range []struct {
		tp                media.Type
		rawString         string
		expectedMinString string
	}{
		{media.CSSType, " body { color: blue; }  ", "body{color:blue}"},
		{media.RSSType, " <hello>  Hugo!   </hello>  ", "<hello>Hugo!</hello>"}, // RSS should be handled as XML
		{media.JSONType, rawJSON, minJSON},
		{media.JavascriptType, rawJS, minJS},
		// JS Regex minifiers
		{media.Type{MainType: "application", SubType: "ecmascript"}, rawJS, minJS},
		{media.Type{MainType: "application", SubType: "javascript"}, rawJS, minJS},
		{media.Type{MainType: "application", SubType: "x-javascript"}, rawJS, minJS},
		{media.Type{MainType: "application", SubType: "x-ecmascript"}, rawJS, minJS},
		{media.Type{MainType: "text", SubType: "ecmascript"}, rawJS, minJS},
		{media.Type{MainType: "text", SubType: "javascript"}, rawJS, minJS},
		{media.Type{MainType: "text", SubType: "x-javascript"}, rawJS, minJS},
		{media.Type{MainType: "text", SubType: "x-ecmascript"}, rawJS, minJS},
		// JSON Regex minifiers
		{media.Type{MainType: "application", SubType: "json"}, rawJSON, minJSON},
		{media.Type{MainType: "application", SubType: "x-json"}, rawJSON, minJSON},
		{media.Type{MainType: "application", SubType: "ld+json"}, rawJSON, minJSON},
		{media.Type{MainType: "text", SubType: "json"}, rawJSON, minJSON},
		{media.Type{MainType: "text", SubType: "x-json"}, rawJSON, minJSON},
		{media.Type{MainType: "text", SubType: "ld+json"}, rawJSON, minJSON},
	} {
		var b bytes.Buffer

		c.Assert(m.Minify(test.tp, &b, strings.NewReader(test.rawString)), qt.IsNil)
		c.Assert(b.String(), qt.Equals, test.expectedMinString)
	}

}

func TestConfigureMinify(t *testing.T) {
	c := qt.New(t)
	v := viper.New()
	v.Set("minify", map[string]interface{}{
		"disablexml": true,
		"tdewolff": map[string]interface{}{
			"html": map[string]interface{}{
				"keepwhitespace": true,
			},
		},
	})
	m, _ := New(media.DefaultTypes, output.DefaultFormats, v)

	for _, test := range []struct {
		tp                media.Type
		rawString         string
		expectedMinString string
		errorExpected     bool
	}{
		{media.HTMLType, "<hello> Hugo! </hello>", "<hello> Hugo! </hello>", false}, // configured minifier
		{media.CSSType, " body { color: blue; }  ", "body{color:blue}", false},      // default minifier
		{media.XMLType, " <hello>  Hugo!   </hello>  ", "", true},                   // disable Xml minificatin
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
	v := viper.New()
	m, _ := New(media.DefaultTypes, output.DefaultFormats, v)

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
		m1 := make(map[string]interface{})
		m2 := make(map[string]interface{})
		c.Assert(json.Unmarshal([]byte(test), &m1), qt.IsNil)
		c.Assert(m.Minify(media.JSONType, &b, strings.NewReader(test)), qt.IsNil)
		c.Assert(json.Unmarshal(b.Bytes(), &m2), qt.IsNil)
		c.Assert(m1, qt.DeepEquals, m2)
	}

}

func TestBugs(t *testing.T) {
	c := qt.New(t)
	v := viper.New()
	m, _ := New(media.DefaultTypes, output.DefaultFormats, v)

	for _, test := range []struct {
		tp                media.Type
		rawString         string
		expectedMinString string
	}{
		// https://github.com/gohugoio/hugo/issues/5506
		{media.CSSType, " body { color: rgba(000, 000, 000, 0.7); }", "body{color:rgba(0,0,0,.7)}"},
	} {
		var b bytes.Buffer

		c.Assert(m.Minify(test.tp, &b, strings.NewReader(test.rawString)), qt.IsNil)
		c.Assert(b.String(), qt.Equals, test.expectedMinString)
	}

}
