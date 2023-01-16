// Copyright 2021 The Hugo Authors. All rights reserved.
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

package create

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestDecodeRemoteOptions(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	for _, test := range []struct {
		name    string
		args    map[string]any
		want    fromRemoteOptions
		wantErr bool
	}{
		{
			"POST",
			map[string]any{
				"meThod": "PoST",
				"headers": map[string]any{
					"foo": "bar",
				},
			},
			fromRemoteOptions{
				Method: "POST",
				Headers: map[string]any{
					"foo": "bar",
				},
			},
			false,
		},
		{
			"Body",
			map[string]any{
				"meThod": "POST",
				"body":   []byte("foo"),
			},
			fromRemoteOptions{
				Method: "POST",
				Body:   []byte("foo"),
			},
			false,
		},
		{
			"Body, string",
			map[string]any{
				"meThod": "POST",
				"body":   "foo",
			},
			fromRemoteOptions{
				Method: "POST",
				Body:   []byte("foo"),
			},
			false,
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			got, err := decodeRemoteOptions(test.args)
			isErr := qt.IsNil
			if test.wantErr {
				isErr = qt.IsNotNil
			}

			c.Assert(err, isErr)
			c.Assert(got, qt.DeepEquals, test.want)
		})

	}
}

func TestOptionsNewRequest(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	opts := fromRemoteOptions{
		Method: "GET",
		Body:   []byte("foo"),
	}

	req, err := opts.NewRequest("https://example.com/api")

	c.Assert(err, qt.IsNil)
	c.Assert(req.Method, qt.Equals, "GET")
	c.Assert(req.Header["User-Agent"], qt.DeepEquals, []string{"Hugo Static Site Generator"})

	opts = fromRemoteOptions{
		Method: "GET",
		Body:   []byte("foo"),
		Headers: map[string]any{
			"User-Agent": "foo",
		},
	}

	req, err = opts.NewRequest("https://example.com/api")

	c.Assert(err, qt.IsNil)
	c.Assert(req.Method, qt.Equals, "GET")
	c.Assert(req.Header["User-Agent"], qt.DeepEquals, []string{"foo"})

}

func TestCalculateResourceID(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	c.Assert(calculateResourceID("foo", nil), qt.Equals, "5917621528921068675")
	c.Assert(calculateResourceID("foo", map[string]any{"bar": "baz"}), qt.Equals, "7294498335241413323")

	c.Assert(calculateResourceID("foo", map[string]any{"key": "1234", "bar": "baz"}), qt.Equals, "14904296279238663669")
	c.Assert(calculateResourceID("asdf", map[string]any{"key": "1234", "bar": "asdf"}), qt.Equals, "14904296279238663669")
	c.Assert(calculateResourceID("asdf", map[string]any{"key": "12345", "bar": "asdf"}), qt.Equals, "12191037851845371770")
}
