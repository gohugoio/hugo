package hugolib

import (
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
	"github.com/spf13/viper"

	"io/ioutil"
	"os"

	"log"

	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/stretchr/testify/require"
)

func newTestPathSpec(fs *hugofs.Fs, v *viper.Viper) *helpers.PathSpec {
	l := helpers.NewDefaultLanguage(v)
	return helpers.NewPathSpec(fs, l)
}

func newTestCfg() (*viper.Viper, *hugofs.Fs) {

	v := viper.New()
	fs := hugofs.NewMem(v)

	v.SetFs(fs.Source)

	loadDefaultSettingsFor(v)

	// Default is false, but true is easier to use as default in tests
	v.Set("defaultContentLanguageInSubdir", true)

	return v, fs

}

// newTestSite creates a new site in the  English language with in-memory Fs.
// The site will have a template system loaded and ready to use.
// Note: This is only used in single site tests.
func newTestSite(t testing.TB, configKeyValues ...interface{}) *Site {

	cfg, fs := newTestCfg()

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}

	d := deps.DepsCfg{Language: helpers.NewLanguage("en", cfg), Fs: fs}

	s, err := NewSiteForCfg(d)

	if err != nil {
		t.Fatalf("Failed to create Site: %s", err)
	}
	return s
}

func newDebugLogger() *jww.Notepad {
	return jww.NewNotepad(jww.LevelDebug, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
}

func createWithTemplateFromNameValues(additionalTemplates ...string) func(templ tpl.Template) error {

	return func(templ tpl.Template) error {
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
	return buildSingleSiteExpected(t, false, depsCfg, buildCfg)
}

func buildSingleSiteExpected(t *testing.T, expectBuildError bool, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	h, err := NewHugoSites(depsCfg)

	require.NoError(t, err)
	require.Len(t, h.Sites, 1)

	if expectBuildError {
		require.Error(t, h.Build(buildCfg))
		return nil

	}

	require.NoError(t, h.Build(buildCfg))

	return h.Sites[0]
}

func writeSourcesToSource(t *testing.T, base string, fs *hugofs.Fs, sources ...source.ByteSource) {
	for _, src := range sources {
		writeSource(t, fs, filepath.Join(base, src.Name), string(src.Content))
	}
}
