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
