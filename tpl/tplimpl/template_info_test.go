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
package tplimpl

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/tpl"
)

func TestTemplateInfoShortcode(t *testing.T) {
	c := qt.New(t)
	d := newD(c)
	h := d.Tmpl().(*templateExec)

	c.Assert(h.AddTemplate("shortcodes/mytemplate.html", `
{{ .Inner }}
`), qt.IsNil)

	c.Assert(h.postTransform(), qt.IsNil)

	tt, found, _ := d.Tmpl().LookupVariant("mytemplate", tpl.TemplateVariants{})

	c.Assert(found, qt.Equals, true)
	tti, ok := tt.(tpl.Info)
	c.Assert(ok, qt.Equals, true)
	c.Assert(tti.ParseInfo().IsInner, qt.Equals, true)

}

// TODO(bep) move and use in other places
func newD(c *qt.C) *deps.Deps {
	v := newTestConfig()
	fs := hugofs.NewMem(v)

	depsCfg := newDepsConfig(v)
	depsCfg.Fs = fs
	d, err := deps.New(depsCfg)
	c.Assert(err, qt.IsNil)

	provider := DefaultTemplateProvider
	provider.Update(d)

	return d

}
