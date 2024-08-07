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

type KatexOptions struct {
	Output       string `json:"output"` // html, mathml (default), htmlAndMathml
	DisplayMode  bool   `json:"displayMode"`
	ThrowOnError bool   `json:"throwOnError"`
}

type KatexOutput struct {
	Output string `json:"output"`
}
