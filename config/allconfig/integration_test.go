package allconfig_test

import (
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/hugolib"
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

	//b.AssertFileContent("public/p1/index.html", "Title: p1")

	sites := b.H.Sites
	b.Assert(len(sites), qt.Equals, 2)

	configs := b.H.Configs
	mods := configs.Modules
	b.Assert(len(mods), qt.Equals, 1)
	mod := mods[0]
	b.Assert(mod.Mounts(), qt.HasLen, 8)

	enConcp := sites[0].Conf
	enConf := enConcp.GetConfig().(*allconfig.Config)

	b.Assert(enConcp.BaseURL().String(), qt.Equals, "https://example.com")
	modConf := enConf.Module
	b.Assert(modConf.Mounts, qt.HasLen, 8)
	b.Assert(modConf.Mounts[0].Source, qt.Equals, filepath.FromSlash("content/en"))
	b.Assert(modConf.Mounts[0].Target, qt.Equals, "content")
	b.Assert(modConf.Mounts[0].Lang, qt.Equals, "en")
	b.Assert(modConf.Mounts[1].Source, qt.Equals, filepath.FromSlash("content/sv"))
	b.Assert(modConf.Mounts[1].Target, qt.Equals, "content")
	b.Assert(modConf.Mounts[1].Lang, qt.Equals, "sv")

}
