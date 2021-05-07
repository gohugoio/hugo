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

func TestGet(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var ns = plugin.New(&deps.Deps{
		Cfg: cfg,
		Site: page.NewDummyHugoSite(cfg),
	})

	HelloFmt := map[string]string{
		"english": "Hello %s",
		"french":  "Salutation %s",
		"spanish": "Hola %s",
	}

	for _, test := range []struct {
		args   [2]interface{}
		expect interface{}
	}{
		{[2]interface{}{pluginName, `HelloFmt`}, &HelloFmt},
		{[2]interface{}{filepath.FromSlash(pluginName), `HelloFmt`}, &HelloFmt},
	} {
		v, err := ns.Get(test.args[0], test.args[1])

		c.Assert(err, qt.IsNil)
		c.Assert(v, qt.DeepEquals, test.expect)
	}

	// Errors

	for _, test := range [][2]interface{}{
		{filepath.FromSlash(`does-not-exist`), `field`},
		{filepath.FromSlash(`not-plugin/path`), `field`},
		{filepath.FromSlash(`sample`), `field`},
		{nil, `field`},
		{filepath.FromSlash(pluginName), `UnknownField`},
		{filepath.FromSlash(pluginName), `helloFmt`},
		{filepath.FromSlash(pluginName), `Language`},
	} {
		p, err := ns.Get(test[0], test[1])

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

func TestHas(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var ns = plugin.New(&deps.Deps{
		Cfg: cfg,
		Site: page.NewDummyHugoSite(cfg),
	})

	for _, test := range []struct{
		PluginName interface{}
		VariableName interface{}
		Expect interface{}
	}{
		{pluginName, `Hello`, true},
		{filepath.FromSlash(pluginName), `Hello`, true},
		{filepath.FromSlash(pluginName), `HelloFmt`, true},
		{filepath.FromSlash(pluginName), `Language`, false},
		{filepath.FromSlash(pluginName), `doesNotExist`, false},
		{filepath.FromSlash(pluginName), `not - a + variable`, false},
		{filepath.FromSlash(pluginName), nil, false},
	} {
		ok, err := ns.Has(test.PluginName, test.VariableName)

		c.Assert(err, qt.IsNil)
		c.Assert(ok, qt.Equals, test.Expect)
	}
}

func TestCall(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	var ns = plugin.New(&deps.Deps{
		Cfg: cfg,
		Site: page.NewDummyHugoSite(cfg),
	})

	for _, test := range []struct{
		PluginName interface{}
		FunctionName interface{}
		Arguments []interface{}
		Expect interface{}
	}{
		{pluginName, `Hello`, []interface{}{"holyhope"}, "Hello holyhope"},
		{filepath.FromSlash(pluginName), `Hello`, []interface{}{"holyhope"}, "Hello holyhope"},
	} {
		fn, err := ns.Get(test.PluginName, test.FunctionName)
		c.Assert(err, qt.IsNil)

		result, err := ns.Call(fn, test.Arguments...)

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.Expect)
	}

	/* Errors */

	for _, test := range []struct{
		PluginName interface{}
		FunctionName interface{}
		Arguments []interface{}
	}{
		{pluginName, `Hello`, []interface{}{(*int)(nil)}},
		{filepath.FromSlash(pluginName), `Hello`, []interface{}{3}},
	} {
		fn, err := ns.Get(test.PluginName, test.FunctionName)
		c.Assert(err, qt.IsNil)

		_, err = ns.Call(fn, test.Arguments...)
		c.Assert(err, qt.Not(qt.IsNil))
	}
}
