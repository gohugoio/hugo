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

import "strings"

// RootParamProvider is for internal use only.
type RootParamProvider interface {
	// IsRootParam reports whether key in Params was set from outside of a params section,
	// e.g. from a front matter field such as title, or by Hugo itself.
	IsRootParam(key string) bool
}

// reservedParamMethods maps reserved front matter fields that are also stored in Params
// to the Page method returning the same value.
var reservedParamMethods = map[string]string{
	"description":    "Description",
	"draft":          "Draft",
	"keywords":       "Keywords",
	"layout":         "Layout",
	"linktitle":      "LinkTitle",
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
