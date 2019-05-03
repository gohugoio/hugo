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

package modules

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

type ThemeConfig struct {
	// The theme name as set in the configuration.
	// This maps either to a folder below /themes or
	// to a Go module Path.
	Name string

	// Set if the source lives in a Go module.
	Module *Module

	// The absolute path to this theme.
	Dir string

	// Optional configuration filename (e.g. "/themes/mytheme/config.json").
	ConfigFilename string

	// Optional config read from the ConfigFile above.
	Cfg config.Provider
}

// Collects and creates a module tree.
type collector struct {
	*Handler

	*collected
}

func (c *collector) initModules() error {
	c.collected = &collected{
		seen:     make(map[string]bool),
		vendored: make(map[string]string),
	}

	return c.loadModules()
}

const vendorModulesFilename = "modules.txt"

func (c *collector) collectModulesTXT(dir string) error {
	vendorDir := filepath.Join(dir, vendord)
	filename := filepath.Join(vendorDir, vendorModulesFilename)

	f, err := c.fs.Open(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// # github.com/alecthomas/chroma v0.6.3
		line := scanner.Text()
		line = strings.Trim(line, "# ")
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return errors.Errorf("invalid modules list: %q", filename)
		}
		path := parts[0]
		if _, found := c.vendored[path]; !found {
			c.vendored[path] = filepath.Join(vendorDir, path)
		}

	}
	return nil
}

func (c *collector) getVendoredDir(path string) string {
	return c.vendored[path]
}

func (c *collector) loadModules() error {
	modules, err := c.List()
	if err != nil {
		return err
	}
	c.modules = modules
	return nil
}

type collected struct {
	seen map[string]bool

	// Maps module path to a _vendor dir. These values are fetched from
	// _vendor/modules.txt, and the first (top-most) will win.
	vendored map[string]string

	// Set if a Go modules enabled project.
	modules Modules

	themes []ThemeConfig
}

// TODO(bep) mod rename these types.
// TODO(bep) mod plan for vendor:
// - iterate /vendor and create "virtual module" (VendorDir?)
// - no-vendor
func (c *collector) isSeen(theme string) bool {
	loki := strings.ToLower(theme)
	if c.seen[loki] {
		return true
	}
	c.seen[loki] = true
	return false
}

func (c *collector) addAndRecurse(dir string, themes ...string) error {
	for i := 0; i < len(themes); i++ {
		theme := themes[i]
		if !c.isSeen(theme) {
			tc, err := c.add(dir, theme)
			if err != nil {
				return err
			}
			if err := c.addThemeNamesFromTheme(tc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *collector) add(dir, name string) (ThemeConfig, error) {
	var tc ThemeConfig
	var mod *Module

	if err := c.collectModulesTXT(dir); err != nil {
		return ThemeConfig{}, err
	}

	// Try _vendor first.
	// TODO(bep) mod config flag
	moduleDir := c.getVendoredDir(name)

	if moduleDir == "" {
		mod = c.modules.GetByPath(name)
		if mod != nil {
			moduleDir = mod.Dir
		}

		if moduleDir == "" {
			if c.GoModulesFilename != "" && c.IsProbablyModule(name) {
				// Try to "go get" it and reload the module configuration.
				if err := c.Get(name); err != nil {
					return ThemeConfig{}, err
				}
				if err := c.loadModules(); err != nil {
					return ThemeConfig{}, err
				}

				mod = c.modules.GetByPath(name)
				if mod != nil {
					moduleDir = mod.Dir
				}
			}

			// Fall back to /themes/<mymodule>
			if moduleDir == "" {
				moduleDir = filepath.Join(c.themesDir, name)
				if found, _ := afero.Exists(c.fs, moduleDir); !found {
					return ThemeConfig{}, errors.Errorf("module %q not found", name)
				}
			}
		}
	}

	if found, _ := afero.Exists(c.fs, moduleDir); !found {
		return ThemeConfig{}, errors.Errorf("%q not found", moduleDir)
	}

	tc = ThemeConfig{
		Name:   name,
		Dir:    moduleDir,
		Module: mod,
	}

	if err := c.applyThemeConfig(&tc); err != nil {
		return tc, err
	}

	c.themes = append(c.themes, tc)
	return tc, nil

}

type ThemesConfig struct {
	Themes []ThemeConfig

	// Set if this is a Go modules enabled project.
	GoModulesFilename string
}

func (h *Handler) Collect() (ThemesConfig, error) {
	if len(h.imports) == 0 {
		return ThemesConfig{}, nil
	}

	c := &collector{
		Handler: h,
	}

	if err := c.collect(); err != nil {
		return ThemesConfig{}, err
	}

	return ThemesConfig{
		Themes:            c.themes,
		GoModulesFilename: c.GoModulesFilename,
	}, nil

}

func (c *collector) collect() error {
	if err := c.initModules(); err != nil {
		return err
	}

	for _, imp := range c.imports {
		if err := c.addAndRecurse(c.workingDir, imp); err != nil {
			return err
		}
	}

	return nil
}

func (c *collector) applyThemeConfig(tc *ThemeConfig) error {

	var (
		configFilename string
		cfg            config.Provider
		exists         bool
	)

	// Viper supports more, but this is the sub-set supported by Hugo.
	for _, configFormats := range config.ValidConfigFileExtensions {
		configFilename = filepath.Join(tc.Dir, "config."+configFormats)
		exists, _ = afero.Exists(c.fs, configFilename)
		if exists {
			break
		}
	}

	if !exists {
		// No theme config set.
		return nil
	}

	if configFilename != "" {
		var err error
		cfg, err = config.FromFile(c.fs, configFilename)
		if err != nil {
			return err
		}
	}

	tc.ConfigFilename = configFilename
	tc.Cfg = cfg

	return nil

}

func (c *collector) addThemeNamesFromTheme(theme ThemeConfig) error {
	if theme.Cfg != nil && theme.Cfg.IsSet("theme") {
		v := theme.Cfg.Get("theme")
		switch vv := v.(type) {
		case []string:
			return c.addAndRecurse(theme.Dir, vv...)
		case []interface{}:
			return c.addAndRecurse(theme.Dir, cast.ToStringSlice(vv)...)
		default:
			return c.addAndRecurse(theme.Dir, cast.ToString(vv))
		}
	}

	return nil
}
