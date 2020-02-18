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

package converter

import (
	"bytes"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/markup_config"
	"github.com/gohugoio/hugo/markup/tableofcontents"
	"github.com/spf13/afero"
)

// ProviderConfig configures a new Provider.
type ProviderConfig struct {
	MarkupConfig markup_config.Config

	Cfg       config.Provider // Site config
	ContentFs afero.Fs
	Logger    *loggers.Logger
	Highlight func(code, lang, optsStr string) (string, error)
}

// ProviderProvider creates converter providers.
type ProviderProvider interface {
	New(cfg ProviderConfig) (Provider, error)
}

// Provider creates converters.
type Provider interface {
	New(ctx DocumentContext) (Converter, error)
	Name() string
}

// NewProvider creates a new Provider with the given name.
func NewProvider(name string, create func(ctx DocumentContext) (Converter, error)) Provider {
	return newConverter{
		name:   name,
		create: create,
	}
}

type newConverter struct {
	name   string
	create func(ctx DocumentContext) (Converter, error)
}

func (n newConverter) New(ctx DocumentContext) (Converter, error) {
	return n.create(ctx)
}

func (n newConverter) Name() string {
	return n.name
}

var NopConverter = new(nopConverter)

type nopConverter int

func (nopConverter) Convert(ctx RenderContext) (Result, error) {
	return &bytes.Buffer{}, nil
}

func (nopConverter) Supports(feature identity.Identity) bool {
	return false
}

// Converter wraps the Convert method that converts some markup into
// another format, e.g. Markdown to HTML.
type Converter interface {
	Convert(ctx RenderContext) (Result, error)
	Supports(feature identity.Identity) bool
}

// Result represents the minimum returned from Convert.
type Result interface {
	Bytes() []byte
}

// DocumentInfo holds additional information provided by some converters.
type DocumentInfo interface {
	AnchorSuffix() string
}

// TableOfContentsProvider provides the content as a ToC structure.
type TableOfContentsProvider interface {
	TableOfContents() tableofcontents.Root
}

// AnchorNameSanitizer tells how a converter sanitizes anchor names.
type AnchorNameSanitizer interface {
	SanitizeAnchorName(s string) string
}

// Bytes holds a byte slice and implements the Result interface.
type Bytes []byte

// Bytes returns itself
func (b Bytes) Bytes() []byte {
	return b
}

// DocumentContext holds contextual information about the document to convert.
type DocumentContext struct {
	Document        interface{} // May be nil. Usually a page.Page
	DocumentID      string
	DocumentName    string
	ConfigOverrides map[string]interface{}
}

// RenderContext holds contextual information about the content to render.
type RenderContext struct {
	Src         []byte
	RenderTOC   bool
	RenderHooks *hooks.Render
}

var (
	FeatureRenderHooks = identity.NewPathIdentity("markup", "renderingHooks")
)
