package allconfig_test

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/media"
)

func TestDirsMount(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
disableKinds = ["taxonomy", "term"]
[languages]
[languages.en]
weight = 1
[languages.sv]
weight = 2
[[module.mounts]]
source = 'content/en'
target = 'content'
lang = 'en'
[[module.mounts]]
source = 'content/sv'
target = 'content'
lang = 'sv'
-- content/en/p1.md --
---
title: "p1"
---
-- content/sv/p1.md --
---
title: "p1"
---
-- layouts/_default/single.html --
Title: {{ .Title }}
	`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t, TxtarString: files},
	).Build()

	// b.AssertFileContent("public/p1/index.html", "Title: p1")

	sites := b.H.Sites
	b.Assert(len(sites), qt.Equals, 2)

	configs := b.H.Configs
	mods := configs.Modules
	b.Assert(len(mods), qt.Equals, 1)
	mod := mods[0]
	b.Assert(mod.Mounts(), qt.HasLen, 8)

	enConcp := sites[0].Conf
	enConf := enConcp.GetConfig().(*allconfig.Config)

	b.Assert(enConcp.BaseURL().String(), qt.Equals, "https://example.com/")
	modConf := enConf.Module
	b.Assert(modConf.Mounts, qt.HasLen, 8)
	b.Assert(modConf.Mounts[0].Source, qt.Equals, filepath.FromSlash("content/en"))
	b.Assert(modConf.Mounts[0].Target, qt.Equals, "content")
	b.Assert(modConf.Mounts[0].Lang, qt.Equals, "en")
	b.Assert(modConf.Mounts[1].Source, qt.Equals, filepath.FromSlash("content/sv"))
	b.Assert(modConf.Mounts[1].Target, qt.Equals, "content")
	b.Assert(modConf.Mounts[1].Lang, qt.Equals, "sv")
}

func TestConfigAliases(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
logI18nWarnings = true
logPathWarnings = true
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t, TxtarString: files},
	).Build()

	conf := b.H.Configs.Base

	b.Assert(conf.PrintI18nWarnings, qt.Equals, true)
	b.Assert(conf.PrintPathWarnings, qt.Equals, true)
}

func TestRedefineContentTypes(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"
[mediaTypes]
[mediaTypes."text/html"]
suffixes = ["html", "xhtml"]
`

	b := hugolib.Test(t, files)

	conf := b.H.Configs.Base
	contentTypes := conf.ContentTypes.Config

	b.Assert(contentTypes.HTML.Suffixes(), qt.DeepEquals, []string{"html", "xhtml"})
	b.Assert(contentTypes.Markdown.Suffixes(), qt.DeepEquals, []string{"md", "mdown", "markdown"})
}

func TestPaginationConfig(t *testing.T) {
	files := `
-- hugo.toml --
 [languages.en]
 weight = 1
 [languages.en.pagination]
 pagerSize = 20
 [languages.de]
 weight = 2
 [languages.de.pagination]
 path = "page-de"

`

	b := hugolib.Test(t, files)

	confEn := b.H.Sites[0].Conf.Pagination()
	confDe := b.H.Sites[1].Conf.Pagination()

	b.Assert(confEn.Path, qt.Equals, "page")
	b.Assert(confEn.PagerSize, qt.Equals, 20)
	b.Assert(confDe.Path, qt.Equals, "page-de")
	b.Assert(confDe.PagerSize, qt.Equals, 10)
}

func TestPaginationConfigDisableAliases(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[pagination]
disableAliases = true
pagerSize = 2
-- layouts/_default/list.html --
{{ $paginator := .Paginate  site.RegularPages }}
{{ template "_internal/pagination.html" . }}
{{ range $paginator.Pages }}
  {{ .Title }}
{{ end }}
-- content/p1.md --
---
title: "p1"
---
-- content/p2.md --
---
title: "p2"
---
-- content/p3.md --
---
title: "p3"
---
`

	b := hugolib.Test(t, files)

	b.AssertFileExists("public/page/1/index.html", false)
	b.AssertFileContent("public/page/2/index.html", "pagination-default")
}

func TestMapUglyURLs(t *testing.T) {
	files := `
-- hugo.toml --
[uglyurls]
  posts = true
`

	b := hugolib.Test(t, files)

	c := b.H.Configs.Base

	b.Assert(c.C.IsUglyURLSection("posts"), qt.IsTrue)
	b.Assert(c.C.IsUglyURLSection("blog"), qt.IsFalse)
}

// Issue 13199
func TestInvalidOutputFormat(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
[outputs]
home = ['html','foo']
-- layouts/index.html --
x
`

	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.IsNotNil)
	b.Assert(err.Error(), qt.Contains, `failed to create config: unknown output format "foo" for kind "home"`)
}

// Issue 13201
func TestLanguageConfigSlice(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
[languages.en]
title = 'TITLE_EN'
weight = 2
[languages.de]
title = 'TITLE_DE'
weight = 1
[languages.fr]
title = 'TITLE_FR'
weight = 3
`

	b := hugolib.Test(t, files)
	b.Assert(b.H.Configs.LanguageConfigSlice[0].Title, qt.Equals, `TITLE_DE`)
}

func TestContentTypesDefault(t *testing.T) {
	files := `
-- hugo.toml --
baseURL = "https://example.com"


`

	b := hugolib.Test(t, files)

	ct := b.H.Configs.Base.ContentTypes
	c := ct.Config
	s := ct.SourceStructure.(map[string]media.ContentTypeConfig)

	b.Assert(c.IsContentFile("foo.md"), qt.Equals, true)
	b.Assert(len(s), qt.Equals, 6)
}
