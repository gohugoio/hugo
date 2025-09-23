// Package plaintext provides a passthrough/plain text markup converter.
package plaintext

import (
	"html"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/converter"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct{}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("plaintext", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &plainConverter{}, nil
	}), nil
}

type plainConverter struct{}

// Convert returns the input content HTML-escaped and wrapped in a <pre> tag
// to preserve whitespace and line breaks. No additional markup processing
// is performed.
func (c *plainConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	escaped := html.EscapeString(string(ctx.Src))
	out := []byte("<pre>" + escaped + "</pre>")
	return converter.Bytes(out), nil
}

func (c *plainConverter) Supports(feature identity.Identity) bool {
	return false
}
