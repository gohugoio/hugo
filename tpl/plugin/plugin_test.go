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

package plugin_test

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/plugin"
	"github.com/spf13/viper"
)

const pluginName = "hello"

var cfg = viper.New()

func init() {
	cfg.Set("pluginDir", "../../resources/testdata/plugin")
}

type tstNoStringer struct{}

func TestOpen(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var ns = plugin.New(&deps.Deps{
		Cfg: cfg,
		Site: page.NewDummyHugoSite(cfg),
	})

	for _, test := range []interface{}{
		pluginName,
		filepath.FromSlash(pluginName),
	} {
		result, err := ns.Open(test)

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Not(qt.IsNil))

		vVar, err := result.Lookup(`Hello`)
		c.Assert(err, qt.IsNil)
		c.Assert(vVar, qt.Not(qt.IsNil))
	}

	// Errors

	for _, test := range []interface{}{
		filepath.FromSlash(`does-not-exist`),
		filepath.FromSlash(`not-plugin/path`),
		filepath.FromSlash(`sample`),
		nil,
	} {
		p, err := ns.Open(test)

		c.Assert(err, qt.Not(qt.IsNil))
		c.Assert(p, qt.IsNil)
	}
}


func TestExist(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var ns = plugin.New(&deps.Deps{
		Cfg: cfg,
		Site: page.NewDummyHugoSite(cfg),
	})

	for _, test := range []struct{
		PluginName interface{}
		Exists interface{}
	}{
		{pluginName, true},
		{filepath.FromSlash(pluginName), true},
		{filepath.FromSlash(`does-not-exist`), false},
		{filepath.FromSlash(`not-plugin/path`), false},
		{filepath.FromSlash(`sample`), false},
		{nil, false},
	} {
		ok, err := ns.Exist(test.PluginName)

		c.Assert(err, qt.IsNil)
		c.Assert(ok, qt.Equals, test.Exists)
	}
}
