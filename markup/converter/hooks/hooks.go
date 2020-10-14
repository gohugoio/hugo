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

package hooks

import (
	"io"

	"github.com/gohugoio/hugo/identity"
)

type LinkContext interface {
	Page() interface{}
	Destination() string
	Title() string
	Text() string
	PlainText() string
}

type LinkRenderer interface {
	RenderLink(w io.Writer, ctx LinkContext) error
	identity.Provider
}

// HeadingContext contains accessors to all attributes that a HeadingRenderer
// can use to render a heading.
type HeadingContext interface {
	// Page is the page containing the heading.
	Page() interface{}
	// Level is the level of the header (i.e. 1 for top-level, 2 for sub-level, etc.).
	Level() int
	// Anchor is the HTML id assigned to the heading.
	Anchor() string
	// Text is the rendered (HTML) heading text, excluding the heading marker.
	Text() string
	// PlainText is the unrendered version of Text.
	PlainText() string
}

// HeadingRenderer describes a uniquely identifiable rendering hook.
type HeadingRenderer interface {
	// Render writes the rendered content to w using the data in w.
	RenderHeading(w io.Writer, ctx HeadingContext) error
	identity.Provider
}

type Renderers struct {
	LinkRenderer    LinkRenderer
	ImageRenderer   LinkRenderer
	HeadingRenderer HeadingRenderer
}

func (r *Renderers) Eq(other interface{}) bool {
	ro, ok := other.(*Renderers)
	if !ok {
		return false
	}
	if r == nil || ro == nil {
		return r == nil
	}

	if r.ImageRenderer.GetIdentity() != ro.ImageRenderer.GetIdentity() {
		return false
	}

	if r.LinkRenderer.GetIdentity() != ro.LinkRenderer.GetIdentity() {
		return false
	}

	if r.HeadingRenderer.GetIdentity() != ro.HeadingRenderer.GetIdentity() {
		return false
	}

	return true
}
