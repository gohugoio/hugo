// Copyright 2025 The Hugo Authors. All rights reserved.
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

package plaintext

import (
	"testing"

	"github.com/gohugoio/hugo/markup/converter"

	qt "github.com/frankban/quicktest"
)

func TestConvert(t *testing.T) {
	c := qt.New(t)

	p, err := Provider.New(converter.ProviderConfig{})
	c.Assert(err, qt.IsNil)

	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)

	input := "Just some plain text and some # markdown and some <p>html content</p>."
	b, err := conv.Convert(converter.RenderContext{Src: []byte(input)})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, input)
}
