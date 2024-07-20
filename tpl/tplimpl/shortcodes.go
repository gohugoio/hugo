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

package tplimpl

import (
	"strings"

	"github.com/gohugoio/hugo/tpl"
)

// Currently lang, outFormat, suffix
const numTemplateVariants = 3

type shortcodeVariant struct {
	// The possible variants: lang, outFormat, suffix
	// gtag
	// gtag.html
	// gtag.no.html
	// gtag.no.amp.html
	// A slice of length numTemplateVariants.
	variants []string

	ts *templateState
}

type shortcodeTemplates struct {
	variants []shortcodeVariant
}

func (s *shortcodeTemplates) indexOf(variants []string) int {
L:
	for i, v1 := range s.variants {
		for i, v2 := range v1.variants {
			if v2 != variants[i] {
				continue L
			}
		}
		return i
	}
	return -1
}

func (s *shortcodeTemplates) fromVariants(variants tpl.TemplateVariants) (shortcodeVariant, bool) {
	return s.fromVariantsSlice([]string{
		variants.Language,
		strings.ToLower(variants.OutputFormat.Name),
		variants.OutputFormat.MediaType.FirstSuffix.Suffix,
	})
}

func (s *shortcodeTemplates) fromVariantsSlice(variants []string) (shortcodeVariant, bool) {
	var (
		bestMatch       shortcodeVariant
		bestMatchWeight int
	)

	for _, variant := range s.variants {
		w := s.compareVariants(variants, variant.variants)
		if bestMatchWeight == 0 || w > bestMatchWeight {
			bestMatch = variant
			bestMatchWeight = w
		}
	}

	return bestMatch, true
}

// calculate a weight for two string slices of same length.
// higher value means "better match".
func (s *shortcodeTemplates) compareVariants(a, b []string) int {
	weight := 0
	k := len(a)
	for i, av := range a {
		bv := b[i]
		if av == bv {
			// Add more weight to the left side (language...).
			weight = weight + k - i
		} else {
			weight--
		}
	}
	return weight
}

func templateVariants(name string) []string {
	_, variants := templateNameAndVariants(name)
	return variants
}

func templateNameAndVariants(name string) (string, []string) {
	variants := make([]string, numTemplateVariants)

	parts := strings.Split(name, ".")

	if len(parts) <= 1 {
		// No variants.
		return name, variants
	}

	name = parts[0]
	parts = parts[1:]
	lp := len(parts)
	start := len(variants) - lp

	for i, j := start, 0; i < len(variants); i, j = i+1, j+1 {
		variants[i] = parts[j]
	}

	if lp > 1 && lp < len(variants) {
		for i := lp - 1; i > 0; i-- {
			variants[i-1] = variants[i]
		}
	}

	if lp == 1 {
		// Suffix only. Duplicate it into the output format field to
		// make HTML win over AMP.
		variants[len(variants)-2] = variants[len(variants)-1]
	}

	return name, variants
}

func resolveTemplateType(name string) templateType {
	if isShortcode(name) {
		return templateShortcode
	}

	if strings.Contains(name, "partials/") {
		return templatePartial
	}

	return templateUndefined
}

func isShortcode(name string) bool {
	return strings.Contains(name, shortcodesPathPrefix)
}

func isInternal(name string) bool {
	return strings.HasPrefix(name, internalPathPrefix)
}
