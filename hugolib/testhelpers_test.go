package hugolib

import (
	"image/jpeg"
	"io"
	"math/rand"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/google/go-cmp/cmp"

	"github.com/gohugoio/hugo/parser"
	"github.com/pkg/errors"

	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/sanity-io/litter"
	"github.com/spf13/afero"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl"
	"github.com/spf13/viper"

	"os"

	"github.com/gohugoio/hugo/resources/resource"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/hugofs"
)

var (
	deepEqualsPages         = qt.CmpEquals(cmp.Comparer(func(p1, p2 *pageState) bool { return p1 == p2 }))
	deepEqualsOutputFormats = qt.CmpEquals(cmp.Comparer(func(o1, o2 output.Format) bool {
		return o1.Name == o2.Name && o1.MediaType.Type() == o2.MediaType.Type()
	}))
)

type sitesBuilder struct {
	Cfg     config.Provider
	environ []string

	Fs      *hugofs.Fs
	T       testing.TB
	depsCfg deps.DepsCfg

	*qt.C

	logger *loggers.Logger
	rnd    *rand.Rand
	dumper litter.Options

	// Used to test partial rebuilds.
	changedFiles []string
	removedFiles []string

	// Aka the Hugo server mode.
	running bool

	H *HugoSites

	theme string

	// Default toml
	configFormat  string
	configFileSet bool
	viperSet      bool

	// Default is empty.
	// TODO(bep) revisit this and consider always setting it to something.
	// Consider this in relation to using the BaseFs.PublishFs to all publishing.
	workingDir string

	addNothing bool
	// Base data/content
	contentFilePairs  []filenameContent
	templateFilePairs []filenameContent
	i18nFilePairs     []filenameContent
	dataFilePairs     []filenameContent

	// Additional data/content.
	// As in "use the base, but add these on top".
	contentFilePairsAdded  []filenameContent
	templateFilePairsAdded []filenameContent
	i18nFilePairsAdded     []filenameContent
	dataFilePairsAdded     []filenameContent
}

type filenameContent struct {
	filename string
	content  string
}

func newTestSitesBuilder(t testing.TB) *sitesBuilder {
	v := viper.New()
	fs := hugofs.NewMem(v)

	litterOptions := litter.Options{
		HidePrivateFields: true,
		StripPackageNames: true,
		Separator:         " ",
	}

	return &sitesBuilder{T: t, C: qt.New(t), Fs: fs, configFormat: "toml",
		dumper: litterOptions, rnd: rand.New(rand.NewSource(time.Now().Unix()))}
}

func newTestSitesBuilderFromDepsCfg(t testing.TB, d deps.DepsCfg) *sitesBuilder {
	c := qt.New(t)

	litterOptions := litter.Options{
		HidePrivateFields: true,
		StripPackageNames: true,
		Separator:         " ",
	}

	b := &sitesBuilder{T: t, C: c, depsCfg: d, Fs: d.Fs, dumper: litterOptions, rnd: rand.New(rand.NewSource(time.Now().Unix()))}
	workingDir := d.Cfg.GetString("workingDir")

	b.WithWorkingDir(workingDir)

	return b.WithViper(d.Cfg.(*viper.Viper))

}

func (s *sitesBuilder) Running() *sitesBuilder {
	s.running = true
	return s
}

func (s *sitesBuilder) WithNothingAdded() *sitesBuilder {
	s.addNothing = true
	return s
}

func (s *sitesBuilder) WithLogger(logger *loggers.Logger) *sitesBuilder {
	s.logger = logger
	return s
}

func (s *sitesBuilder) WithWorkingDir(dir string) *sitesBuilder {
	s.workingDir = filepath.FromSlash(dir)
	return s
}

func (s *sitesBuilder) WithEnviron(env ...string) *sitesBuilder {
	for i := 0; i < len(env); i += 2 {
		s.environ = append(s.environ, fmt.Sprintf("%s=%s", env[i], env[i+1]))
	}
	return s
}

func (s *sitesBuilder) WithConfigTemplate(data interface{}, format, configTemplate string) *sitesBuilder {
	s.T.Helper()

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
	s.T.Helper()
	if s.configFileSet {
		s.T.Fatal("WithViper: use Viper or config.toml, not both")
	}
	defer func() {
		s.viperSet = true
	}()

	// Write to a config file to make sure the tests follow the same code path.
	var buff bytes.Buffer
	m := v.AllSettings()
	s.Assert(parser.InterfaceToConfig(m, metadecoders.TOML, &buff), qt.IsNil)
	return s.WithConfigFile("toml", buff.String())
}

func (s *sitesBuilder) WithConfigFile(format, conf string) *sitesBuilder {
	s.T.Helper()
	if s.viperSet {
		s.T.Fatal("WithConfigFile: use Viper or config.toml, not both")
	}
	s.configFileSet = true
	filename := s.absFilename("config." + format)
	writeSource(s.T, s.Fs, filename, conf)
	s.configFormat = format
	return s
}

func (s *sitesBuilder) WithThemeConfigFile(format, conf string) *sitesBuilder {
	s.T.Helper()
	if s.theme == "" {
		s.theme = "test-theme"
	}
	filename := filepath.Join("themes", s.theme, "config."+format)
	writeSource(s.T, s.Fs, s.absFilename(filename), conf)
	return s
}

func (s *sitesBuilder) WithSourceFile(filenameContent ...string) *sitesBuilder {
	s.T.Helper()
	for i := 0; i < len(filenameContent); i += 2 {
		writeSource(s.T, s.Fs, s.absFilename(filenameContent[i]), filenameContent[i+1])
	}
	return s
}

func (s *sitesBuilder) absFilename(filename string) string {
	filename = filepath.FromSlash(filename)
	if filepath.IsAbs(filename) {
		return filename
	}
	if s.workingDir != "" && !strings.HasPrefix(filename, s.workingDir) {
		filename = filepath.Join(s.workingDir, filename)
	}
	return filename
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
	s.T.Helper()
	return s.WithSimpleConfigFileAndBaseURL("http://example.com/")
}

func (s *sitesBuilder) WithSimpleConfigFileAndBaseURL(baseURL string) *sitesBuilder {
	s.T.Helper()
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

func (s *sitesBuilder) WithSunset(in string) {
	// Write a real image into one of the bundle above.
	src, err := os.Open(filepath.FromSlash("testdata/sunset.jpg"))
	s.Assert(err, qt.IsNil)

	out, err := s.Fs.Source.Create(filepath.FromSlash(filepath.Join(s.workingDir, in)))
	s.Assert(err, qt.IsNil)

	_, err = io.Copy(out, src)
	s.Assert(err, qt.IsNil)

	out.Close()
	src.Close()
}

func (s *sitesBuilder) createFilenameContent(pairs []string) []filenameContent {
	var slice []filenameContent
	s.appendFilenameContent(&slice, pairs...)
	return slice
}

func (s *sitesBuilder) appendFilenameContent(slice *[]filenameContent, pairs ...string) {
	if len(pairs)%2 != 0 {
		panic("file content mismatch")
	}
	for i := 0; i < len(pairs); i += 2 {
		c := filenameContent{
			filename: pairs[i],
			content:  pairs[i+1],
		}
		*slice = append(*slice, c)
	}
}

func (s *sitesBuilder) WithContent(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.contentFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithContentAdded(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.contentFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) WithTemplates(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.templateFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithTemplatesAdded(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.templateFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) WithData(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.dataFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithDataAdded(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.dataFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) WithI18n(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.i18nFilePairs, filenameContent...)
	return s
}

func (s *sitesBuilder) WithI18nAdded(filenameContent ...string) *sitesBuilder {
	s.appendFilenameContent(&s.i18nFilePairsAdded, filenameContent...)
	return s
}

func (s *sitesBuilder) EditFiles(filenameContent ...string) *sitesBuilder {
	for i := 0; i < len(filenameContent); i += 2 {
		filename, content := filepath.FromSlash(filenameContent[i]), filenameContent[i+1]
		absFilename := s.absFilename(filename)
		s.changedFiles = append(s.changedFiles, absFilename)
		writeSource(s.T, s.Fs, absFilename, content)

	}
	return s
}

func (s *sitesBuilder) RemoveFiles(filenames ...string) *sitesBuilder {
	for _, filename := range filenames {
		absFilename := s.absFilename(filename)
		s.removedFiles = append(s.removedFiles, absFilename)
		s.Assert(s.Fs.Source.Remove(absFilename), qt.IsNil)
	}
	return s
}

func (s *sitesBuilder) writeFilePairs(folder string, files []filenameContent) *sitesBuilder {
	// We have had some "filesystem ordering" bugs that we have not discovered in
	// our tests running with the in memory filesystem.
	// That file system is backed by a map so not sure how this helps, but some
	// randomness in tests doesn't hurt.
	// TODO(bep) this turns out to be more confusing than helpful.
	//s.rnd.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

	for _, fc := range files {
		target := folder
		// TODO(bep) clean  up this magic.
		if strings.HasPrefix(fc.filename, folder) {
			target = ""
		}

		if s.workingDir != "" {
			target = filepath.Join(s.workingDir, target)
		}

		writeSource(s.T, s.Fs, filepath.Join(target, fc.filename), fc.content)
	}
	return s
}

func (s *sitesBuilder) CreateSites() *sitesBuilder {
	if err := s.CreateSitesE(); err != nil {
		herrors.PrintStackTraceFromErr(err)
		s.Fatalf("Failed to create sites: %s", err)
	}

	return s
}

func (s *sitesBuilder) LoadConfig() error {
	if !s.configFileSet {
		s.WithSimpleConfigFile()
	}

	cfg, _, err := LoadConfig(ConfigSourceDescriptor{
		WorkingDir: s.workingDir,
		Fs:         s.Fs.Source,
		Logger:     s.logger,
		Environ:    s.environ,
		Filename:   "config." + s.configFormat}, func(cfg config.Provider) error {

		return nil
	})

	if err != nil {
		return err
	}

	s.Cfg = cfg

	return nil
}

func (s *sitesBuilder) CreateSitesE() error {
	if !s.addNothing {
		if _, ok := s.Fs.Source.(*afero.OsFs); ok {
			for _, dir := range []string{
				"content/sect",
				"layouts/_default",
				"layouts/_default/_markup",
				"layouts/partials",
				"layouts/shortcodes",
				"data",
				"i18n",
			} {
				if err := os.MkdirAll(filepath.Join(s.workingDir, dir), 0777); err != nil {
					return errors.Wrapf(err, "failed to create %q", dir)
				}
			}
		}

		s.addDefaults()
		s.writeFilePairs("content", s.contentFilePairsAdded)
		s.writeFilePairs("layouts", s.templateFilePairsAdded)
		s.writeFilePairs("data", s.dataFilePairsAdded)
		s.writeFilePairs("i18n", s.i18nFilePairsAdded)

		s.writeFilePairs("i18n", s.i18nFilePairs)
		s.writeFilePairs("data", s.dataFilePairs)
		s.writeFilePairs("content", s.contentFilePairs)
		s.writeFilePairs("layouts", s.templateFilePairs)

	}

	if err := s.LoadConfig(); err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	s.Fs.Destination = hugofs.NewCreateCountingFs(s.Fs.Destination)

	depsCfg := s.depsCfg
	depsCfg.Fs = s.Fs
	depsCfg.Cfg = s.Cfg
	depsCfg.Logger = s.logger
	depsCfg.Running = s.running

	sites, err := NewHugoSites(depsCfg)
	if err != nil {
		return errors.Wrap(err, "failed to create sites")
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
	s.T.Helper()
	return s.build(cfg, false)
}

func (s *sitesBuilder) BuildFail(cfg BuildCfg) *sitesBuilder {
	s.T.Helper()
	return s.build(cfg, true)
}

func (s *sitesBuilder) changeEvents() []fsnotify.Event {

	var events []fsnotify.Event

	for _, v := range s.changedFiles {
		events = append(events, fsnotify.Event{
			Name: v,
			Op:   fsnotify.Write,
		})
	}
	for _, v := range s.removedFiles {
		events = append(events, fsnotify.Event{
			Name: v,
			Op:   fsnotify.Remove,
		})
	}

	return events
}

func (s *sitesBuilder) build(cfg BuildCfg, shouldFail bool) *sitesBuilder {
	s.Helper()
	defer func() {
		s.changedFiles = nil
	}()

	if s.H == nil {
		s.CreateSites()
	}

	err := s.H.Build(cfg, s.changeEvents()...)

	if err == nil {
		logErrorCount := s.H.NumLogErrors()
		if logErrorCount > 0 {
			err = fmt.Errorf("logged %d errors", logErrorCount)
		}
	}
	if err != nil && !shouldFail {
		herrors.PrintStackTraceFromErr(err)
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

		listTemplateCommon = "{{ $p := .Paginator }}{{ $p.PageNumber }}|{{ .Title }}|{{ i18n \"hello\" }}|{{ .Permalink }}|Pager: {{ template \"_internal/pagination.html\" . }}|Kind: {{ .Kind }}|Content: {{ .Content }}|Len Pages: {{ len .Pages }}|Len RegularPages: {{ len .RegularPages }}| HasParent: {{ if .Parent }}YES{{ else }}NO{{ end }}"

		defaultTemplates = []string{
			"_default/single.html", "Single: {{ .Title }}|{{ i18n \"hello\" }}|{{.Language.Lang}}|RelPermalink: {{ .RelPermalink }}|Permalink: {{ .Permalink }}|{{ .Content }}|Resources: {{ range .Resources }}{{ .MediaType }}: {{ .RelPermalink}} -- {{ end }}|Summary: {{ .Summary }}|Truncated: {{ .Truncated }}|Parent: {{ .Parent.Title }}",
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
			// Special templates
			"404.html", "404|{{ .Lang }}|{{ .Title }}",
			"robots.txt", "robots|{{ .Lang }}|{{ .Title }}",
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
		s.writeFilePairs("content", s.createFilenameContent(defaultContent))
	}

	if len(s.templateFilePairs) == 0 {
		s.writeFilePairs("layouts", s.createFilenameContent(defaultTemplates))
	}
	if len(s.dataFilePairs) == 0 {
		s.writeFilePairs("data", s.createFilenameContent(defaultData))
	}
	if len(s.i18nFilePairs) == 0 {
		s.writeFilePairs("i18n", s.createFilenameContent(defaultI18n))
	}
}

func (s *sitesBuilder) Fatalf(format string, args ...interface{}) {
	s.T.Helper()
	s.T.Fatalf(format, args...)
}

func (s *sitesBuilder) AssertFileContentFn(filename string, f func(s string) bool) {
	s.T.Helper()
	content := s.FileContent(filename)
	if !f(content) {
		s.Fatalf("Assert failed for %q in content\n%s", filename, content)
	}
}

func (s *sitesBuilder) AssertHome(matches ...string) {
	s.AssertFileContent("public/index.html", matches...)
}

func (s *sitesBuilder) AssertFileContent(filename string, matches ...string) {
	s.T.Helper()
	content := s.FileContent(filename)
	for _, m := range matches {
		lines := strings.Split(m, "\n")
		for _, match := range lines {
			match = strings.TrimSpace(match)
			if match == "" {
				continue
			}
			if !strings.Contains(content, match) {
				s.Fatalf("No match for %q in content for %s\n%s\n%q", match, filename, content, content)
			}
		}
	}
}

func (s *sitesBuilder) AssertImage(width, height int, filename string) {
	filename = filepath.Join(s.workingDir, filename)
	f, err := s.Fs.Destination.Open(filename)
	s.Assert(err, qt.IsNil)
	defer f.Close()
	cfg, err := jpeg.DecodeConfig(f)
	s.Assert(err, qt.IsNil)
	s.Assert(cfg.Width, qt.Equals, width)
	s.Assert(cfg.Height, qt.Equals, height)
}

func (s *sitesBuilder) AssertNoDuplicateWrites() {
	s.Helper()
	d := s.Fs.Destination.(hugofs.DuplicatesReporter)
	s.Assert(d.ReportDuplicates(), qt.Equals, "")
}

func (s *sitesBuilder) FileContent(filename string) string {
	s.T.Helper()
	filename = filepath.FromSlash(filename)
	if !strings.HasPrefix(filename, s.workingDir) {
		filename = filepath.Join(s.workingDir, filename)
	}
	return readDestination(s.T, s.Fs, filename)
}

func (s *sitesBuilder) AssertObject(expected string, object interface{}) {
	s.T.Helper()
	got := s.dumper.Sdump(object)
	expected = strings.TrimSpace(expected)

	if expected != got {
		fmt.Println(got)
		diff := htesting.DiffStrings(expected, got)
		s.Fatalf("diff:\n%s\nexpected\n%s\ngot\n%s", diff, expected, got)
	}
}

func (s *sitesBuilder) AssertFileContentRe(filename string, matches ...string) {
	content := readDestination(s.T, s.Fs, filename)
	for _, match := range matches {
		r := regexp.MustCompile("(?s)" + match)
		if !r.MatchString(content) {
			s.Fatalf("No match for %q in content for %s\n%q", match, filename, content)
		}
	}
}

func (s *sitesBuilder) CheckExists(filename string) bool {
	return destinationExists(s.Fs, filepath.Clean(filename))
}

func (s *sitesBuilder) GetPage(ref string) page.Page {
	p, err := s.H.Sites[0].getPageNew(nil, ref)
	s.Assert(err, qt.IsNil)
	return p
}

func (s *sitesBuilder) GetPageRel(p page.Page, ref string) page.Page {
	p, err := s.H.Sites[0].getPageNew(p, ref)
	s.Assert(err, qt.IsNil)
	return p
}

func newTestHelper(cfg config.Provider, fs *hugofs.Fs, t testing.TB) testHelper {
	return testHelper{
		Cfg: cfg,
		Fs:  fs,
		C:   qt.New(t),
	}
}

type testHelper struct {
	Cfg config.Provider
	Fs  *hugofs.Fs
	*qt.C
}

func (th testHelper) assertFileContent(filename string, matches ...string) {
	th.Helper()
	filename = th.replaceDefaultContentLanguageValue(filename)
	content := readDestination(th, th.Fs, filename)
	for _, match := range matches {
		match = th.replaceDefaultContentLanguageValue(match)
		th.Assert(strings.Contains(content, match), qt.Equals, true, qt.Commentf(match+" not in: \n"+content))
	}
}

func (th testHelper) assertFileContentRegexp(filename string, matches ...string) {
	filename = th.replaceDefaultContentLanguageValue(filename)
	content := readDestination(th, th.Fs, filename)
	for _, match := range matches {
		match = th.replaceDefaultContentLanguageValue(match)
		r := regexp.MustCompile(match)
		matches := r.MatchString(content)
		if !matches {
			fmt.Println(match+":\n", content)
		}
		th.Assert(matches, qt.Equals, true)
	}
}

func (th testHelper) assertFileNotExist(filename string) {
	exists, err := helpers.Exists(filename, th.Fs.Destination)
	th.Assert(err, qt.IsNil)
	th.Assert(exists, qt.Equals, false)
}

func (th testHelper) replaceDefaultContentLanguageValue(value string) string {
	defaultInSubDir := th.Cfg.GetBool("defaultContentLanguageInSubDir")
	replace := th.Cfg.GetString("defaultContentLanguage") + "/"

	if !defaultInSubDir {
		value = strings.Replace(value, replace, "", 1)

	}
	return value
}

func loadTestConfig(fs afero.Fs, withConfig ...func(cfg config.Provider) error) (*viper.Viper, error) {
	v, _, err := LoadConfig(ConfigSourceDescriptor{Fs: fs}, withConfig...)
	return v, err
}

func newTestCfgBasic() (*viper.Viper, *hugofs.Fs) {
	mm := afero.NewMemMapFs()
	v := viper.New()
	v.Set("defaultContentLanguageInSubdir", true)

	fs := hugofs.NewFrom(hugofs.NewBaseFileDecorator(mm), v)

	return v, fs

}

func newTestCfg(withConfig ...func(cfg config.Provider) error) (*viper.Viper, *hugofs.Fs) {
	mm := afero.NewMemMapFs()

	v, err := loadTestConfig(mm, func(cfg config.Provider) error {
		// Default is false, but true is easier to use as default in tests
		cfg.Set("defaultContentLanguageInSubdir", true)

		for _, w := range withConfig {
			w(cfg)
		}

		return nil
	})

	if err != nil && err != ErrNoConfigFile {
		panic(err)
	}

	fs := hugofs.NewFrom(hugofs.NewBaseFileDecorator(mm), v)

	return v, fs

}

func newTestSitesFromConfig(t testing.TB, afs afero.Fs, tomlConfig string, layoutPathContentPairs ...string) (testHelper, *HugoSites) {
	if len(layoutPathContentPairs)%2 != 0 {
		t.Fatalf("Layouts must be provided in pairs")
	}

	c := qt.New(t)

	writeToFs(t, afs, filepath.Join("content", ".gitkeep"), "")
	writeToFs(t, afs, "config.toml", tomlConfig)

	cfg, err := LoadConfigDefault(afs)
	c.Assert(err, qt.IsNil)

	fs := hugofs.NewFrom(afs, cfg)
	th := newTestHelper(cfg, fs, t)

	for i := 0; i < len(layoutPathContentPairs); i += 2 {
		writeSource(t, fs, layoutPathContentPairs[i], layoutPathContentPairs[i+1])
	}

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	c.Assert(err, qt.IsNil)

	return th, h
}

func createWithTemplateFromNameValues(additionalTemplates ...string) func(templ tpl.TemplateManager) error {

	return func(templ tpl.TemplateManager) error {
		for i := 0; i < len(additionalTemplates); i += 2 {
			err := templ.AddTemplate(additionalTemplates[i], additionalTemplates[i+1])
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// TODO(bep) replace these with the builder
func buildSingleSite(t testing.TB, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	t.Helper()
	return buildSingleSiteExpected(t, false, false, depsCfg, buildCfg)
}

func buildSingleSiteExpected(t testing.TB, expectSiteInitEror, expectBuildError bool, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	t.Helper()
	b := newTestSitesBuilderFromDepsCfg(t, depsCfg).WithNothingAdded()

	err := b.CreateSitesE()

	if expectSiteInitEror {
		b.Assert(err, qt.Not(qt.IsNil))
		return nil
	} else {
		b.Assert(err, qt.IsNil)
	}

	h := b.H

	b.Assert(len(h.Sites), qt.Equals, 1)

	if expectBuildError {
		b.Assert(h.Build(buildCfg), qt.Not(qt.IsNil))
		return nil

	}

	b.Assert(h.Build(buildCfg), qt.IsNil)

	return h.Sites[0]
}

func writeSourcesToSource(t *testing.T, base string, fs *hugofs.Fs, sources ...[2]string) {
	for _, src := range sources {
		writeSource(t, fs, filepath.Join(base, src[0]), src[1])
	}
}

func getPage(in page.Page, ref string) page.Page {
	p, err := in.GetPage(ref)
	if err != nil {
		panic(err)
	}
	return p
}

func content(c resource.ContentProvider) string {
	cc, err := c.Content()
	if err != nil {
		panic(err)
	}

	ccs, err := cast.ToStringE(cc)
	if err != nil {
		panic(err)
	}
	return ccs
}

func dumpPages(pages ...page.Page) {
	fmt.Println("---------")
	for _, p := range pages {
		var meta interface{}
		if p.File() != nil && p.File().FileInfo() != nil {
			meta = p.File().FileInfo().Meta()
		}
		fmt.Printf("Kind: %s Title: %-10s RelPermalink: %-10s Path: %-10s sections: %s Lang: %s Meta: %v\n",
			p.Kind(), p.Title(), p.RelPermalink(), p.Path(), p.SectionsPath(), p.Lang(), meta)
	}
}

func dumpSPages(pages ...*pageState) {
	for i, p := range pages {
		fmt.Printf("%d: Kind: %s Title: %-10s RelPermalink: %-10s Path: %-10s sections: %s\n",
			i+1,
			p.Kind(), p.Title(), p.RelPermalink(), p.Path(), p.SectionsPath())
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
	return (os.Getenv("CI") != "" || os.Getenv("CI_LOCAL") != "") && os.Getenv("CIRCLE_BRANCH") == ""
}

// See https://github.com/golang/go/issues/19280
// Not in use.
var parallelEnabled = true

func parallel(t *testing.T) {
	if parallelEnabled {
		t.Parallel()
	}
}

func skipSymlink(t *testing.T) {
	if runtime.GOOS == "windows" && os.Getenv("CI") == "" {
		t.Skip("skip symlink test on local Windows (needs admin)")
	}

}

func captureStderr(f func() error) (string, error) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := f()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}
