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
	c := qt.New(t)

	for _, test := range []struct {
		name    string
		args    map[string]interface{}
		want    fromRemoteOptions
		wantErr bool
	}{
		{
			"POST",
			map[string]interface{}{
				"meThod": "PoST",
				"headers": map[string]interface{}{
					"foo": "bar",
				},
			},
			fromRemoteOptions{
				Method: "POST",
				Headers: map[string]interface{}{
					"foo": "bar",
				},
			},
			false,
		},
		{
			"Body",
			map[string]interface{}{
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
			map[string]interface{}{
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
