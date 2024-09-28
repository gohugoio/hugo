package hugolib

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/htesting"

	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/google/go-cmp/cmp"

	"github.com/gohugoio/hugo/parser"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/sanity-io/litter"
	"github.com/spf13/afero"
	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/resources/resource"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
)

var (
	deepEqualsPages         = qt.CmpEquals(cmp.Comparer(func(p1, p2 *pageState) bool { return p1 == p2 }))
	deepEqualsOutputFormats = qt.CmpEquals(cmp.Comparer(func(o1, o2 output.Format) bool {
		return o1.Name == o2.Name && o1.MediaType.Type == o2.MediaType.Type
	}))
)

type sitesBuilder struct {
	Cfg     config.Provider
	Configs *allconfig.Configs

	environ []string

	Fs      *hugofs.Fs
	T       testing.TB
	depsCfg deps.DepsCfg

	*qt.C

	logger loggers.Logger
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
	configSet     bool

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
	v := config.New()
	v.Set("publishDir", "public")
	v.Set("disableLiveReload", true)
	fs := hugofs.NewFromOld(afero.NewMemMapFs(), v)

	litterOptions := litter.Options{
		HidePrivateFields: true,
		StripPackageNames: true,
		Separator:         " ",
	}

	return &sitesBuilder{
		T: t, C: qt.New(t), Fs: fs, configFormat: "toml",
		dumper: litterOptions, rnd: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func newTestSitesBuilderFromDepsCfg(t testing.TB, d deps.DepsCfg) *sitesBuilder {
	c := qt.New(t)

	litterOptions := litter.Options{
		HidePrivateFields: true,
		StripPackageNames: true,
		Separator:         " ",
	}

	b := &sitesBuilder{T: t, C: c, depsCfg: d, Fs: d.Fs, dumper: litterOptions, rnd: rand.New(rand.NewSource(time.Now().Unix()))}
	workingDir := d.Configs.LoadingInfo.BaseConfig.WorkingDir

	b.WithWorkingDir(workingDir)

	return b
}

func (s *sitesBuilder) Running() *sitesBuilder {
	s.running = true
	return s
}

func (s *sitesBuilder) WithNothingAdded() *sitesBuilder {
	s.addNothing = true
	return s
}

func (s *sitesBuilder) WithLogger(logger loggers.Logger) *sitesBuilder {
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

func (s *sitesBuilder) WithConfigTemplate(data any, format, configTemplate string) *sitesBuilder {
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

func (s *sitesBuilder) WithViper(v config.Provider) *sitesBuilder {
	s.T.Helper()
	if s.configFileSet {
		s.T.Fatal("WithViper: use Viper or config.toml, not both")
	}
	defer func() {
		s.configSet = true
	}()

	// Write to a config file to make sure the tests follow the same code path.
	var buff bytes.Buffer
	m := v.Get("").(maps.Params)
	s.Assert(parser.InterfaceToConfig(m, metadecoders.TOML, &buff), qt.IsNil)
	return s.WithConfigFile("toml", buff.String())
}

func (s *sitesBuilder) WithConfigFile(format, conf string) *sitesBuilder {
	s.T.Helper()
	if s.configSet {
		s.T.Fatal("WithConfigFile: use config.Config or config.toml, not both")
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
id = "UA-ga_id"

[privacy]
[privacy.disqus]
disable = false
[privacy.googleAnalytics]
respectDoNotTrack = true
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
	return s.WithSimpleConfigFileAndSettings(map[string]any{"baseURL": baseURL})
}

func (s *sitesBuilder) WithSimpleConfigFileAndSettings(settings any) *sitesBuilder {
	s.T.Helper()
	var buf bytes.Buffer
	parser.InterfaceToConfig(settings, metadecoders.TOML, &buf)
	config := buf.String() + commonConfigSections
	return s.WithConfigFile("toml", config)
}

func (s *sitesBuilder) WithDefaultMultiSiteConfig() *sitesBuilder {
	defaultMultiSiteConfig := `
baseURL = "http://example.com/blog"

disablePathToLower = true
defaultContentLanguage = "en"
defaultContentLanguageInSubdir = true

[pagination]
pagerSize = 1

[permalinks]
other = "/somewhere/else/:filename"

[Taxonomies]
tag = "tags"

[Languages]
[Languages.en]
weight = 10
title = "In English"
languageName = "English"
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
[Languages.nn.pagination]
path = "side"
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
[Languages.nb.pagination]
path = "side"
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
	// s.rnd.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

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
		s.Fatalf("Failed to create sites: %s", err)
	}

	s.Assert(s.Fs.PublishDir, qt.IsNotNil)
	s.Assert(s.Fs.WorkingDirReadOnly, qt.IsNotNil)

	return s
}

func (s *sitesBuilder) LoadConfig() error {
	if !s.configFileSet {
		s.WithSimpleConfigFile()
	}

	flags := config.New()
	flags.Set("internal", map[string]any{
		"running": s.running,
		"watch":   s.running,
	})

	if s.workingDir != "" {
		flags.Set("workingDir", s.workingDir)
	}

	res, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{
		Fs:       s.Fs.Source,
		Logger:   s.logger,
		Flags:    flags,
		Environ:  s.environ,
		Filename: "config." + s.configFormat,
	})
	if err != nil {
		return err
	}

	s.Cfg = res.LoadingInfo.Cfg
	s.Configs = res

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
				if err := os.MkdirAll(filepath.Join(s.workingDir, dir), 0o777); err != nil {
					return fmt.Errorf("failed to create %q: %w", dir, err)
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
		return fmt.Errorf("failed to load config: %w", err)
	}

	s.Fs.PublishDir = hugofs.NewCreateCountingFs(s.Fs.PublishDir)

	depsCfg := s.depsCfg
	depsCfg.Fs = s.Fs
	if depsCfg.Configs.IsZero() {
		depsCfg.Configs = s.Configs
	}
	depsCfg.TestLogger = s.logger

	sites, err := NewHugoSites(depsCfg)
	if err != nil {
		return fmt.Errorf("failed to create sites: %w", err)
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
			"index.html", "{{ $p := .Paginator }}Default Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}|String Resource: {{ ( \"Hugo Pipes\" | resources.FromString \"text/pipes.txt\").RelPermalink  }}|String Resource Permalink: {{ ( \"Hugo Pipes\" | resources.FromString \"text/pipes.txt\").Permalink  }}",
			"index.fr.html", "{{ $p := .Paginator }}French Home Page {{ $p.PageNumber }}: {{ .Title }}|{{ .IsHome }}|{{ i18n \"hello\" }}|{{ .Permalink }}|{{  .Site.Data.hugo.slogan }}|String Resource: {{ ( \"Hugo Pipes\" | resources.FromString \"text/pipes.txt\").RelPermalink  }}|String Resource Permalink: {{ ( \"Hugo Pipes\" | resources.FromString \"text/pipes.txt\").Permalink  }}",
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

func (s *sitesBuilder) Fatalf(format string, args ...any) {
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

// Helper to migrate tests to new format.
func (s *sitesBuilder) DumpTxtar() string {
	var sb strings.Builder

	skipRe := regexp.MustCompile(`^(public|resources|package-lock.json|go.sum)`)

	afero.Walk(s.Fs.Source, s.workingDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, s.workingDir+"/")
		if skipRe.MatchString(rel) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info == nil || info.IsDir() {
			return nil
		}
		sb.WriteString(fmt.Sprintf("-- %s --\n", rel))
		b, err := afero.ReadFile(s.Fs.Source, path)
		s.Assert(err, qt.IsNil)
		sb.WriteString(strings.TrimSpace(string(b)))
		sb.WriteString("\n")
		return nil
	})

	return sb.String()
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
				s.Assert(content, qt.Contains, match, qt.Commentf(match+" not in: \n"+content))
			}
		}
	}
}

func (s *sitesBuilder) AssertFileDoesNotExist(filename string) {
	if s.CheckExists(filename) {
		s.Fatalf("File %q exists but must not exist.", filename)
	}
}

func (s *sitesBuilder) AssertImage(width, height int, filename string) {
	f, err := s.Fs.WorkingDirReadOnly.Open(filename)
	s.Assert(err, qt.IsNil)
	defer f.Close()
	cfg, err := jpeg.DecodeConfig(f)
	s.Assert(err, qt.IsNil)
	s.Assert(cfg.Width, qt.Equals, width)
	s.Assert(cfg.Height, qt.Equals, height)
}

func (s *sitesBuilder) AssertNoDuplicateWrites() {
	s.Helper()
	hugofs.WalkFilesystems(s.Fs.PublishDir, func(fs afero.Fs) bool {
		if dfs, ok := fs.(hugofs.DuplicatesReporter); ok {
			s.Assert(dfs.ReportDuplicates(), qt.Equals, "")
		}
		return false
	})
}

func (s *sitesBuilder) FileContent(filename string) string {
	s.Helper()
	filename = filepath.FromSlash(filename)
	return readWorkingDir(s.T, s.Fs, filename)
}

func (s *sitesBuilder) AssertObject(expected string, object any) {
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
	content := readWorkingDir(s.T, s.Fs, filename)
	for _, match := range matches {
		r := regexp.MustCompile("(?s)" + match)
		if !r.MatchString(content) {
			s.Fatalf("No match for %q in content for %s\n%q", match, filename, content)
		}
	}
}

func (s *sitesBuilder) CheckExists(filename string) bool {
	return workingDirExists(s.Fs, filepath.Clean(filename))
}

func (s *sitesBuilder) GetPage(ref string) page.Page {
	p, err := s.H.Sites[0].getPage(nil, ref)
	s.Assert(err, qt.IsNil)
	return p
}

func (s *sitesBuilder) GetPageRel(p page.Page, ref string) page.Page {
	p, err := s.H.Sites[0].getPage(p, ref)
	s.Assert(err, qt.IsNil)
	return p
}

func (s *sitesBuilder) NpmInstall() hexec.Runner {
	sc := security.DefaultConfig
	var err error
	sc.Exec.Allow, err = security.NewWhitelist("npm")
	s.Assert(err, qt.IsNil)
	ex := hexec.New(sc, s.workingDir)
	command, err := ex.New("npm", "install")
	s.Assert(err, qt.IsNil)
	return command
}

func newTestHelperFromProvider(cfg config.Provider, fs *hugofs.Fs, t testing.TB) (testHelper, *allconfig.Configs) {
	res, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{
		Flags: cfg,
		Fs:    fs.Source,
	})
	if err != nil {
		t.Fatal(err)
	}
	return newTestHelper(res.Base, fs, t), res
}

func newTestHelper(cfg *allconfig.Config, fs *hugofs.Fs, t testing.TB) testHelper {
	return testHelper{
		Cfg: cfg,
		Fs:  fs,
		C:   qt.New(t),
	}
}

type testHelper struct {
	Cfg *allconfig.Config
	Fs  *hugofs.Fs
	*qt.C
}

func (th testHelper) assertFileContent(filename string, matches ...string) {
	th.Helper()
	filename = th.replaceDefaultContentLanguageValue(filename)
	content := readWorkingDir(th, th.Fs, filename)
	for _, match := range matches {
		match = th.replaceDefaultContentLanguageValue(match)
		th.Assert(strings.Contains(content, match), qt.Equals, true, qt.Commentf(match+" not in: \n"+content))
	}
}

func (th testHelper) assertFileNotExist(filename string) {
	exists, err := helpers.Exists(filename, th.Fs.PublishDir)
	th.Assert(err, qt.IsNil)
	th.Assert(exists, qt.Equals, false)
}

func (th testHelper) replaceDefaultContentLanguageValue(value string) string {
	defaultInSubDir := th.Cfg.DefaultContentLanguageInSubdir
	replace := th.Cfg.DefaultContentLanguage + "/"

	if !defaultInSubDir {
		value = strings.Replace(value, replace, "", 1)
	}
	return value
}

func loadTestConfigFromProvider(cfg config.Provider) (*allconfig.Configs, error) {
	workingDir := cfg.GetString("workingDir")
	fs := afero.NewMemMapFs()
	if workingDir != "" {
		fs.MkdirAll(workingDir, 0o755)
	}
	res, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{Flags: cfg, Fs: fs})
	return res, err
}

func newTestCfg(withConfig ...func(cfg config.Provider) error) (config.Provider, *hugofs.Fs) {
	mm := afero.NewMemMapFs()
	cfg := config.New()
	cfg.Set("defaultContentLanguageInSubdir", false)
	cfg.Set("publishDir", "public")

	fs := hugofs.NewFromOld(hugofs.NewBaseFileDecorator(mm), cfg)

	return cfg, fs
}

func newTestSitesFromConfig(t testing.TB, afs afero.Fs, tomlConfig string, layoutPathContentPairs ...string) (testHelper, *HugoSites) {
	if len(layoutPathContentPairs)%2 != 0 {
		t.Fatalf("Layouts must be provided in pairs")
	}

	c := qt.New(t)

	writeToFs(t, afs, filepath.Join("content", ".gitkeep"), "")
	writeToFs(t, afs, "config.toml", tomlConfig)

	cfg, err := allconfig.LoadConfig(allconfig.ConfigSourceDescriptor{Fs: afs})
	c.Assert(err, qt.IsNil)

	fs := hugofs.NewFrom(afs, cfg.LoadingInfo.BaseConfig)
	th := newTestHelper(cfg.Base, fs, t)

	for i := 0; i < len(layoutPathContentPairs); i += 2 {
		writeSource(t, fs, layoutPathContentPairs[i], layoutPathContentPairs[i+1])
	}

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Configs: cfg})

	c.Assert(err, qt.IsNil)

	return th, h
}

// TODO(bep) replace these with the builder
func buildSingleSite(t testing.TB, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	t.Helper()
	return buildSingleSiteExpected(t, false, false, depsCfg, buildCfg)
}

func buildSingleSiteExpected(t testing.TB, expectSiteInitError, expectBuildError bool, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	t.Helper()
	b := newTestSitesBuilderFromDepsCfg(t, depsCfg).WithNothingAdded()

	err := b.CreateSitesE()

	if expectSiteInitError {
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
	cc, err := c.Content(context.Background())
	if err != nil {
		panic(err)
	}

	ccs, err := cast.ToStringE(cc)
	if err != nil {
		panic(err)
	}
	return ccs
}
