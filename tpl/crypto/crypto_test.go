// Copyright 2017 The Hugo Authors. All rights reserved.
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

package crypto

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestMD5(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	ns := New()

	for i, test := range []struct {
		in     interface{}
		expect interface{}
	}{
		{"Hello world, gophers!", "b3029f756f98f79e7f1b7f1d1f0dd53b"},
		{"Lorem ipsum dolor", "06ce65ac476fc656bea3fca5d02cfd81"},
		{t, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.in)

		result, err := ns.MD5(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestSHA1(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for i, test := range []struct {
		in     interface{}
		expect interface{}
	}{
		{"Hello world, gophers!", "c8b5b0e33d408246e30f53e32b8f7627a7a649d4"},
		{"Lorem ipsum dolor", "45f75b844be4d17b3394c6701768daf39419c99b"},
		{t, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.in)

		result, err := ns.SHA1(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}

func TestSHA256(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	ns := New()

	for i, test := range []struct {
		in     interface{}
		expect interface{}
	}{
		{"Hello world, gophers!", "6ec43b78da9669f50e4e422575c54bf87536954ccd58280219c393f2ce352b46"},
		{"Lorem ipsum dolor", "9b3e1beb7053e0f900a674dd1c99aca3355e1275e1b03d3cb1bc977f5154e196"},
		{t, false},
	} {
		errMsg := qt.Commentf("[%d] %v", i, test.in)

		result, err := ns.SHA256(test.in)

		if b, ok := test.expect.(bool); ok && !b {
			c.Assert(err, qt.Not(qt.IsNil), errMsg)
			continue
		}

		c.Assert(err, qt.IsNil, errMsg)
		c.Assert(result, qt.Equals, test.expect, errMsg)
	}
}
