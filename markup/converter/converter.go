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
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"
)

// ProviderConfig configures a new Provider.
type ProviderConfig struct {
	Cfg       config.Provider // Site config
	ContentFs afero.Fs
	Logger    *loggers.Logger
	Highlight func(code, lang, optsStr string) (string, error)
}

// NewProvider creates converter providers.
type NewProvider interface {
	New(cfg ProviderConfig) (Provider, error)
}

// Provider creates converters.
type Provider interface {
	New(ctx DocumentContext) (Converter, error)
}

// NewConverter is an adapter that can be used as a ConverterProvider.
type NewConverter func(ctx DocumentContext) (Converter, error)

// New creates a new Converter for the given ctx.
func (n NewConverter) New(ctx DocumentContext) (Converter, error) {
	return n(ctx)
}

// Converter wraps the Convert method that converts some markup into
// another format, e.g. Markdown to HTML.
type Converter interface {
	Convert(ctx RenderContext) (Result, error)
}

// Result represents the minimum returned from Convert.
type Result interface {
	Bytes() []byte
}

// DocumentInfo holds additional information provided by some converters.
type DocumentInfo interface {
	AnchorSuffix() string
}

// Bytes holds a byte slice and implements the Result interface.
type Bytes []byte

// Bytes returns itself
func (b Bytes) Bytes() []byte {
	return b
}

// DocumentContext holds contextual information about the document to convert.
type DocumentContext struct {
	DocumentID      string
	DocumentName    string
	ConfigOverrides map[string]interface{}
}

// RenderContext holds contextual information about the content to render.
type RenderContext struct {
	Src       []byte
	RenderTOC bool
}
