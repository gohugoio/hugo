// Copyright 2024 The Hugo Authors. All rights reserved.
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

package media

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/config"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

// DefaultTypes is the default media types supported by Hugo.
var DefaultTypes Types

func init() {
	// Apply delimiter to all.
	for _, m := range defaultMediaTypesConfig {
		m.(map[string]any)["delimiter"] = "."
	}

	ns, err := DecodeTypes(nil)
	if err != nil {
		panic(err)
	}
	DefaultTypes = ns.Config

	// Initialize the Builtin types with values from DefaultTypes.
	v := reflect.ValueOf(&Builtin).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fieldName := v.Type().Field(i).Name
		builtinType := f.Interface().(Type)
		if builtinType.Type == "" {
			panic(fmt.Errorf("builtin type %q is empty", fieldName))
		}
		defaultType, found := DefaultTypes.GetByType(builtinType.Type)
		if !found {
			panic(fmt.Errorf("missing default type for field builtin type: %q", fieldName))
		}
		f.Set(reflect.ValueOf(defaultType))
	}
}

func init() {
	DefaultContentTypes = ContentTypes{
		HTML:             Builtin.HTMLType,
		Markdown:         Builtin.MarkdownType,
		AsciiDoc:         Builtin.AsciiDocType,
		Pandoc:           Builtin.PandocType,
		ReStructuredText: Builtin.ReStructuredTextType,
		EmacsOrgMode:     Builtin.EmacsOrgModeType,
	}

	DefaultContentTypes.init()
}

var DefaultContentTypes ContentTypes

// ContentTypes holds the media types that are considered content in Hugo.
type ContentTypes struct {
	HTML             Type
	Markdown         Type
	AsciiDoc         Type
	Pandoc           Type
	ReStructuredText Type
	EmacsOrgMode     Type

	// Created in init().
	types        Types
	extensionSet map[string]bool
}

func (t *ContentTypes) init() {
	t.types = Types{t.HTML, t.Markdown, t.AsciiDoc, t.Pandoc, t.ReStructuredText, t.EmacsOrgMode}
	t.extensionSet = make(map[string]bool)
	for _, mt := range t.types {
		for _, suffix := range mt.Suffixes() {
			t.extensionSet[suffix] = true
		}
	}
}

func (t ContentTypes) IsContentSuffix(suffix string) bool {
	return t.extensionSet[suffix]
}

// IsContentFile returns whether the given filename is a content file.
func (t ContentTypes) IsContentFile(filename string) bool {
	return t.IsContentSuffix(strings.TrimPrefix(filepath.Ext(filename), "."))
}

// IsIndexContentFile returns whether the given filename is an index content file.
func (t ContentTypes) IsIndexContentFile(filename string) bool {
	if !t.IsContentFile(filename) {
		return false
	}

	base := filepath.Base(filename)

	return strings.HasPrefix(base, "index.") || strings.HasPrefix(base, "_index.")
}

// IsHTMLSuffix returns whether the given suffix is a HTML media type.
func (t ContentTypes) IsHTMLSuffix(suffix string) bool {
	for _, s := range t.HTML.Suffixes() {
		if s == suffix {
			return true
		}
	}
	return false
}

// Types is a slice of media types.
func (t ContentTypes) Types() Types {
	return t.types
}

// FromTypes creates a new ContentTypes updated with the values from the given Types.
func (t ContentTypes) FromTypes(types Types) ContentTypes {
	if tt, ok := types.GetByType(t.HTML.Type); ok {
		t.HTML = tt
	}
	if tt, ok := types.GetByType(t.Markdown.Type); ok {
		t.Markdown = tt
	}
	if tt, ok := types.GetByType(t.AsciiDoc.Type); ok {
		t.AsciiDoc = tt
	}
	if tt, ok := types.GetByType(t.Pandoc.Type); ok {
		t.Pandoc = tt
	}
	if tt, ok := types.GetByType(t.ReStructuredText.Type); ok {
		t.ReStructuredText = tt
	}
	if tt, ok := types.GetByType(t.EmacsOrgMode.Type); ok {
		t.EmacsOrgMode = tt
	}

	t.init()

	return t
}

// Hold the configuration for a given media type.
type MediaTypeConfig struct {
	// The file suffixes used for this media type.
	Suffixes []string
	// Delimiter used before suffix.
	Delimiter string
}

// DecodeTypes decodes the given map of media types.
func DecodeTypes(in map[string]any) (*config.ConfigNamespace[map[string]MediaTypeConfig, Types], error) {
	buildConfig := func(v any) (Types, any, error) {
		m, err := maps.ToStringMapE(v)
		if err != nil {
			return nil, nil, err
		}
		if m == nil {
			m = map[string]any{}
		}
		m = maps.CleanConfigStringMap(m)
		// Merge with defaults.
		maps.MergeShallow(m, defaultMediaTypesConfig)

		var types Types

		for k, v := range m {
			mediaType, err := FromString(k)
			if err != nil {
				return nil, nil, err
			}
			if err := mapstructure.WeakDecode(v, &mediaType); err != nil {
				return nil, nil, err
			}
			mm := maps.ToStringMap(v)
			suffixes, _, found := maps.LookupEqualFold(mm, "suffixes")
			if found {
				mediaType.SuffixesCSV = strings.TrimSpace(strings.ToLower(strings.Join(cast.ToStringSlice(suffixes), ",")))
			}
			if mediaType.SuffixesCSV != "" && mediaType.Delimiter == "" {
				mediaType.Delimiter = DefaultDelimiter
			}
			InitMediaType(&mediaType)
			types = append(types, mediaType)
		}

		sort.Sort(types)

		return types, m, nil
	}

	ns, err := config.DecodeNamespace[map[string]MediaTypeConfig](in, buildConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode media types: %w", err)
	}
	return ns, nil
}

// TODO(bep) get rid of this.
var DefaultPathParser = &paths.PathParser{
	IsContentExt: func(ext string) bool {
		return DefaultContentTypes.IsContentSuffix(ext)
	},
}
