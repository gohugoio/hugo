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

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/tpl"
	"github.com/stretchr/testify/require"
)

func TestTemplateInfoShortcode(t *testing.T) {
	assert := require.New(t)
	d := newD(assert)
	h := d.Tmpl.(tpl.TemplateHandler)

	assert.NoError(h.AddTemplate("shortcodes/mytemplate.html", `
{{ .Inner }}
`))
	tt, found, _ := d.Tmpl.LookupVariant("mytemplate", tpl.TemplateVariants{})

	assert.True(found)
	tti, ok := tt.(tpl.TemplateInfoProvider)
	assert.True(ok)
	assert.True(tti.TemplateInfo().IsInner)

}

// TODO(bep) move and use in other places
func newD(assert *require.Assertions) *deps.Deps {
	v := newTestConfig()
	fs := hugofs.NewMem(v)

	depsCfg := newDepsConfig(v)
	depsCfg.Fs = fs
	d, err := deps.New(depsCfg)
	assert.NoError(err)

	provider := DefaultTemplateProvider
	provider.Update(d)

	return d

}
