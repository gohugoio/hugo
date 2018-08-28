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

package paths

import (
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type ThemeConfig struct {
	// The theme name as provided by the folder name below /themes.
	Name string

	// Optional configuration filename (e.g. "/themes/mytheme/config.json").
	ConfigFilename string

	// Optional config read from the ConfigFile above.
	Cfg config.Provider
}

// Create file system, an ordered theme list from left to right, no duplicates.
type themesCollector struct {
	themesDir string
	fs        afero.Fs
	seen      map[string]bool
	themes    []ThemeConfig
}

func (c *themesCollector) isSeen(theme string) bool {
	loki := strings.ToLower(theme)
	if c.seen[loki] {
		return true
	}
	c.seen[loki] = true
	return false
}

func (c *themesCollector) addAndRecurse(themes ...string) error {
	for i := 0; i < len(themes); i++ {
		theme := themes[i]
		configFilename := c.getConfigFileIfProvided(theme)
		if !c.isSeen(theme) {
			tc, err := c.add(theme, configFilename)
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

func (c *themesCollector) add(name, configFilename string) (ThemeConfig, error) {
	var cfg config.Provider
	var tc ThemeConfig

	if configFilename != "" {
		v := viper.New()
		v.SetFs(c.fs)
		v.AutomaticEnv()
		v.SetEnvPrefix("hugo")
		v.SetConfigFile(configFilename)

		err := v.ReadInConfig()
		if err != nil {
			return tc, err
		}
		cfg = v

	}

	tc = ThemeConfig{Name: name, ConfigFilename: configFilename, Cfg: cfg}
	c.themes = append(c.themes, tc)
	return tc, nil

}

func collectThemeNames(p *Paths) ([]ThemeConfig, error) {
	return CollectThemes(p.Fs.Source, p.AbsPathify(p.ThemesDir), p.Themes())

}

func CollectThemes(fs afero.Fs, themesDir string, themes []string) ([]ThemeConfig, error) {
	if len(themes) == 0 {
		return nil, nil
	}

	c := &themesCollector{
		fs:        fs,
		themesDir: themesDir,
		seen:      make(map[string]bool)}

	for i := 0; i < len(themes); i++ {
		theme := themes[i]
		if err := c.addAndRecurse(theme); err != nil {
			return nil, err
		}
	}

	return c.themes, nil

}

func (c *themesCollector) getConfigFileIfProvided(theme string) string {
	configDir := filepath.Join(c.themesDir, theme)

	var (
		configFilename string
		exists         bool
	)

	// Viper supports more, but this is the sub-set supported by Hugo.
	for _, configFormats := range []string{"toml", "yaml", "yml", "json"} {
		configFilename = filepath.Join(configDir, "config."+configFormats)
		exists, _ = afero.Exists(c.fs, configFilename)
		if exists {
			break
		}
	}

	if !exists {
		// No theme config set.
		return ""
	}

	return configFilename

}

func (c *themesCollector) addThemeNamesFromTheme(theme ThemeConfig) error {
	if theme.Cfg != nil && theme.Cfg.IsSet("theme") {
		v := theme.Cfg.Get("theme")
		switch vv := v.(type) {
		case []string:
			return c.addAndRecurse(vv...)
		case []interface{}:
			return c.addAndRecurse(cast.ToStringSlice(vv)...)
		default:
			return c.addAndRecurse(cast.ToString(vv))
		}
	}

	return nil
}
