package hugolib

import (
	"path/filepath"
	"testing"

	"fmt"
	"regexp"
	"strings"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/viper"

	"io/ioutil"
	"os"

	"log"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/stretchr/testify/require"
)

const ()

type sitesBuilder struct {
	Cfg config.Provider
	Fs  *hugofs.Fs
	T   testing.TB

	H *HugoSites

	// We will add some default if not set.
	templatesAdded bool
	i18nAdded      bool
	dataAdded      bool
	contentAdded   bool
}

func newTestSitesBuilder(t testing.TB) *sitesBuilder {
	v := viper.New()
	fs := hugofs.NewMem(v)

	return &sitesBuilder{T: t, Fs: fs}
}

func (s *sitesBuilder) WithTOMLConfig(conf string) *sitesBuilder {
	writeSource(s.T, s.Fs, "config.toml", conf)
	return s
}

func (s *sitesBuilder) WithDefaultMultiSiteConfig() *sitesBuilder {
	var defaultMultiSiteConfig = `
baseURL = "http://example.com/blog"

paginate = 1
disablePathToLower = true
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true

[permalinks]
other = "/somewhere/else/:filename"

[blackfriday]
angledQuotes = true

[Taxonomies]
tag = "tags"

[Languages]
[Languages.en]
weight = 10
title = "In English"
languageName = "English"
[Languages.en.blackfriday]
angledQuotes = false
[[Languages.en.menu.main]]
url    = "/"
name   = "Home"
weight = 0

[Languages.fr]
weight = 20
title = "Le Français"
languageName = "Français"
[Languages.fr.Taxonomies]
plaque = "plaques"

[Languages.nn]
weight = 30
title = "På nynorsk"
languageName = "Nynorsk"
paginatePath = "side"
[Languages.nn.Taxonomies]
lag = "lag"
[[Languages.nn.menu.main]]
url    = "/"
name   = "Heim"
weight = 1

[Languages.nb]
weight = 40
title = "På bokmål"
languageName = "Bokmål"
paginatePath = "side"
[Languages.nb.Taxonomies]
lag = "lag"
`

	return s.WithTOMLConfig(defaultMultiSiteConfig)

}

func (s *sitesBuilder) WithContent(filenameContent ...string) *sitesBuilder {
	s.contentAdded = true
	for i := 0; i < len(filenameContent); i += 2 {
		filename, content := filenameContent[i], filenameContent[i+1]
		writeSource(s.T, s.Fs, filepath.Join("content", filename), content)
	}
	return s
}

func (s *sitesBuilder) WithTemplates(filenameContent ...string) *sitesBuilder {
	s.templatesAdded = true
	for i := 0; i < len(filenameContent); i += 2 {
		filename, content := filenameContent[i], filenameContent[i+1]
		writeSource(s.T, s.Fs, filepath.Join("layouts", filename), content)
	}
	return s
}

func (s *sitesBuilder) CreateSites() *sitesBuilder {
	if !s.templatesAdded {
		s.addDefaultTemplates()
	}
	if !s.i18nAdded {
		s.addDefaultI18n()
	}
	if !s.dataAdded {
		s.addDefaultData()
	}
	if !s.contentAdded {
		s.addDefaultContent()
	}

	if s.Cfg == nil {
		cfg, err := LoadConfig(s.Fs.Source, "", "config.toml")
		if err != nil {
			s.T.Fatalf("Failed to load config: %s", err)
		}
		s.Cfg = cfg
	}

	sites, err := NewHugoSites(deps.DepsCfg{Fs: s.Fs, Cfg: s.Cfg})
	if err != nil {
		s.T.Fatalf("Failed to create sites: %s", err)
	}
	s.H = sites

	return s
}

func (s *sitesBuilder) Build(cfg BuildCfg) *sitesBuilder {
	if s.H == nil {
		s.T.Fatal("Need to run builder.CreateSites first")
	}
	err := s.H.Build(cfg)
	if err != nil {
		s.T.Fatalf("Build failed: %s", err)
	}

	return s
}

func (s *sitesBuilder) addDefaultTemplates() {
	fs := s.Fs
	t := s.T

	// Layouts

	writeSource(t, fs, filepath.Join("layouts", "_default/single.html"), "Single: {{ .Title }}|{{ i18n \"hello\" }}|{{.Lang}}|{{ .Content }}")
	writeSource(t, fs, filepath.Join("layouts", "_default/list.html"), "{{ $p := .Paginator }}List Page {{ $p.PageNumber }}: {{ .Title }}|{{ i18n \"hello\" }}|{{ .Permalink }}|Pager: {{ template \"_internal/pagination.html\" . }}")
	writeSource(t, fs, filepath.Join("layouts", "index.html"), "{{ $p := .Paginator }}Default Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}")
	writeSource(t, fs, filepath.Join("layouts", "index.fr.html"), "{{ $p := .Paginator }}French Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}")

	// Shortcodes
	writeSource(t, fs, filepath.Join("layouts", "shortcodes", "shortcode.html"), "Shortcode: {{ i18n \"hello\" }}")
	// A shortcode in multiple languages
	writeSource(t, fs, filepath.Join("layouts", "shortcodes", "lingo.html"), "LingoDefault")
	writeSource(t, fs, filepath.Join("layouts", "shortcodes", "lingo.fr.html"), "LingoFrench")
}

func (s *sitesBuilder) addDefaultI18n() {
	fs := s.Fs
	t := s.T

	writeSource(t, fs, filepath.Join("i18n", "en.yaml"), `
hello:
  other: "Hello"
`)
	writeSource(t, fs, filepath.Join("i18n", "fr.yaml"), `
hello:
  other: "Bonjour"
`)

}

func (s *sitesBuilder) addDefaultData() {
	fs := s.Fs
	t := s.T

	writeSource(t, fs, filepath.FromSlash("data/hugo.toml"), "slogan = \"Hugo Rocks!\"")
}

func (s *sitesBuilder) addDefaultContent() {
	fs := s.Fs
	t := s.T

	contentTemplate := `---
title: doc1
weight: 1
tags:
 - tag1
date: "2018-02-28"
---
# doc1
*some "content"*

{{< shortcode >}}

{{< lingo >}}
`

	writeSource(t, fs, filepath.FromSlash("content/sect/doc1.en.md"), contentTemplate)
	writeSource(t, fs, filepath.FromSlash("content/sect/doc1.fr.md"), contentTemplate)
	writeSource(t, fs, filepath.FromSlash("content/sect/doc1.nb.md"), contentTemplate)
	writeSource(t, fs, filepath.FromSlash("content/sect/doc1.nn.md"), contentTemplate)
}

func (s *sitesBuilder) AssertFileContent(filename string, matches ...string) {
	content := readDestination(s.T, s.Fs, filename)
	for _, match := range matches {
		if !strings.Contains(content, match) {
			s.T.Fatalf("No match for %q in content for %s\n%q", match, filename, content)
		}
	}
}

func (s *sitesBuilder) AssertFileContentRe(filename string, matches ...string) {
	content := readDestination(s.T, s.Fs, filename)
	for _, match := range matches {
		r := regexp.MustCompile(match)
		if !r.MatchString(content) {
			s.T.Fatalf("No match for %q in content for %s\n%q", match, filename, content)
		}
	}
}

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

func newWarningLogger() *jww.Notepad {
	return jww.NewNotepad(jww.LevelWarn, jww.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
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

func writeSourcesToSource(t *testing.T, base string, fs *hugofs.Fs, sources ...[2]string) {
	for _, src := range sources {
		writeSource(t, fs, filepath.Join(base, src[0]), src[1])
	}
}

func dumpPages(pages ...*Page) {
	for i, p := range pages {
		fmt.Printf("%d: Kind: %s Title: %-10s RelPermalink: %-10s Path: %-10s sections: %s Len Sections(): %d\n",
			i+1,
			p.Kind, p.title, p.RelPermalink(), p.Path(), p.sections, len(p.Sections()))
	}
}

func isCI() bool {
	return os.Getenv("CI") != ""
}
