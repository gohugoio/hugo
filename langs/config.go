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

package langs

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"sort"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"
)

// LanguageConfig holds the configuration for a single language.
// This is what is read from the config file.
type LanguageConfig struct {
	// The language name, e.g. "English".
	LanguageName string

	// The language code, e.g. "en-US".
	LanguageCode string

	// The language title. When set, this will
	// override site.Title for this language.
	Title string

	// The language direction, e.g. "ltr" or "rtl".
	LanguageDirection string

	// The language weight. When set to a non-zero value, this will
	// be the main sort criteria for the language.
	Weight int

	// Set to true to disable this language.
	Disabled bool
}

type LanguageInternal struct {
	// Name is the name of the role, extracted from the key in the config.
	Name string

	// Whether this role is the default role.
	// This will be rendered in the root.
	// There is only be one default role.
	Default bool

	LanguageConfig
}

type LanguagesInternal struct {
	LanguageConfigs map[string]LanguageConfig
	Sorted          []LanguageInternal
}

func (ls LanguagesInternal) IndexDefault() int {
	for i, role := range ls.Sorted {
		if role.Default {
			return i
		}
	}
	panic("no default role found")
}

func (ls LanguagesInternal) ResolveName(i int) string {
	if i < 0 || i >= len(ls.Sorted) {
		panic(fmt.Sprintf("index %d out of range for languages", i))
	}
	return ls.Sorted[i].Name
}

func (ls LanguagesInternal) ResolveIndex(name string) int {
	for i, role := range ls.Sorted {
		if role.Name == name {
			return i
		}
	}
	panic(fmt.Sprintf("no language found for name %q", name))
}

func (ls LanguagesInternal) Len() int {
	return len(ls.Sorted)
}

// IndexMatch returns an iterator for the roles that match the filter.
func (ls LanguagesInternal) IndexMatch(match predicate.P[string]) (iter.Seq[int], error) {
	return func(yield func(i int) bool) {
		for i, l := range ls.Sorted {
			if match(l.Name) {
				if !yield(i) {
					return
				}
			}
		}
	}, nil
}

// ForEachIndex returns an iterator for the indices of the languages.
func (ls LanguagesInternal) ForEachIndex() iter.Seq[int] {
	return func(yield func(i int) bool) {
		for i := range ls.Sorted {
			if !yield(i) {
				return
			}
		}
	}
}

func (ls *LanguagesInternal) init(defaultContentLanguage string, disabledLanguages []string) (string, error) {
	const en = "en"

	if len(ls.LanguageConfigs) == 0 {
		// Add a default language.
		if defaultContentLanguage == "" {
			defaultContentLanguage = en
		}
		ls.LanguageConfigs[defaultContentLanguage] = LanguageConfig{}
	}

	var (
		defaultSeen bool
		enIdx       int = -1
	)
	for k, v := range ls.LanguageConfigs {
		if !v.Disabled && slices.Contains(disabledLanguages, k) {
			// This language is disabled.
			v.Disabled = true
			ls.LanguageConfigs[k] = v
		}

		if k == "" {
			return "", errors.New("language name cannot be empty")
		}

		if err := paths.ValidateIdentifier(k); err != nil {
			return "", fmt.Errorf("language name %q is invalid: %s", k, err)
		}

		var isDefault bool
		if k == defaultContentLanguage {
			isDefault = true
			defaultSeen = true
		}

		if isDefault && v.Disabled {
			return "", fmt.Errorf("default language %q is disabled", k)
		}

		if !v.Disabled {
			ls.Sorted = append(ls.Sorted, LanguageInternal{Name: k, Default: isDefault, LanguageConfig: v})
		}
	}

	// Sort by weight if set, then by name.
	sort.SliceStable(ls.Sorted, func(i, j int) bool {
		ri, rj := ls.Sorted[i], ls.Sorted[j]
		if ri.Weight == rj.Weight {
			return ri.Name < rj.Name
		}
		if rj.Weight == 0 {
			return true
		}
		if ri.Weight == 0 {
			return false
		}
		return ri.Weight < rj.Weight
	})

	for i, l := range ls.Sorted {
		if l.Name == en {
			enIdx = i
			break
		}
	}

	if !defaultSeen {
		if defaultContentLanguage != "" {
			// Set by the user, but not found in the config.
			return "", fmt.Errorf("defaultContentLanguage %q not found in languages configuration", defaultContentLanguage)
		}
		// Not set by the user, so we use the first language in the config.
		defaultIdx := 0
		if enIdx != -1 {
			defaultIdx = enIdx
		}
		d := ls.Sorted[defaultIdx]
		d.Default = true
		ls.LanguageConfigs[d.Name] = d.LanguageConfig
		ls.Sorted[defaultIdx] = d
		defaultContentLanguage = d.Name

	}

	return defaultContentLanguage, nil
}

func DecodeConfig(defaultContentLanguage string, disabledLanguages []string, m map[string]any) (*config.ConfigNamespace[map[string]LanguageConfig, LanguagesInternal], string, error) {
	v, err := config.DecodeNamespace[map[string]LanguageConfig](m, func(in any) (LanguagesInternal, any, error) {
		var languages LanguagesInternal
		var conf map[string]LanguageConfig
		if err := mapstructure.Decode(m, &conf); err != nil {
			return languages, nil, err
		}
		languages.LanguageConfigs = conf
		var err error
		if defaultContentLanguage, err = languages.init(defaultContentLanguage, disabledLanguages); err != nil {
			return languages, nil, err
		}
		return languages, languages.LanguageConfigs, nil
	})

	return v, defaultContentLanguage, err
}
