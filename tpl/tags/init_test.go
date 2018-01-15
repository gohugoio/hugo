package tags_test

import (
	"testing"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
	"github.com/gohugoio/hugo/tpl/tags"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	var found bool
	var ns *internal.TemplateFuncsNamespace

	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns = nsf(&deps.Deps{})
		if tags.Name == ns.Name {
			found = true
			break
		}
	}

	require.True(t, found)
	require.IsType(t, &tags.Namespace{}, ns.Context())
}
