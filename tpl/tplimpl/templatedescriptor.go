// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/resources/kinds"
)

const baseNameBaseof = "baseof"

// This is used both as a key and in lookups.
type TemplateDescriptor struct {
	// Group 1.
	Kind   string // page, home, section, taxonomy, term (and only those)
	Layout string // list, single, baseof, mycustomlayout.

	// Group 2.
	OutputFormat string // rss, csv ...
	MediaType    string // text/html, text/plain, ...
	Lang         string // en, nn, fr, ...

	Variant1 string // contextual variant, e.g. "link" in render hooks."
	Variant2 string // contextual variant, e.g. "id" in render.

	// Misc.
	LayoutMustMatch bool // If set, we only look for the exact layout.
	IsPlainText     bool // Whether this is a plain text template.
}

func (d *TemplateDescriptor) normalizeFromFile() {
	// fmt.Println("normalizeFromFile", "kind:", d.Kind, "layout:", d.Layout, "of:", d.OutputFormat)

	if d.Layout == d.OutputFormat {
		d.Layout = ""
	}

	if d.Kind == kinds.KindTemporary {
		d.Kind = ""
	}

	if d.Layout == d.Kind {
		d.Layout = ""
	}
}

type descriptorHandler struct {
	opts StoreOptions
}

// Note that this in this setup is usually a descriptor constructed from a page,
// so we want to find the best match for that page.
func (s descriptorHandler) compareDescriptors(category Category, isEmbedded bool, this, other TemplateDescriptor) weight {
	if this.LayoutMustMatch && this.Layout != other.Layout {
		return weightNoMatch
	}

	w := this.doCompare(category, isEmbedded, other)

	if w.w1 <= 0 {
		if category == CategoryMarkup && (this.Variant1 == other.Variant1) && (this.Variant2 == other.Variant2 || this.Variant2 != "" && other.Variant2 == "") {
			// See issue 13242.
			if this.OutputFormat != other.OutputFormat && this.OutputFormat == s.opts.DefaultOutputFormat {
				return w
			}

			w.w1 = 1
			return w
		}
	}

	return w
}

//lint:ignore ST1006 this vs other makes it easier to reason about.
func (this TemplateDescriptor) doCompare(category Category, isEmbedded bool, other TemplateDescriptor) weight {
	w := weightNoMatch

	// HTML in plain text is OK, but not the other way around.
	if other.IsPlainText && !this.IsPlainText {
		return w
	}
	if other.Kind != "" && other.Kind != this.Kind {
		return w
	}

	if other.Layout != "" && other.Layout != layoutAll && other.Layout != this.Layout {
		if isLayoutCustom(this.Layout) {
			if this.Kind == "" {
				this.Layout = ""
			} else if this.Kind == kinds.KindPage {
				this.Layout = layoutSingle
			} else {
				this.Layout = layoutList
			}
		}

		// Test again.
		if other.Layout != this.Layout {
			return w
		}
	}

	if other.Lang != "" && other.Lang != this.Lang {
		return w
	}

	if other.OutputFormat != "" && other.OutputFormat != this.OutputFormat {
		if this.MediaType != other.MediaType {
			return w
		}

		// We want e.g. home page in amp output format (media type text/html) to
		// find a template even if one isn't specified for that output format,
		// when one exist for the html output format (same media type).
		if category != CategoryBaseof && (this.Kind == "" || (this.Kind != other.Kind && (this.Layout != other.Layout && other.Layout != layoutAll))) {
			return w
		}

		// Continue.
	}

	// One example of variant1 and 2 is for render codeblocks:
	// variant1=codeblock, variant2=go (language).
	if other.Variant1 != "" && other.Variant1 != this.Variant1 {
		return w
	}

	if isEmbedded {
		if other.Variant2 != "" && other.Variant2 != this.Variant2 {
			return w
		}
	} else {
		// If both are set and different, no match.
		if other.Variant2 != "" && this.Variant2 != "" && other.Variant2 != this.Variant2 {
			return w
		}
	}

	const (
		weightKind         = 3 // page, home, section, taxonomy, term (and only those)
		weightcustomLayout = 4 // custom layout (mylayout, set in e.g. front matter)
		weightLayout       = 2 // standard layouts (single,list,all)
		weightOutputFormat = 2 // a configured output format (e.g. rss, html, json)
		weightMediaType    = 1 // a configured media type (e.g. text/html, text/plain)
		weightLang         = 1 // a configured language (e.g. en, nn, fr, ...)
		weightVariant1     = 4 // currently used for render hooks, e.g. "link", "image"
		weightVariant2     = 2 // currently used for render hooks, e.g. the language "go" in code blocks.

		// We will use the values for group 2 and 3
		// if the distance up to the template is shorter than
		// the one we're comparing with.
		// E.g for a page in /posts/mypage.md with the
		// two templates /layouts/posts/single.html and /layouts/page.html,
		// the first one is the best match even if the second one
		// has a higher w1 value.
		weight2Group1 = 1 // kind, standardl layout (single,list,all)
		weight2Group2 = 2 // custom layout (mylayout)

		weight3 = 1 // for media type, lang, output format.
	)

	// Now we now know that the other descriptor is a subset of this.
	// Now calculate the weights.
	w.w1++

	if other.Kind != "" && other.Kind == this.Kind {
		w.w1 += weightKind
		w.w2 = weight2Group1
	}

	if other.Layout != "" && other.Layout == this.Layout || other.Layout == layoutAll {
		if isLayoutCustom(this.Layout) {
			w.w1 += weightcustomLayout
			w.w2 = weight2Group2
		} else {
			w.w1 += weightLayout
			w.w2 = weight2Group1
		}
	}

	if other.Lang != "" && other.Lang == this.Lang {
		w.w1 += weightLang
		w.w3 += weight3
	}

	if other.OutputFormat != "" && other.OutputFormat == this.OutputFormat {
		w.w1 += weightOutputFormat
		w.w3 += weight3
	}

	if other.MediaType != "" && other.MediaType == this.MediaType {
		w.w1 += weightMediaType
		w.w3 += weight3
	}

	if other.Variant1 != "" && other.Variant1 == this.Variant1 {
		w.w1 += weightVariant1
	}

	if other.Variant1 != "" && other.Variant2 == this.Variant2 {
		w.w1 += weightVariant2
	}

	return w
}

func (d TemplateDescriptor) IsZero() bool {
	return d == TemplateDescriptor{}
}

//lint:ignore ST1006 this vs other makes it easier to reason about.
func (this TemplateDescriptor) isKindInLayout(layout string) bool {
	if this.Kind == "" {
		return true
	}
	if this.Kind != kinds.KindPage {
		return layout != layoutSingle
	}
	return layout != layoutList
}
