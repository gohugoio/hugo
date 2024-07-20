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

package transform_test

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl/transform"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/media"

	qt "github.com/frankban/quicktest"
)

const (
	testJSON = `
	
{
    "ROOT_KEY": {
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
}

	`
)

var _ resource.ReadSeekCloserResource = (*testContentResource)(nil)

type testContentResource struct {
	content string
	mime    media.Type

	key string
}

func (t testContentResource) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	return hugio.NewReadSeekerNoOpCloserFromString(t.content), nil
}

func (t testContentResource) MediaType() media.Type {
	return t.mime
}

func (t testContentResource) Key() string {
	return t.key
}

func TestUnmarshal(t *testing.T) {
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t},
	).Build()

	ns := transform.New(b.H.Deps)

	assertSlogan := func(m map[string]any) {
		b.Assert(m["slogan"], qt.Equals, "Hugo Rocks!")
	}

	for _, test := range []struct {
		data    any
		options any
		expect  any
	}{
		{`{ "slogan": "Hugo Rocks!" }`, nil, func(m map[string]any) {
			assertSlogan(m)
		}},
		{`slogan: "Hugo Rocks!"`, nil, func(m map[string]any) {
			assertSlogan(m)
		}},
		{`slogan = "Hugo Rocks!"`, nil, func(m map[string]any) {
			assertSlogan(m)
		}},
		{testContentResource{key: "r1", content: `slogan: "Hugo Rocks!"`, mime: media.Builtin.YAMLType}, nil, func(m map[string]any) {
			assertSlogan(m)
		}},
		{testContentResource{key: "r1", content: `{ "slogan": "Hugo Rocks!" }`, mime: media.Builtin.JSONType}, nil, func(m map[string]any) {
			assertSlogan(m)
		}},
		{testContentResource{key: "r1", content: `slogan = "Hugo Rocks!"`, mime: media.Builtin.TOMLType}, nil, func(m map[string]any) {
			assertSlogan(m)
		}},
		{testContentResource{key: "r1", content: `<root><slogan>Hugo Rocks!</slogan></root>"`, mime: media.Builtin.XMLType}, nil, func(m map[string]any) {
			assertSlogan(m)
		}},
		{testContentResource{key: "r1", content: `1997,Ford,E350,"ac, abs, moon",3000.00
1999,Chevy,"Venture ""Extended Edition""","",4900.00`, mime: media.Builtin.CSVType}, nil, func(r [][]string) {
			b.Assert(len(r), qt.Equals, 2)
			first := r[0]
			b.Assert(len(first), qt.Equals, 5)
			b.Assert(first[1], qt.Equals, "Ford")
		}},
		{testContentResource{key: "r1", content: `a;b;c`, mime: media.Builtin.CSVType}, map[string]any{"delimiter": ";"}, func(r [][]string) {
			b.Assert([][]string{{"a", "b", "c"}}, qt.DeepEquals, r)
		}},
		{"a,b,c", nil, func(r [][]string) {
			b.Assert([][]string{{"a", "b", "c"}}, qt.DeepEquals, r)
		}},
		{"a;b;c", map[string]any{"delimiter": ";"}, func(r [][]string) {
			b.Assert([][]string{{"a", "b", "c"}}, qt.DeepEquals, r)
		}},
		{testContentResource{key: "r1", content: `
% This is a comment
a;b;c`, mime: media.Builtin.CSVType}, map[string]any{"DElimiter": ";", "Comment": "%"}, func(r [][]string) {
			b.Assert([][]string{{"a", "b", "c"}}, qt.DeepEquals, r)
		}},
		// errors
		{"thisisnotavaliddataformat", nil, false},
		{testContentResource{key: "r1", content: `invalid&toml"`, mime: media.Builtin.TOMLType}, nil, false},
		{testContentResource{key: "r1", content: `unsupported: MIME"`, mime: media.Builtin.CalendarType}, nil, false},
		{"thisisnotavaliddataformat", nil, false},
		{`{ notjson }`, nil, false},
		{tstNoStringer{}, nil, false},
	} {

		ns.Reset()

		var args []any

		if test.options != nil {
			args = []any{test.options, test.data}
		} else {
			args = []any{test.data}
		}

		result, err := ns.Unmarshal(args...)

		if bb, ok := test.expect.(bool); ok && !bb {
			b.Assert(err, qt.Not(qt.IsNil))
		} else if fn, ok := test.expect.(func(m map[string]any)); ok {
			b.Assert(err, qt.IsNil)
			m, ok := result.(map[string]any)
			b.Assert(ok, qt.Equals, true)
			fn(m)
		} else if fn, ok := test.expect.(func(r [][]string)); ok {
			b.Assert(err, qt.IsNil)
			r, ok := result.([][]string)
			b.Assert(ok, qt.Equals, true)
			fn(r)
		} else {
			b.Assert(err, qt.IsNil)
			b.Assert(result, qt.Equals, test.expect)
		}

	}
}

func BenchmarkUnmarshalString(b *testing.B) {
	bb := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: b},
	).Build()

	ns := transform.New(bb.H.Deps)

	const numJsons = 100

	var jsons [numJsons]string
	for i := 0; i < numJsons; i++ {
		jsons[i] = strings.Replace(testJSON, "ROOT_KEY", fmt.Sprintf("root%d", i), 1)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := ns.Unmarshal(ctx, jsons[rand.Intn(numJsons)])
		if err != nil {
			b.Fatal(err)
		}
		if result == nil {
			b.Fatal("no result")
		}
	}
}

func BenchmarkUnmarshalResource(b *testing.B) {
	bb := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: b},
	).Build()

	ns := transform.New(bb.H.Deps)

	const numJsons = 100

	var jsons [numJsons]testContentResource
	for i := 0; i < numJsons; i++ {
		key := fmt.Sprintf("root%d", i)
		jsons[i] = testContentResource{key: key, content: strings.Replace(testJSON, "ROOT_KEY", key, 1), mime: media.Builtin.JSONType}
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := ns.Unmarshal(ctx, jsons[rand.Intn(numJsons)])
		if err != nil {
			b.Fatal(err)
		}
		if result == nil {
			b.Fatal("no result")
		}
	}
}
