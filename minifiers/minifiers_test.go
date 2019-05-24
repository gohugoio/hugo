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
	"strings"
	"testing"

	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/output"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert := require.New(t)
	m := New(media.DefaultTypes, output.DefaultFormats)

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

		assert.NoError(m.Minify(test.tp, &b, strings.NewReader(test.rawString)))
		assert.Equal(test.expectedMinString, b.String())
	}

}

func TestBugs(t *testing.T) {
	assert := require.New(t)
	m := New(media.DefaultTypes, output.DefaultFormats)

	for _, test := range []struct {
		tp                media.Type
		rawString         string
		expectedMinString string
	}{
		// https://github.com/gohugoio/hugo/issues/5506
		{media.CSSType, " body { color: rgba(000, 000, 000, 0.7); }", "body{color:rgba(0,0,0,.7)}"},
	} {
		var b bytes.Buffer

		assert.NoError(m.Minify(test.tp, &b, strings.NewReader(test.rawString)))
		assert.Equal(test.expectedMinString, b.String())
	}

}
