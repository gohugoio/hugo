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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/common/paths"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/spf13/afero"
)

var (
	// See issue #8979 for context.
	// Hugo has always used config.toml etc. as the default config file name.
	// But hugo.toml is a more descriptive name, but we need to check for both.
	DefaultConfigNames = []string{"hugo", "config"}

	DefaultConfigNamesSet = make(map[string]bool)

	ValidConfigFileExtensions                    = []string{"toml", "yaml", "yml", "json"}
	validConfigFileExtensionsMap map[string]bool = make(map[string]bool)
)

func init() {
	for _, name := range DefaultConfigNames {
		DefaultConfigNamesSet[name] = true
	}

	for _, ext := range ValidConfigFileExtensions {
		validConfigFileExtensionsMap[ext] = true
	}
}

// IsValidConfigFilename returns whether filename is one of the supported
// config formats in Hugo.
func IsValidConfigFilename(filename string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	return validConfigFileExtensionsMap[ext]
}

func FromTOMLConfigString(config string) Provider {
	cfg, err := FromConfigString(config, "toml")
	if err != nil {
		panic(err)
	}
	return cfg
}

// FromConfigString creates a config from the given YAML, JSON or TOML config. This is useful in tests.
func FromConfigString(config, configType string) (Provider, error) {
	m, err := readConfig(metadecoders.FormatFromString(configType), []byte(config))
	if err != nil {
		return nil, err
	}
	return NewFrom(m), nil
}

// FromFile loads the configuration from the given filename.
func FromFile(fs afero.Fs, filename string) (Provider, error) {
	m, err := loadConfigFromFile(fs, filename)
	if err != nil {
		fe := herrors.UnwrapFileError(err)
		if fe != nil {
			pos := fe.Position()
			pos.Filename = filename
			fe.UpdatePosition(pos)
			return nil, err
		}
		return nil, herrors.NewFileErrorFromFile(err, filename, fs, nil)
	}
	return NewFrom(m), nil
}

// FromFileToMap is the same as FromFile, but it returns the config values
// as a simple map.
func FromFileToMap(fs afero.Fs, filename string) (map[string]any, error) {
	return loadConfigFromFile(fs, filename)
}

func readConfig(format metadecoders.Format, data []byte) (map[string]any, error) {
	m, err := metadecoders.Default.UnmarshalToMap(data, format)
	if err != nil {
		return nil, err
	}

	RenameKeys(m)

	return m, nil
}

func loadConfigFromFile(fs afero.Fs, filename string) (map[string]any, error) {
	m, err := metadecoders.Default.UnmarshalFileToMap(fs, filename)
	if err != nil {
		return nil, err
	}
	RenameKeys(m)
	return m, nil
}

func LoadConfigFromDir(sourceFs afero.Fs, configDir, environment string) (Provider, []string, error) {
	defaultConfigDir := filepath.Join(configDir, "_default")
	environmentConfigDir := filepath.Join(configDir, environment)
	cfg := New()

	var configDirs []string
	// Merge from least to most specific.
	for _, dir := range []string{defaultConfigDir, environmentConfigDir} {
		if _, err := sourceFs.Stat(dir); err == nil {
			configDirs = append(configDirs, dir)
		}
	}

	if len(configDirs) == 0 {
		return nil, nil, nil
	}

	// Keep track of these so we can watch them for changes.
	var dirnames []string

	for _, configDir := range configDirs {
		err := afero.Walk(sourceFs, configDir, func(path string, fi os.FileInfo, err error) error {
			if fi == nil || err != nil {
				return nil
			}

			if fi.IsDir() {
				dirnames = append(dirnames, path)
				return nil
			}

			if !IsValidConfigFilename(path) {
				return nil
			}

			name := paths.Filename(filepath.Base(path))

			item, err := metadecoders.Default.UnmarshalFileToMap(sourceFs, path)
			if err != nil {
				// This will be used in error reporting, use the most specific value.
				dirnames = []string{path}
				return fmt.Errorf("failed to unmarshal config for path %q: %w", path, err)
			}

			var keyPath []string
			if !DefaultConfigNamesSet[name] {
				// Can be params.jp, menus.en etc.
				name, lang := paths.FileAndExtNoDelimiter(name)

				keyPath = []string{name}

				if lang != "" {
					keyPath = []string{"languages", lang}
					switch name {
					case "menu", "menus":
						keyPath = append(keyPath, "menus")
					case "params":
						keyPath = append(keyPath, "params")
					}
				}
			}

			root := item
			if len(keyPath) > 0 {
				root = make(map[string]any)
				m := root
				for i, key := range keyPath {
					if i >= len(keyPath)-1 {
						m[key] = item
					} else {
						nm := make(map[string]any)
						m[key] = nm
						m = nm
					}
				}
			}

			// Migrate menu => menus etc.
			RenameKeys(root)

			// Set will overwrite keys with the same name, recursively.
			cfg.Set("", root)

			return nil
		})
		if err != nil {
			return nil, dirnames, err
		}

	}

	return cfg, dirnames, nil
}

var keyAliases maps.KeyRenamer

func init() {
	var err error
	keyAliases, err = maps.NewKeyRenamer(
		// Before 0.53 we used singular for "menu".
		"{menu,languages/*/menu}", "menus",
	)

	if err != nil {
		panic(err)
	}
}

// RenameKeys renames config keys in m recursively according to a global Hugo
// alias definition.
func RenameKeys(m map[string]any) {
	keyAliases.Rename(m)
}
