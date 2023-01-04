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

package metadecoders

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestFormatFromString(t *testing.T) {
	c := qt.New(t)
	for _, test := range []struct {
		s      string
		expect Format
	}{
		{"json", JSON},
		{"yaml", YAML},
		{"yml", YAML},
		{"xml", XML},
		{"toml", TOML},
		{"config.toml", TOML},
		{"tOMl", TOML},
		{"org", ORG},
		{"foo", ""},
	} {
		c.Assert(FormatFromString(test.s), qt.Equals, test.expect)
	}
}

func TestFormatFromContentString(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	for i, test := range []struct {
		data   string
		expect any
	}{
		{`foo = "bar"`, TOML},
		{`   foo = "bar"`, TOML},
		{`foo="bar"`, TOML},
		{`foo: "bar"`, YAML},
		{`foo:"bar"`, YAML},
		{`{ "foo": "bar"`, JSON},
		{`a,b,c"`, CSV},
		{`<foo>bar</foo>"`, XML},
		{`asdfasdf`, Format("")},
		{``, Format("")},
	} {
		errMsg := qt.Commentf("[%d] %s", i, test.data)

		result := Default.FormatFromContentString(test.data)

		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}
