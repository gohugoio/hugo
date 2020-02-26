// Copyright 2020 The Hugo Authors. All rights reserved.
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

package postcss

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

// Issue 6166
func TestDecodeOptions(t *testing.T) {
	c := qt.New(t)
	opts1, err := DecodeOptions(map[string]interface{}{
		"no-map": true,
	})

	c.Assert(err, qt.IsNil)
	c.Assert(opts1.NoMap, qt.Equals, true)

	opts2, err := DecodeOptions(map[string]interface{}{
		"noMap": true,
	})

	c.Assert(err, qt.IsNil)
	c.Assert(opts2.NoMap, qt.Equals, true)

}

func TestShouldImport(t *testing.T) {
	c := qt.New(t)

	for _, test := range []struct {
		input  string
		expect bool
	}{
		{input: `@import "navigation.css";`, expect: true},
		{input: `@import "navigation.css"; /* Using a string */`, expect: true},
		{input: `@import "navigation.css"`, expect: true},
		{input: `@import 'navigation.css';`, expect: true},
		{input: `@import url("navigation.css");`, expect: false},
		{input: `@import url('https://fonts.googleapis.com/css?family=Open+Sans:400,400i,800,800i&display=swap');`, expect: false},
	} {
		c.Assert(shouldImport(test.input), qt.Equals, test.expect)
	}
}
