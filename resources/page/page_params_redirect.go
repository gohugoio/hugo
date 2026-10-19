// Copyright 2026 The Hugo Authors. All rights reserved.
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

package page

import (
	"fmt"
	"iter"
	"strings"

	"github.com/bep/logg"
	"github.com/gohugoio/hugo/common/hugo"
)

// RootParamProvider is for internal use only.
type RootParamProvider interface {
	// IsRootParam reports whether key in Params was set from outside of a params section,
	// e.g. from a front matter field such as title, or by Hugo itself.
	IsRootParam(key string) bool
}

// reservedParamMethods maps reserved front matter fields that are also stored in Params
// to the Page method returning the same value.
var reservedParamMethods = map[string]string{
	"date":           "Date",
	"description":    "Description",
	"draft":          "Draft",
	"expirydate":     "ExpiryDate",
	"keywords":       "Keywords",
	"lastmod":        "Lastmod",
	"layout":         "Layout",
	"linktitle":      "LinkTitle",
	"publishdate":    "PublishDate",
	"slug":           "Slug",
	"summary":        "RawSummary",
	"title":          "Title",
	"translationkey": "TranslationKey",
	"type":           "Type",
	"weight":         "Weight",
}

// ReservedParamMethod returns the name of the Page method to use instead of
// looking up key in the Params of v, and whether key should be redirected.
// A key is only redirected if it was set from a reserved front matter field,
// so e.g. params.title in front matter is still looked up in Params.
func ReservedParamMethod(v any, key string) (string, bool) {
	key = strings.ToLower(key)
	method, found := reservedParamMethods[key]
	if !found {
		return "", false
	}
	p, ok := v.(RootParamProvider)
	if !ok || !p.IsRootParam(key) {
		return "", false
	}
	return method, true
}

// FindReservedParam returns the first page in pages where key would be
// redirected by ReservedParamMethod, and the name of the method to use.
func FindReservedParam(key string, pages iter.Seq[any]) (Page, string, bool) {
	if _, found := reservedParamMethods[strings.ToLower(key)]; !found {
		return nil, "", false
	}
	for v := range pages {
		if method, ok := ReservedParamMethod(v, key); ok {
			if p, ok := v.(Page); ok {
				return p, method, true
			}
		}
	}
	return nil, "", false
}

// DeprecateReservedParam logs that the reserved front matter field key
// was looked up in Params via usage, e.g. `ByParam "title"`.
func DeprecateReservedParam(log logg.Logger, usage, key, alternative string) {
	hugo.DeprecateWithLogger(
		fmt.Sprintf("Using %s to get the %s front matter field", usage, strings.ToLower(key)),
		fmt.Sprintf("Use %s instead.", alternative),
		"0.169.0",
		log,
	)
}

func (p Pages) deprecateReservedParam(usage, key string, alternative func(method string) string) {
	seq := func(yield func(any) bool) {
		for _, pp := range p {
			if !yield(pp) {
				return
			}
		}
	}
	if pp, method, ok := FindReservedParam(key, seq); ok {
		DeprecateReservedParam(pp.Language().Logger(), usage, key, alternative(method))
	}
}
