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

package warpc

import (
	_ "embed"
)

//go:embed wasm/renderkatex.wasm
var katexWasm []byte

// See https://katex.org/docs/options.html
type KatexInput struct {
	Expression string       `json:"expression"`
	Options    KatexOptions `json:"options"`
}

// KatexOptions defines the options for the KaTeX rendering.
// See https://katex.org/docs/options.html
type KatexOptions struct {
	// html, mathml (default), htmlAndMathml
	Output string `json:"output"`

	// If true, display math in display mode, false in inline mode.
	DisplayMode bool `json:"displayMode"`

	// Render \tags on the left side instead of the right.
	Leqno bool `json:"leqno"`

	// If true,  render flush left with a 2em left margin.
	Fleqn bool `json:"fleqn"`

	// The color used for typesetting errors.
	// A color string given in the format "#XXX" or "#XXXXXX"
	ErrorColor string `json:"errorColor"`

	//  A collection of custom macros.
	Macros map[string]string `json:"macros,omitempty"`

	// Specifies a minimum thickness, in ems, for fraction lines.
	MinRuleThickness float64 `json:"minRuleThickness"`

	// If true, KaTeX will throw a ParseError when it encounters an unsupported command.
	ThrowOnError bool `json:"throwOnError"`
}

type KatexOutput struct {
	Output string `json:"output"`
}
