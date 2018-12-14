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
	return newNamespace(dir, themeDir)
}

func newNamespace(userDir, themeDir string) *Namespace {
	return &Namespace{cache: starlight.New(userDir, themeDir)}
}

// Namespace provides template functions for the "plugins" namespace.
// Each plugin type should implement its own function.
type Namespace struct {
	cache *starlight.Cache
}
