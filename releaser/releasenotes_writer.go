// Copyright 2017-present The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package releaser implements a set of utilities and a wrapper around Goreleaser
// to help automate the Hugo release process.
package releaser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	issueLinkTemplate                        = "[#%d](https://github.com/gohugoio/hugo/issues/%d)"
	linkTemplate                             = "[%s](%s)"
	releaseNotesMarkdownTemplatePatchRelease = `
{{ if eq (len .All) 1 }}
This is a bug-fix release with one important fix.
{{ else }}
This is a bug-fix release with a couple of important fixes.
{{ end }}
{{ range .All }}
{{- if .GitHubCommit -}}
* {{ .Subject }} {{ . | commitURL }} {{ . | authorURL }} {{ range .Issues }}{{ . | issue }}{{ end }}
{{ else -}}
* {{ .Subject }} {{ range .Issues }}{{ . | issue }}{{ end }}
{{ end -}}
{{- end }}


`
	releaseNotesMarkdownTemplate = `
{{- $contribsPerAuthor := .All.ContribCountPerAuthor -}}
{{- $docsContribsPerAuthor := .Docs.ContribCountPerAuthor -}}

This release represents **{{ len .All }} contributions by {{ len $contribsPerAuthor }} contributors** to the main Hugo code base.

{{- if  gt (len $contribsPerAuthor) 3 -}}
{{- $u1 := index $contribsPerAuthor 0 -}}
{{- $u2 := index $contribsPerAuthor 1 -}}
{{- $u3 := index $contribsPerAuthor 2 -}}
{{- $u4 := index $contribsPerAuthor 3 -}}
{{- $u1.AuthorLink }} leads the Hugo development with a significant amount of contributions, but also a big shoutout to {{ $u2.AuthorLink }}, {{ $u3.AuthorLink }}, and {{ $u4.AuthorLink }} for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.
{{ end }}
Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **{{ len .Docs }} contributions by {{ len $docsContribsPerAuthor }} contributors**.
{{- if  gt (len $docsContribsPerAuthor) 3 -}}
{{- $u1 := index $docsContribsPerAuthor 0 -}}
{{- $u2 := index $docsContribsPerAuthor 1 -}}
{{- $u3 := index $docsContribsPerAuthor 2 -}}
{{- $u4 := index $docsContribsPerAuthor 3 }} A special thanks to {{ $u1.AuthorLink }}, {{ $u2.AuthorLink }}, {{ $u3.AuthorLink }}, and {{ $u4.AuthorLink }} for their work on the documentation site.
{{ end }}

Hugo now has:

{{ with .Repo -}}
* {{ .Stars }}+ [stars](https://github.com/gohugoio/hugo/stargazers)
* {{ len .Contributors }}+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
{{- end -}}
{{ with .ThemeCount }}
* {{ . }}+ [themes](http://themes.gohugo.io/)
{{ end }}
{{ with .Notes }}
## Notes
{{ template "change-section" . }}
{{- end -}}
## Enhancements
{{ template "change-headers"  .Enhancements -}}
## Fixes
{{ template "change-headers"  .Fixes -}}

{{ define "change-headers" }}
{{ $tmplChanges := index . "templateChanges" -}}
{{- $outChanges := index . "outChanges" -}}
{{- $coreChanges := index . "coreChanges" -}}
{{- $otherChanges := index . "otherChanges" -}}
{{- with $tmplChanges -}}
### Templates
{{ template "change-section" . }}
{{- end -}}
{{- with $outChanges -}}
### Output
{{ template "change-section"  . }}
{{- end -}}
{{- with $coreChanges -}}
### Core
{{ template "change-section" . }}
{{- end -}}
{{- with $otherChanges -}}
### Other
{{ template "change-section"  . }}
{{- end -}}
{{ end }}


{{ define "change-section" }}
{{ range . }}
{{- if .GitHubCommit -}}
* {{ .Subject }} {{ . | commitURL }} {{ . | authorURL }} {{ range .Issues }}{{ . | issue }}{{ end }}
{{ else -}}
* {{ .Subject }} {{ range .Issues }}{{ . | issue }}{{ end }}
{{ end -}}
{{- end }}
{{ end }}
`
)

var templateFuncs = template.FuncMap{
	"isPatch": func(c changeLog) bool {
		return !strings.HasSuffix(c.Version, "0")
	},
	"issue": func(id int) string {
		return fmt.Sprintf(issueLinkTemplate, id, id)
	},
	"commitURL": func(info gitInfo) string {
		if info.GitHubCommit.HTMLURL == "" {
			return ""
		}
		return fmt.Sprintf(linkTemplate, info.Hash, info.GitHubCommit.HTMLURL)
	},
	"authorURL": func(info gitInfo) string {
		if info.GitHubCommit.Author.Login == "" {
			return ""
		}
		return fmt.Sprintf(linkTemplate, "@"+info.GitHubCommit.Author.Login, info.GitHubCommit.Author.HTMLURL)
	},
}

func writeReleaseNotes(version string, infosMain, infosDocs gitInfos, to io.Writer) error {
	client := newGitHubAPI("hugo")
	changes := gitInfosToChangeLog(infosMain, infosDocs)
	changes.Version = version
	repo, err := client.fetchRepo()
	if err == nil {
		changes.Repo = &repo
	}
	themeCount, err := fetchThemeCount()
	if err == nil {
		changes.ThemeCount = themeCount
	}

	mtempl := releaseNotesMarkdownTemplate

	if !strings.HasSuffix(version, "0") {
		mtempl = releaseNotesMarkdownTemplatePatchRelease
	}

	tmpl, err := template.New("").Funcs(templateFuncs).Parse(mtempl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(to, changes)
	if err != nil {
		return err
	}

	return nil

}

func fetchThemeCount() (int, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/gohugoio/hugoThemes/master/.gitmodules")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	return bytes.Count(b, []byte("submodule")), nil
}

func writeReleaseNotesToTmpFile(version string, infosMain, infosDocs gitInfos) (string, error) {
	f, err := ioutil.TempFile("", "hugorelease")
	if err != nil {
		return "", err
	}

	defer f.Close()

	if err := writeReleaseNotes(version, infosMain, infosDocs, f); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func getReleaseNotesDocsTempDirAndName(version string, final bool) (string, string) {
	if final {
		return hugoFilepath("temp"), fmt.Sprintf("%s-relnotes-ready.md", version)
	}
	return hugoFilepath("temp"), fmt.Sprintf("%s-relnotes.md", version)
}

func getReleaseNotesDocsTempFilename(version string, final bool) string {
	return filepath.Join(getReleaseNotesDocsTempDirAndName(version, final))
}

func (r *ReleaseHandler) releaseNotesState(version string) (releaseNotesState, error) {
	docsTempPath, name := getReleaseNotesDocsTempDirAndName(version, false)
	_, err := os.Stat(filepath.Join(docsTempPath, name))

	if err == nil {
		return releaseNotesCreated, nil
	}

	docsTempPath, name = getReleaseNotesDocsTempDirAndName(version, true)
	_, err = os.Stat(filepath.Join(docsTempPath, name))

	if err == nil {
		return releaseNotesReady, nil
	}

	if !os.IsNotExist(err) {
		return releaseNotesNone, err
	}

	return releaseNotesNone, nil

}

func (r *ReleaseHandler) writeReleaseNotesToTemp(version string, isPatch bool, infosMain, infosDocs gitInfos) (string, error) {

	docsTempPath, name := getReleaseNotesDocsTempDirAndName(version, isPatch)

	var (
		w io.WriteCloser
	)

	if !r.try {
		os.Mkdir(docsTempPath, os.ModePerm)

		f, err := os.Create(filepath.Join(docsTempPath, name))
		if err != nil {
			return "", err
		}

		name = f.Name()

		defer f.Close()

		w = f

	} else {
		w = os.Stdout
	}

	if err := writeReleaseNotes(version, infosMain, infosDocs, w); err != nil {
		return "", err
	}

	return name, nil

}

func (r *ReleaseHandler) writeReleaseNotesToDocs(title, description, sourceFilename string) (string, error) {
	targetFilename := "index.md"
	bundleDir := strings.TrimSuffix(filepath.Base(sourceFilename), "-ready.md")
	contentDir := hugoFilepath("docs/content/en/news/" + bundleDir)
	targetFullFilename := filepath.Join(contentDir, targetFilename)

	if r.try {
		fmt.Printf("Write release notes to /docs: Bundle %q Dir: %q\n", bundleDir, contentDir)
		return targetFullFilename, nil
	}

	if err := os.MkdirAll(contentDir, os.ModePerm); err != nil {
		return "", nil
	}

	b, err := ioutil.ReadFile(sourceFilename)
	if err != nil {
		return "", err
	}

	f, err := os.Create(targetFullFilename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fmTail := ""
	if !strings.HasSuffix(title, ".0") {
		// Bug fix release
		fmTail = `
images:
- images/blog/hugo-bug-poster.png
`
	}

	if _, err := f.WriteString(fmt.Sprintf(`
---
date: %s
title: %q
description: %q
categories: ["Releases"]%s
---

	`, time.Now().Format("2006-01-02"), title, description, fmTail)); err != nil {
		return "", err
	}

	if _, err := f.Write(b); err != nil {
		return "", err
	}

	return targetFullFilename, nil

}
