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

// Package asciidocext converts AsciiDoc to HTML using Asciidoctor
// external binary. The `asciidoc` module is reserved for a future golang
// implementation.
package asciidocext

import (
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/markup/asciidocext/internal"
	"github.com/gohugoio/hugo/markup/converter"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct{}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("asciidocext", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &internal.AsciiDocConverter{
			Ctx: ctx,
			Cfg: cfg,
		}, nil
	}), nil
}

// Supports reports whether the AsciiDoc converter is installed. Only used in
// tests.
func Supports() (bool, error) {
	hasAsciiDoc, err := internal.HasAsciiDoc()
	if htesting.SupportsAll() {
		if !hasAsciiDoc {
			panic(err)
		}
		return true, nil

	}

	return hasAsciiDoc, err
}

// SupportsGoATDiagrams reports whether the AsciiDoc converter can render GoAT
// diagrams. Only used in tests.
func SupportsGoATDiagrams() (bool, error) {
	supportsGoATDiagrams, err := internal.CanRenderGoATDiagrams()
	if htesting.SupportsAll() {
		if !supportsGoATDiagrams {
			panic(err)
		}
		return true, nil
	}

	return supportsGoATDiagrams, err
}
