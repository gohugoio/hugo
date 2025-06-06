// Copyright 2020 The Hugo Authors. All rights reserved.
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

package constants

// Error/Warning IDs.
// Do not change these values.
const (
	// IDs for remote errors in tpl/data.
	ErrRemoteGetJSON = "error-remote-getjson"
	ErrRemoteGetCSV  = "error-remote-getcsv"

	WarnFrontMatterParamsOverrides = "warning-frontmatter-params-overrides"
	WarnRenderShortcodesInHTML     = "warning-rendershortcodes-in-html"
	WarnGoldmarkRawHTML            = "warning-goldmark-raw-html"
	WarnPartialSuperfluousPrefix   = "warning-partial-superfluous-prefix"
	WarnHomePageIsLeafBundle       = "warning-home-page-is-leaf-bundle"
)

// Field/method names with special meaning.
const (
	FieldRelPermalink = "RelPermalink"
	FieldPermalink    = "Permalink"
)

// IsFieldRelOrPermalink returns whether the given name is a RelPermalink or Permalink.
func IsFieldRelOrPermalink(name string) bool {
	return name == FieldRelPermalink || name == FieldPermalink
}

// Resource transformations.
const (
	ResourceTransformationFingerprint = "fingerprint"
)

// IsResourceTransformationPermalinkHash returns whether the given name is a resource transformation that changes the permalink based on the content.
func IsResourceTransformationPermalinkHash(name string) bool {
	return name == ResourceTransformationFingerprint
}
