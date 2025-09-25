// Package plaintext provides a passthrough/plain text markup converter.
package plaintext

import (
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

// Convert with no additional markup processing
func (c *plainConverter) Convert(ctx converter.RenderContext) (converter.ResultRender, error) {
	return converter.Bytes(ctx.Src), nil
}

func (c *plainConverter) Supports(feature identity.Identity) bool {
	return false
}
