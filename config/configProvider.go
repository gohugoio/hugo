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
	"time"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/common/urls"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"
)

// AllProvider is a sub set of all config settings.
type AllProvider interface {
	Language() *langs.Language
	Languages() langs.Languages
	LanguagesDefaultFirst() langs.Languages
	LanguagePrefix() string
	BaseURL() urls.BaseURL
	BaseURLLiveReload() urls.BaseURL
	PathParser() *paths.PathParser
	Environment() string
	IsMultihost() bool
	IsMultilingual() bool
	NoBuildLock() bool
	BaseConfig() BaseConfig
	Dirs() CommonDirs
	Quiet() bool
	DirsBase() CommonDirs
	ContentTypes() ContentTypesProvider
	GetConfigSection(string) any
	GetConfig() any
	CanonifyURLs() bool
	DisablePathToLower() bool
	RemovePathAccents() bool
	IsUglyURLs(section string) bool
	DefaultContentLanguage() string
	DefaultContentLanguageInSubdir() bool
	IsLangDisabled(string) bool
	SummaryLength() int
	Pagination() Pagination
	BuildExpired() bool
	BuildFuture() bool
	BuildDrafts() bool
	Running() bool
	Watching() bool
	NewIdentityManager(name string) identity.Manager
	FastRenderMode() bool
	PrintUnusedTemplates() bool
	EnableMissingTranslationPlaceholders() bool
	TemplateMetrics() bool
	TemplateMetricsHints() bool
	PrintI18nWarnings() bool
	CreateTitle(s string) string
	IgnoreFile(s string) bool
	NewContentEditor() string
	Timeout() time.Duration
	StaticDirs() []string
	IgnoredLogs() map[string]bool
	WorkingDir() string
	EnableEmoji() bool
}

// We cannot import the media package as that would create a circular dependency.
// This interface defineds a sub set of what media.ContentTypes provides.
type ContentTypesProvider interface {
	IsContentSuffix(suffix string) bool
	IsContentFile(filename string) bool
	IsIndexContentFile(filename string) bool
	IsHTMLSuffix(suffix string) bool
}

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
	Keys() []string
	Merge(key string, value any)
	SetDefaults(params maps.Params)
	SetDefaultMergeStrategy()
	WalkParams(walkFn func(params ...maps.KeyParams) bool)
	IsSet(key string) bool
}

// GetStringSlicePreserveString returns a string slice from the given config and key.
// It differs from the GetStringSlice method in that if the config value is a string,
// we do not attempt to split it into fields.
func GetStringSlicePreserveString(cfg Provider, key string) []string {
	sd := cfg.Get(key)
	return types.ToStringSlicePreserveString(sd)
}
