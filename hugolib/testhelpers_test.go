package hugolib

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"unicode/utf8"

	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/gohugoio/hugo/langs"
	"github.com/sanity-io/litter"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/viper"

	"os"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sitesBuilder struct {
	Cfg config.Provider
	Fs  *hugofs.Fs
	T   testing.TB

	logger *loggers.Logger

	dumper litter.Options

	// Aka the Hugo server mode.
	running bool

	H *HugoSites

	theme string

	// Default toml
	configFormat string

	// Default is empty.
	// TODO(bep) revisit this and consider always setting it to something.
	// Consider this in relation to using the BaseFs.PublishFs to all publishing.
	workingDir string

	// Base data/content
	contentFilePairs  []string
	templateFilePairs []string
	i18nFilePairs     []string
	dataFilePairs     []string

	// Additional data/content.
	// As in "use the base, but add these on top".
	contentFilePairsAdded  []string
	templateFilePairsAdded []string
	i18nFilePairsAdded     []string
	dataFilePairsAdded     []string
}

func newTestSitesBuilder(t testing.TB) *sitesBuilder {
	v := viper.New()
	fs := hugofs.NewMem(v)

	litterOptions := litter.Options{
		HidePrivateFields: true,
		StripPackageNames: true,
		Separator:         " ",
	}

	return &sitesBuilder{T: t, Fs: fs, configFormat: "toml", dumper: litterOptions}
}

func createTempDir(prefix string) (string, func(), error) {
	workDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		return "", nil, err
	}

	if runtime.GOOS == "darwin" && !strings.HasPrefix(workDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		workDir = "/private" + workDir
	}
	return workDir, func() { os.RemoveAll(workDir) }, nil
}

func (s *sitesBuilder) Running() *sitesBuilder {
	s.running = true
	return s
}

func (s *sitesBuilder) WithLogger(logger *loggers.Logger) *sitesBuilder {
	s.logger = logger
	return s
}

func (s *sitesBuilder) WithWorkingDir(dir string) *sitesBuilder {
	s.workingDir = dir
	return s
}

func (s *sitesBuilder) WithConfigTemplate(data interface{}, format, configTemplate string) *sitesBuilder {
	if format == "" {
		format = "toml"
	}

	templ, err := template.New("test").Parse(configTemplate)
	if err != nil {
		s.Fatalf("Template parse failed: %s", err)
	}
	var b bytes.Buffer
	templ.Execute(&b, data)
	return s.WithConfigFile(format, b.String())
}

func (s *sitesBuilder) WithViper(v *viper.Viper) *sitesBuilder {
	loadDefaultSettingsFor(v)
	s.Cfg = v

	return s
}

func (s *sitesBuilder) WithConfigFile(format, conf string) *sitesBuilder {
	writeSource(s.T, s.Fs, "config."+format, conf)
	s.configFormat = format
	return s
}

func (s *sitesBuilder) WithThemeConfigFile(format, conf string) *sitesBuilder {
	if s.theme == "" {
		s.theme = "test-theme"
	}
	filename := filepath.Join("themes", s.theme, "config."+format)
	writeSource(s.T, s.Fs, filename, conf)
	return s
}

func (s *sitesBuilder) WithSourceFile(filename, content string) *sitesBuilder {
	writeSource(s.T, s.Fs, filepath.FromSlash(filename), content)
	return s
}

const commonConfigSections = `

[services]
[services.disqus]
shortname = "disqus_shortname"
[services.googleAnalytics]
id = "ga_id"

[privacy]
[privacy.disqus]
disable = false
[privacy.googleAnalytics]
respectDoNotTrack = true
anonymizeIP = true
[privacy.instagram]
simple = true
[privacy.twitter]
enableDNT = true
[privacy.vimeo]
disable = false
[privacy.youtube]
disable = false
privacyEnhanced = true

`

func (s *sitesBuilder) WithSimpleConfigFile() *sitesBuilder {
	return s.WithSimpleConfigFileAndBaseURL("http://example.com/")
}

func (s *sitesBuilder) WithSimpleConfigFileAndBaseURL(baseURL string) *sitesBuilder {
	config := fmt.Sprintf("baseURL = %q", baseURL)

	config = config + commonConfigSections
	return s.WithConfigFile("toml", config)
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
` + commonConfigSections

	return s.WithConfigFile("toml", defaultMultiSiteConfig)

}

func (s *sitesBuilder) WithContent(filenameContent ...string) *sitesBuilder {
	s.contentFilePairs = append(s.contentFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithContentAdded(filenameContent ...string) *sitesBuilder {
	s.contentFilePairsAdded = append(s.contentFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) WithTemplates(filenameContent ...string) *sitesBuilder {
	s.templateFilePairs = append(s.templateFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithTemplatesAdded(filenameContent ...string) *sitesBuilder {
	s.templateFilePairsAdded = append(s.templateFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) WithData(filenameContent ...string) *sitesBuilder {
	s.dataFilePairs = append(s.dataFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithDataAdded(filenameContent ...string) *sitesBuilder {
	s.dataFilePairsAdded = append(s.dataFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) WithI18n(filenameContent ...string) *sitesBuilder {
	s.i18nFilePairs = append(s.i18nFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithI18nAdded(filenameContent ...string) *sitesBuilder {
	s.i18nFilePairsAdded = append(s.i18nFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) writeFilePairs(folder string, filenameContent []string) *sitesBuilder {
	if len(filenameContent)%2 != 0 {
		s.Fatalf("expect filenameContent for %q in pairs (%d)", folder, len(filenameContent))
	}
	for i := 0; i < len(filenameContent); i += 2 {
		filename, content := filenameContent[i], filenameContent[i+1]
		target := folder
		// TODO(bep) clean  up this magic.
		if strings.HasPrefix(filename, folder) {
			target = ""
		}

		if s.workingDir != "" {
			target = filepath.Join(s.workingDir, target)
		}

		writeSource(s.T, s.Fs, filepath.Join(target, filename), content)
	}
	return s
}

func (s *sitesBuilder) CreateSites() *sitesBuilder {
	if err := s.CreateSitesE(); err != nil {
		s.Fatalf("Failed to create sites: %s", err)
	}

	return s
}

func (s *sitesBuilder) LoadConfig() error {
	cfg, _, err := LoadConfig(ConfigSourceDescriptor{Fs: s.Fs.Source, Filename: "config." + s.configFormat})
	if err != nil {
		return err
	}
	s.Cfg = cfg
	return nil
}

func (s *sitesBuilder) CreateSitesE() error {
	s.addDefaults()
	s.writeFilePairs("content", s.contentFilePairs)
	s.writeFilePairs("content", s.contentFilePairsAdded)
	s.writeFilePairs("layouts", s.templateFilePairs)
	s.writeFilePairs("layouts", s.templateFilePairsAdded)
	s.writeFilePairs("data", s.dataFilePairs)
	s.writeFilePairs("data", s.dataFilePairsAdded)
	s.writeFilePairs("i18n", s.i18nFilePairs)
	s.writeFilePairs("i18n", s.i18nFilePairsAdded)

	if s.Cfg == nil {
		if err := s.LoadConfig(); err != nil {
			return err
		}
	}

	sites, err := NewHugoSites(deps.DepsCfg{Fs: s.Fs, Cfg: s.Cfg, Logger: s.logger, Running: s.running})
	if err != nil {
		return err
	}
	s.H = sites

	return nil
}

func (s *sitesBuilder) BuildE(cfg BuildCfg) error {
	if s.H == nil {
		s.CreateSites()
	}

	return s.H.Build(cfg)
}

func (s *sitesBuilder) Build(cfg BuildCfg) *sitesBuilder {
	return s.build(cfg, false)
}

func (s *sitesBuilder) BuildFail(cfg BuildCfg) *sitesBuilder {
	return s.build(cfg, true)
}

func (s *sitesBuilder) build(cfg BuildCfg, shouldFail bool) *sitesBuilder {
	if s.H == nil {
		s.CreateSites()
	}

	err := s.H.Build(cfg)

	if err == nil {
		logErrorCount := s.H.NumLogErrors()
		if logErrorCount > 0 {
			err = fmt.Errorf("logged %d errors", logErrorCount)
		}
	}
	if err != nil && !shouldFail {
		s.Fatalf("Build failed: %s", err)
	} else if err == nil && shouldFail {
		s.Fatalf("Expected error")
	}

	return s
}

func (s *sitesBuilder) addDefaults() {

	var (
		contentTemplate = `---
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

		defaultContent = []string{
			"content/sect/doc1.en.md", contentTemplate,
			"content/sect/doc1.fr.md", contentTemplate,
			"content/sect/doc1.nb.md", contentTemplate,
			"content/sect/doc1.nn.md", contentTemplate,
		}

		listTemplateCommon = "{{ $p := .Paginator }}{{ $p.PageNumber }}|{{ .Title }}|{{ i18n \"hello\" }}|{{ .Permalink }}|Pager: {{ template \"_internal/pagination.html\" . }}"

		defaultTemplates = []string{
			"_default/single.html", "Single: {{ .Title }}|{{ i18n \"hello\" }}|{{.Lang}}|{{ .Content }}",
			"_default/list.html", "List Page " + listTemplateCommon,
			"index.html", "{{ $p := .Paginator }}Default Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}|String Resource: {{ ( \"Hugo Pipes\" | resources.FromString \"text/pipes.txt\").RelPermalink  }}",
			"index.fr.html", "{{ $p := .Paginator }}French Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}|String Resource: {{ ( \"Hugo Pipes\" | resources.FromString \"text/pipes.txt\").RelPermalink  }}",
			"_default/terms.html", "Taxonomy Term Page " + listTemplateCommon,
			"_default/taxonomy.html", "Taxonomy List Page " + listTemplateCommon,
			// Shortcodes
			"shortcodes/shortcode.html", "Shortcode: {{ i18n \"hello\" }}",
			// A shortcode in multiple languages
			"shortcodes/lingo.html", "LingoDefault",
			"shortcodes/lingo.fr.html", "LingoFrench",
		}

		defaultI18n = []string{
			"en.yaml", `
hello:
  other: "Hello"
`,
			"fr.yaml", `
hello:
  other: "Bonjour"
`,
		}

		defaultData = []string{
			"hugo.toml", "slogan = \"Hugo Rocks!\"",
		}
	)

	if len(s.contentFilePairs) == 0 {
		s.writeFilePairs("content", defaultContent)
	}
	if len(s.templateFilePairs) == 0 {
		s.writeFilePairs("layouts", defaultTemplates)
	}
	if len(s.dataFilePairs) == 0 {
		s.writeFilePairs("data", defaultData)
	}
	if len(s.i18nFilePairs) == 0 {
		s.writeFilePairs("i18n", defaultI18n)
	}
}

func (s *sitesBuilder) Fatalf(format string, args ...interface{}) {
	Fatalf(s.T, format, args...)
}

func Fatalf(t testing.TB, format string, args ...interface{}) {
	trace := trace()
	format = format + "\n%s"
	args = append(args, trace)
	t.Fatalf(format, args...)
}

func trace() string {
	return strings.Join(assert.CallerInfo(), "\n\r\t\t\t")
}

func (s *sitesBuilder) AssertFileContent(filename string, matches ...string) {
	content := readDestination(s.T, s.Fs, filename)
	for _, match := range matches {
		if !strings.Contains(content, match) {
			s.Fatalf("No match for %q in content for %s\n%s\n%q", match, filename, content, content)
		}
	}
}

func (s *sitesBuilder) AssertObject(expected string, object interface{}) {
	got := s.dumper.Sdump(object)
	expected = strings.TrimSpace(expected)

	if expected != got {
		fmt.Println(got)
		diff := helpers.DiffStrings(expected, got)
		s.Fatalf("diff:\n%s\nexpected\n%s\ngot\n%s", diff, expected, got)
	}
}

func (s *sitesBuilder) AssertFileContentRe(filename string, matches ...string) {
	content := readDestination(s.T, s.Fs, filename)
	for _, match := range matches {
		r := regexp.MustCompile(match)
		if !r.MatchString(content) {
			s.Fatalf("No match for %q in content for %s\n%q", match, filename, content)
		}
	}
}

func (s *sitesBuilder) CheckExists(filename string) bool {
	return destinationExists(s.Fs, filepath.Clean(filename))
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
	l := langs.NewDefaultLanguage(v)
	ps, _ := helpers.NewPathSpec(fs, l)
	return ps
}

func newTestDefaultPathSpec(t *testing.T) *helpers.PathSpec {
	v := viper.New()
	// Easier to reason about in tests.
	v.Set("disablePathToLower", true)
	v.Set("contentDir", "content")
	v.Set("dataDir", "data")
	v.Set("i18nDir", "i18n")
	v.Set("layoutDir", "layouts")
	v.Set("archetypeDir", "archetypes")
	v.Set("assetDir", "assets")
	v.Set("resourceDir", "resources")
	v.Set("publishDir", "public")
	fs := hugofs.NewDefault(v)
	ps, err := helpers.NewPathSpec(fs, v)
	if err != nil {
		t.Fatal(err)
	}
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

	d := deps.DepsCfg{Fs: fs, Cfg: cfg}

	s, err := NewSiteForCfg(d)

	if err != nil {
		Fatalf(t, "Failed to create Site: %s", err)
	}
	return s
}

func newTestSitesFromConfig(t testing.TB, afs afero.Fs, tomlConfig string, layoutPathContentPairs ...string) (testHelper, *HugoSites) {
	if len(layoutPathContentPairs)%2 != 0 {
		Fatalf(t, "Layouts must be provided in pairs")
	}

	writeToFs(t, afs, "config.toml", tomlConfig)

	cfg, err := LoadConfigDefault(afs)
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
	return buildSingleSiteExpected(t, false, false, depsCfg, buildCfg)
}

func buildSingleSiteExpected(t testing.TB, expectSiteInitEror, expectBuildError bool, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	h, err := NewHugoSites(depsCfg)

	if expectSiteInitEror {
		require.Error(t, err)
		return nil
	} else {
		require.NoError(t, err)
	}

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

func printStringIndexes(s string) {
	lines := strings.Split(s, "\n")
	i := 0

	for _, line := range lines {

		for _, r := range line {
			fmt.Printf("%-3s", strconv.Itoa(i))
			i += utf8.RuneLen(r)
		}
		i++
		fmt.Println()
		for _, r := range line {
			fmt.Printf("%-3s", string(r))
		}
		fmt.Println()

	}

}
func isCI() bool {
	return os.Getenv("CI") != ""
}

func isGo111() bool {
	return strings.Contains(runtime.Version(), "1.11")
}
