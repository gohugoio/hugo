// Copyright 2024 The Hugo Authors. All rights reserved.
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

package markup_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/markup"
	"github.com/gohugoio/hugo/markup/converter"
)

func TestConverterRegistry(t *testing.T) {
	c := qt.New(t)
	conf := testconfig.GetTestConfig(nil, nil)
	r, err := markup.NewConverterProvider(converter.ProviderConfig{Conf: conf})

	c.Assert(err, qt.IsNil)
	c.Assert("goldmark", qt.Equals, r.GetMarkupConfig().DefaultMarkdownHandler)

	checkName := func(name string) {
		p := r.Get(name)
		c.Assert(p, qt.Not(qt.IsNil))
		c.Assert(p.Name(), qt.Equals, name)
	}

	c.Assert(r.Get("foo"), qt.IsNil)
	c.Assert(r.Get("markdown").Name(), qt.Equals, "goldmark")

	checkName("goldmark")
	checkName("asciidocext")
	checkName("rst")
	checkName("pandoc")
	checkName("org")
}
