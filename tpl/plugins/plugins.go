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

package plugins

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gohugoio/hugo/deps"
	"github.com/starlight-go/starlight"
)

// New returns a new instance of the plugins-namespaced template functions.  We
// allow the user to specify a plugins directory, and the theme may also have a
// plugin directory.  Plugins in the user's directory will override plugins in
// the theme directory if they have the same name.
func New(d *deps.Deps) *Namespace {
	dir := d.Cfg.GetString("plugin_dir")
	if dir == "" {
		dir = "plugins"
	}
	theme := d.Cfg.GetString("theme")
	themeDir := filepath.Join("./themes", theme, "plugins")
	return &Namespace{cache: starlight.New(dir, themeDir), deps: d}
}

// Namespace provides template functions for the "plugins" namespace.
// Each plugin type should implement its own function.
type Namespace struct {
	cache *starlight.Cache
	deps  *deps.Deps
}

// Starlight runs the given starlight script with the given values as globals.
// The globals are expected to be associated pairs of name (string) and value,
// which will be passed to the script. It's an error if you don't pass the
// values in pairs. The value of the template function is set by calling the
// "output" function in the script with the value you wish to output.
//
// e.g. {{plugins.starlight "foo.star" "page" .Page "site" .Site }}
//
// plugins/hello.star:
//
// output("hello" + site.Title)
func (ns *Namespace) Starlight(filename interface{}, vals ...interface{}) (interface{}, error) {
	s, ok := filename.(string)
	if !ok {
		return nil, fmt.Errorf("expected first argument to be filename (string), but was %v (%T)", filename, filename)
	}

	if len(vals)%2 != 0 {
		return nil, fmt.Errorf("expected values to be pairs of <name> <value>")
	}

	var ret interface{}
	output := func(v interface{}) {
		ret = v
	}
	globals := map[string]interface{}{
		"output": output,
	}

	for i := 0; i < len(vals); i += 2 {
		name, ok := vals[i].(string)
		if !ok {
			return nil, fmt.Errorf("expected argument %d to be string label, but was %v (%T)", i, vals[i], vals[i])
		}
		if name == "output" {
			return nil, errors.New("argument label `output` is a reserved word and cannot be used")
		}
		globals[name] = vals[i+1]
	}

	_, err := ns.cache.Run(s, globals)
	return ret, err
}
