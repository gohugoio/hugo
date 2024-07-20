package hugolib

import (
	"bytes"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
)

func TestContentFactory(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	c.Run("Simple", func(c *qt.C) {
		workingDir := "/my/work"
		b := newTestSitesBuilder(c)
		b.WithWorkingDir(workingDir).WithConfigFile("toml", `

workingDir="/my/work"

[module]
[[module.mounts]]
source = 'mcontent/en'
target = 'content'
lang  = 'en'
[[module.mounts]]
source = 'archetypes'
target = 'archetypes'
	
`)

		b.WithSourceFile(filepath.Join("mcontent/en/bundle", "index.md"), "")

		b.WithSourceFile(filepath.Join("archetypes", "post.md"), `---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: true
---

Hello World.
`)
		b.CreateSites()
		cf := NewContentFactory(b.H)
		abs, err := cf.CreateContentPlaceHolder(filepath.FromSlash("mcontent/en/blog/mypage.md"), false)
		b.Assert(err, qt.IsNil)
		b.Assert(abs, qt.Equals, filepath.FromSlash("/my/work/mcontent/en/blog/mypage.md"))
		b.Build(BuildCfg{SkipRender: true})

		p := b.H.GetContentPage(abs)
		b.Assert(p, qt.Not(qt.IsNil))

		var buf bytes.Buffer
		fi, err := b.H.BaseFs.Archetypes.Fs.Stat("post.md")
		b.Assert(err, qt.IsNil)
		b.Assert(cf.ApplyArchetypeFi(&buf, p, "", fi.(hugofs.FileMetaInfo)), qt.IsNil)

		b.Assert(buf.String(), qt.Contains, `title: "Mypage"`)
	})

	// Issue #9129
	c.Run("Content in both project and theme", func(c *qt.C) {
		b := newTestSitesBuilder(c)
		b.WithConfigFile("toml", `
theme = 'ipsum'		
`)

		themeDir := filepath.Join("themes", "ipsum")
		b.WithSourceFile("content/posts/foo.txt", `Hello.`)
		b.WithSourceFile(filepath.Join(themeDir, "content/posts/foo.txt"), `Hello.`)
		b.CreateSites()
		cf := NewContentFactory(b.H)
		abs, err := cf.CreateContentPlaceHolder(filepath.FromSlash("posts/test.md"), false)
		b.Assert(err, qt.IsNil)
		b.Assert(abs, qt.Equals, filepath.FromSlash("content/posts/test.md"))
	})
}
