// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package helpers

import (
	"sort"
	"strings"
	"sync"

	"github.com/spf13/cast"

	"github.com/spf13/viper"
)

// These are the settings that should only be looked up in the global Viper
// config and not per language.
// This list may not be complete, but contains only settings that we know
// will be looked up in both.
// This isn't perfect, but it is ultimately the user who shoots him/herself in
// the foot.
// See the pathSpec.
var globalOnlySettings = map[string]bool{
	strings.ToLower("defaultContentLanguageInSubdir"): true,
	strings.ToLower("defaultContentLanguage"):         true,
	strings.ToLower("multilingual"):                   true,
}

type Language struct {
	Lang         string
	LanguageName string
	Title        string
	Weight       int
	params       map[string]interface{}
	paramsInit   sync.Once
}

func (l *Language) String() string {
	return l.Lang
}

func NewLanguage(lang string) *Language {
	return &Language{Lang: lang, params: make(map[string]interface{})}
}

func NewDefaultLanguage() *Language {
	defaultLang := viper.GetString("defaultContentLanguage")

	if defaultLang == "" {
		defaultLang = "en"
	}

	return NewLanguage(defaultLang)
}

type Languages []*Language

func NewLanguages(l ...*Language) Languages {
	languages := make(Languages, len(l))
	for i := 0; i < len(l); i++ {
		languages[i] = l[i]
	}
	sort.Sort(languages)
	return languages
}

func (l Languages) Len() int           { return len(l) }
func (l Languages) Less(i, j int) bool { return l[i].Weight < l[j].Weight }
func (l Languages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func (l *Language) Params() map[string]interface{} {
	l.paramsInit.Do(func() {
		// Merge with global config.
		// TODO(bep) consider making this part of a constructor func.

		globalParams := viper.GetStringMap("params")
		for k, v := range globalParams {
			if _, ok := l.params[k]; !ok {
				l.params[k] = v
			}
		}
	})
	return l.params
}

func (l *Language) SetParam(k string, v interface{}) {
	l.params[strings.ToLower(k)] = v
}

// GetBool returns the value associated with the key as a boolean.
func (l *Language) GetBool(key string) bool { return cast.ToBool(l.Get(key)) }

// GetString returns the value associated with the key as a string.
func (l *Language) GetString(key string) string { return cast.ToString(l.Get(key)) }

// GetInt returns the value associated with the key as an int.
func (l *Language) GetInt(key string) int { return cast.ToInt(l.Get(key)) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func (l *Language) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(l.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (l *Language) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(l.Get(key))
}

// Get returns a value associated with the key relying on specified language.
// Get is case-insensitive for a key.
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func (l *Language) Get(key string) interface{} {
	if l == nil {
		panic("language not set")
	}
	key = strings.ToLower(key)
	if !globalOnlySettings[key] {
		if v, ok := l.params[key]; ok {
			return v
		}
	}
	return viper.Get(key)
}
