// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type handler interface {
	addTemplate(name, tpl string) error
}

// #3876
func TestHTMLEscape(t *testing.T) {
	assert := require.New(t)

	data := map[string]string{
		"html":  "<h1>Hi!</h1>",
		"other": "<h1>Hi!</h1>",
	}
	v := viper.New()
	fs := hugofs.NewMem(v)

	//afero.WriteFile(fs.Source, filepath.Join(workingDir, "README.txt"), []byte("Hugo Rocks!"), 0755)

	depsCfg := newDepsConfig(v)
	depsCfg.Fs = fs
	d, err := deps.New(depsCfg)
	assert.NoError(err)

	tpl := `{{ "<h1>Hi!</h1>" | safeHTML }}`

	provider := DefaultTemplateProvider
	provider.Update(d)

	h := d.Tmpl.(handler)

	assert.NoError(h.addTemplate("shortcodes/myShort.html", tpl))

	s, err := d.Tmpl.Lookup("shortcodes/myShort.html").ExecuteToString(data)
	assert.NoError(err)
	assert.Contains(s, "<h1>Hi!</h1>")

	s, err = d.Tmpl.Lookup("shortcodes/myShort").ExecuteToString(data)
	assert.NoError(err)
	assert.Contains(s, "<h1>Hi!</h1>")

}
