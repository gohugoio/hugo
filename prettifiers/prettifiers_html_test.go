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

package prettifiers

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/media"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/output"
	"github.com/spf13/viper"
)

func TestPrettifyUncondensedHTML(t *testing.T) {
	v := viper.New()
	v.Set("prettify", map[string]interface{}{
		"html": map[string]interface{}{
			"condense":   false,
			"inlinetags": []string{"inline"},
		},
	})
	m, _ := New(media.DefaultTypes, output.DefaultFormats, v)

	for _, test := range htmlTable {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			c := qt.New(t)
			want := test.uncondensed[1:] // Strip initial newline, it's only for formatting.
			c.Assert(m.Prettify(media.HTMLType, &b, strings.NewReader(test.input)), qt.IsNil)
			c.Assert(b.String(), qt.Equals, want)
		})
	}
}

func TestPrettifyCondensedHTML(t *testing.T) {
	v := viper.New()
	v.Set("prettify", map[string]interface{}{
		"html": map[string]interface{}{
			"condense":   true,
			"inlinetags": []string{"inline"},
		},
	})
	m, _ := New(media.DefaultTypes, output.DefaultFormats, v)

	for _, test := range htmlTable {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			c := qt.New(t)
			want := test.condensed[1:] // Strip initial newline, it's only for formatting.
			c.Assert(m.Prettify(media.HTMLType, &b, strings.NewReader(test.input)), qt.IsNil)
			c.Assert(b.String(), qt.Equals, want)
		})
	}
}

var htmlTable = []struct {
	name        string
	input       string
	condensed   string
	uncondensed string
}{
	{
		name: "basic",
		input: `
<html><body><h1>
 Hugo!
 </h1></body> </html>`,
		condensed: `
<html>
  <body>
    <h1> Hugo! </h1>
  </body>
</html>
`,
		uncondensed: `
<html>
  <body>
    <h1>
      Hugo!
    </h1>
  </body>
</html>
`},
	{
		name: "inline",
		input: `
<html> <body><inline><inline>
 Hugo!
 </inline></inline> </html>`,
		condensed: `
<html>
  <body>
  <inline><inline> Hugo! </inline></inline>
</html>
`,
		uncondensed: `
<html>
  <body>
  <inline>
    <inline>
      Hugo!
    </inline>
  </inline>
</html>
`},
	{
		name: "block",
		input: `
<html> <body><block><block>
 Hugo!
 </block></block> </html>`,
		condensed: `
<html>
  <body>
  <block>
    <block> Hugo! </block>
  </block>
</html>
`,
		uncondensed: `
<html>
  <body>
  <block>
    <block>
      Hugo!
    </block>
  </block>
</html>
`},
	{
		name: "empty-lines",
		input: `
<html>

<div>Hugo!</div>


</html>`,
		condensed: `
<html>
  <div>Hugo!</div>
</html>
`,
		uncondensed: `
<html>
  <div>
    Hugo!
  </div>
</html>
`},
}
