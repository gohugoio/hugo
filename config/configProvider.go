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

package config

import (
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
)

// Provider provides the configuration settings for Hugo.
type Provider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetParams(key string) maps.Params
	GetStringMap(key string) map[string]any
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	Get(key string) any
	Set(key string, value any)
	Merge(key string, value any)
	SetDefaults(params maps.Params)
	SetDefaultMergeStrategy()
	WalkParams(walkFn func(params ...KeyParams) bool)
	IsSet(key string) bool
}

// GetStringSlicePreserveString returns a string slice from the given config and key.
// It differs from the GetStringSlice method in that if the config value is a string,
// we do not attempt to split it into fields.
func GetStringSlicePreserveString(cfg Provider, key string) []string {
	sd := cfg.Get(key)
	return types.ToStringSlicePreserveString(sd)
}

// SetBaseTestDefaults provides some common config defaults used in tests.
func SetBaseTestDefaults(cfg Provider) Provider {
	setIfNotSet(cfg, "baseURL", "https://example.org")
	setIfNotSet(cfg, "resourceDir", "resources")
	setIfNotSet(cfg, "contentDir", "content")
	setIfNotSet(cfg, "dataDir", "data")
	setIfNotSet(cfg, "i18nDir", "i18n")
	setIfNotSet(cfg, "layoutDir", "layouts")
	setIfNotSet(cfg, "assetDir", "assets")
	setIfNotSet(cfg, "archetypeDir", "archetypes")
	setIfNotSet(cfg, "publishDir", "public")
	setIfNotSet(cfg, "workingDir", "")
	setIfNotSet(cfg, "defaultContentLanguage", "en")
	return cfg
}

func setIfNotSet(cfg Provider, key string, value any) {
	if !cfg.IsSet(key) {
		cfg.Set(key, value)
	}
}
