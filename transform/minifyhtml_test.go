// Copyright 2016 The Hugo Authors. All rights reserved.
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

package transform

import (
	"bytes"
	"strings"
	"testing"
)

func TestMinifyHTML(t *testing.T) {
	testCases := []struct {
		in, expected string
	}{
		{"", ""},
		{"<br>    <br>", "<br> <br>"},
		{"<br>\n\n<br>", "<br>\n<br>"},
		{"<p></p>", "<p>"},
		{
			"<style> body { color: black; any: 0s 0px; } </style>",
			"<style>body{color:#000;any:0s 0}</style>",
		},
		{
			"<html><head></head><body></body></html>",
			"<html><head></head><body></body></html>",
		},
		{
			`<script> console.log( "duck" );
			 // foo
			 debugger;
			 </script>`,
			`<script>console.log("duck");debugger;</script>`,
		},
		{
			`<script type="application/json">
			{
				"a": 10,
				"b": [1, 2, 3]
			}
			</script>`,
			`<script type=application/json>{"a":10,"b":[1,2,3]}</script>`,
		},
	}
	tr := NewChain(MinifyHTML)
	for i, tc := range testCases {
		in := strings.NewReader(tc.in)
		out := new(bytes.Buffer)

		tr.Apply(out, in, []byte("path"))
		if string(out.Bytes()) != tc.expected {
			t.Errorf("%d. Expected %q got %q", i, tc.expected, string(out.Bytes()))
		}
	}
}
