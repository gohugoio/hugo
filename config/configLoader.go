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
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var (
	ValidConfigFileExtensions                    = []string{"toml", "yaml", "yml", "json"}
	validConfigFileExtensionsMap map[string]bool = make(map[string]bool)
)

func init() {
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

// FromConfigString creates a config from the given YAML, JSON or TOML config. This is useful in tests.
func FromConfigString(config, configType string) (Provider, error) {
	v := newViper()
	m, err := readConfig(metadecoders.FormatFromString(configType), []byte(config))
	if err != nil {
		return nil, err
	}

	v.MergeConfigMap(m)

	return v, nil
}

// FromFile loads the configuration from the given filename.
func FromFile(fs afero.Fs, filename string) (Provider, error) {
	m, err := loadConfigFromFile(fs, filename)
	if err != nil {
		return nil, err
	}

	v := newViper()

	err = v.MergeConfigMap(m)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// FromFileToMap is the same as FromFile, but it returns the config values
// as a simple map.
func FromFileToMap(fs afero.Fs, filename string) (map[string]interface{}, error) {
	return loadConfigFromFile(fs, filename)
}

func readConfig(format metadecoders.Format, data []byte) (map[string]interface{}, error) {
	m, err := metadecoders.Default.UnmarshalToMap(data, format)
	if err != nil {
		return nil, err
	}

	RenameKeys(m)

	return m, nil

}

func loadConfigFromFile(fs afero.Fs, filename string) (map[string]interface{}, error) {
	m, err := metadecoders.Default.UnmarshalFileToMap(fs, filename)
	if err != nil {
		return nil, err
	}
	RenameKeys(m)
	return m, nil
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
func RenameKeys(m map[string]interface{}) {
	keyAliases.Rename(m)
}

func newViper() *viper.Viper {
	v := viper.New()

	return v
}
