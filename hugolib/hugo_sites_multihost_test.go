package hugolib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestMultihosts(t *testing.T) {
	t.Parallel()

	var multiSiteTOMLConfigTemplate = `
paginate = 1
disablePathToLower = true
defaultContentLanguage = "{{ .DefaultContentLanguage }}"
defaultContentLanguageInSubdir = {{ .DefaultContentLanguageInSubdir }}

[permalinks]
other = "/somewhere/else/:filename"

[Taxonomies]
tag = "tags"

[Languages]
[Languages.en]
baseURL = "https://example.com"
weight = 10
title = "In English"
languageName = "English"

[Languages.fr]
baseURL = "https://example.fr"
weight = 20
title = "Le Français"
languageName = "Français"

[Languages.nn]
baseURL = "https://example.no"
weight = 30
title = "På nynorsk"
languageName = "Nynorsk"

`

	siteConfig := testSiteConfig{Fs: afero.NewMemMapFs(), DefaultContentLanguage: "fr", DefaultContentLanguageInSubdir: false}
	sites := createMultiTestSites(t, siteConfig, multiSiteTOMLConfigTemplate)
	fs := sites.Fs
	cfg := BuildCfg{Watching: true}
	th := testHelper{sites.Cfg, fs, t}
	assert := require.New(t)

	err := sites.Build(cfg)
	assert.NoError(err)

	th.assertFileContent("public/en/sect/doc1-slug/index.html", "Hello")

	s1 := sites.Sites[0]

	s1h := s1.getPage(KindHome)
	assert.True(s1h.IsTranslated())
	assert.Len(s1h.Translations(), 2)
	assert.Equal("https://example.com/", s1h.Permalink())

	s2 := sites.Sites[1]
	s2h := s2.getPage(KindHome)
	assert.Equal("https://example.fr/", s2h.Permalink())

	th.assertFileContentStraight("public/fr/index.html", "French Home Page")
	th.assertFileContentStraight("public/en/index.html", "Default Home Page")

}
