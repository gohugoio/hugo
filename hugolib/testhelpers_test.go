package hugolib

import (
	"path/filepath"
	"testing"

	"regexp"

	"fmt"
	"strings"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/viper"

	"io/ioutil"
	"os"

	"log"

	"github.com/gohugoio/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/stretchr/testify/require"
)

type testHelper struct {
	Cfg config.Provider
	Fs  *hugofs.Fs
	T   testing.TB
}

func (th testHelper) assertFileContent(filename string, matches ...string) {
	filename = th.replaceDefaultContentLanguageValue(filename)
	content := readDestination(th.T, th.Fs, filename)
	for _, match := range matches {
		match = th.replaceDefaultContentLanguageValue(match)
		require.True(th.T, strings.Contains(content, match), fmt.Sprintf("File no match for\n%q in\n%q:\n%s", strings.Replace(match, "%", "%%", -1), filename, strings.Replace(content, "%", "%%", -1)))
	}
}

// TODO(bep) better name for this. It does no magic replacements depending on defaultontentLanguageInSubDir.
func (th testHelper) assertFileContentStraight(filename string, matches ...string) {
	content := readDestination(th.T, th.Fs, filename)
	for _, match := range matches {
		require.True(th.T, strings.Contains(content, match), fmt.Sprintf("File no match for\n%q in\n%q:\n%s", strings.Replace(match, "%", "%%", -1), filename, strings.Replace(content, "%", "%%", -1)))
	}
}

func (th testHelper) assertFileContentRegexp(filename string, matches ...string) {
	filename = th.replaceDefaultContentLanguageValue(filename)
	content := readDestination(th.T, th.Fs, filename)
	for _, match := range matches {
		match = th.replaceDefaultContentLanguageValue(match)
		r := regexp.MustCompile(match)
		require.True(th.T, r.MatchString(content), fmt.Sprintf("File no match for\n%q in\n%q:\n%s", strings.Replace(match, "%", "%%", -1), filename, strings.Replace(content, "%", "%%", -1)))
	}
}

func (th testHelper) assertFileNotExist(filename string) {
	exists, err := helpers.Exists(filename, th.Fs.Destination)
	require.NoError(th.T, err)
	require.False(th.T, exists)
}

func (th testHelper) replaceDefaultContentLanguageValue(value string) string {
	defaultInSubDir := th.Cfg.GetBool("defaultContentLanguageInSubDir")
	replace := th.Cfg.GetString("defaultContentLanguage") + "/"

	if !defaultInSubDir {
		value = strings.Replace(value, replace, "", 1)

	}
	return value
}

func newTestPathSpec(fs *hugofs.Fs, v *viper.Viper) *helpers.PathSpec {
	l := helpers.NewDefaultLanguage(v)
	ps, _ := helpers.NewPathSpec(fs, l)
	return ps
}

func newTestDefaultPathSpec() *helpers.PathSpec {
	v := viper.New()
	// Easier to reason about in tests.
	v.Set("disablePathToLower", true)
	fs := hugofs.NewDefault(v)
	ps, _ := helpers.NewPathSpec(fs, v)
	return ps
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

	d := deps.DepsCfg{Language: helpers.NewLanguage("en", cfg), Fs: fs, Cfg: cfg}

	s, err := NewSiteForCfg(d)

	if err != nil {
		t.Fatalf("Failed to create Site: %s", err)
	}
	return s
}

func newTestSitesFromConfig(t testing.TB, afs afero.Fs, tomlConfig string, layoutPathContentPairs ...string) (testHelper, *HugoSites) {
	if len(layoutPathContentPairs)%2 != 0 {
		t.Fatalf("Layouts must be provided in pairs")
	}

	writeToFs(t, afs, "config.toml", tomlConfig)

	cfg, err := LoadConfig(afs, "", "config.toml")
	require.NoError(t, err)

	fs := hugofs.NewFrom(afs, cfg)
	th := testHelper{cfg, fs, t}

	for i := 0; i < len(layoutPathContentPairs); i += 2 {
		writeSource(t, fs, layoutPathContentPairs[i], layoutPathContentPairs[i+1])
	}

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	return th, h
}

func newTestSitesFromConfigWithDefaultTemplates(t testing.TB, tomlConfig string) (testHelper, *HugoSites) {
	return newTestSitesFromConfig(t, afero.NewMemMapFs(), tomlConfig,
		"layouts/_default/single.html", "Single|{{ .Title }}|{{ .Content }}",
		"layouts/_default/list.html", "List|{{ .Title }}|{{ .Content }}",
		"layouts/_default/terms.html", "Terms List|{{ .Title }}|{{ .Content }}",
	)
}

func newDebugLogger() *jww.Notepad {
	return jww.NewNotepad(jww.LevelDebug, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
}

func newErrorLogger() *jww.Notepad {
	return jww.NewNotepad(jww.LevelError, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
}
func createWithTemplateFromNameValues(additionalTemplates ...string) func(templ tpl.TemplateHandler) error {

	return func(templ tpl.TemplateHandler) error {
		for i := 0; i < len(additionalTemplates); i += 2 {
			err := templ.AddTemplate(additionalTemplates[i], additionalTemplates[i+1])
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func buildSingleSite(t testing.TB, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	return buildSingleSiteExpected(t, false, depsCfg, buildCfg)
}

func buildSingleSiteExpected(t testing.TB, expectBuildError bool, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
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

func isCI() bool {
	return os.Getenv("CI") != ""
}
