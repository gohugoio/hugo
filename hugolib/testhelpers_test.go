package hugolib

import (
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tplapi"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/require"
)

func newTestDepsConfig() deps.DepsCfg {
	return deps.DepsCfg{Fs: hugofs.NewMem()}
}

func newTestPathSpec() *helpers.PathSpec {
	return helpers.NewPathSpec(hugofs.NewMem(), viper.GetViper())
}

func createWithTemplateFromNameValues(additionalTemplates ...string) func(templ tplapi.Template) error {

	return func(templ tplapi.Template) error {
		for i := 0; i < len(additionalTemplates); i += 2 {
			err := templ.AddTemplate(additionalTemplates[i], additionalTemplates[i+1])
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func buildSingleSite(t *testing.T, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	h, err := NewHugoSitesFromConfiguration(depsCfg)

	require.NoError(t, err)
	require.Len(t, h.Sites, 1)

	require.NoError(t, h.Build(buildCfg))

	return h.Sites[0]
}

func writeSourcesToSource(t *testing.T, base string, fs *hugofs.Fs, sources ...source.ByteSource) {
	for _, src := range sources {
		writeSource(t, fs, filepath.Join(base, src.Name), string(src.Content))
	}
}
