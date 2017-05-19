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

// Package release implements a set of utilities and a wrapper around Goreleaser
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
	issueLinkTemplate            = "[#%d](https://github.com/spf13/hugo/issues/%d)"
	linkTemplate                 = "[%s](%s)"
	releaseNotesMarkdownTemplate = `
{{- $patchRelease := isPatch . -}}
{{- $contribsPerAuthor := .All.ContribCountPerAuthor -}}

{{- if $patchRelease }}
{{ if eq (len .All) 1 }}
This is a bug-fix release with one important fix.
{{ else }}
This is a bug-fix relase with a couple of important fixes.
{{ end }}
{{ else }}
This release represents **{{ len .All }} contributions by {{ len $contribsPerAuthor }} contributors** to the main Hugo code base.
{{ end -}}

{{- if  gt (len $contribsPerAuthor) 3 -}}
{{- $u1 := index $contribsPerAuthor 0 -}}
{{- $u2 := index $contribsPerAuthor 1 -}}
{{- $u3 := index $contribsPerAuthor 2 -}}
{{- $u4 := index $contribsPerAuthor 3 -}}
{{- $u1.AuthorLink }} leads the Hugo development with a significant amount of contributions, but also a big shoutout to {{ $u2.AuthorLink }}, {{ $u3.AuthorLink }}, and {{ $u4.AuthorLink }} for their ongoing contributions.
And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.
{{ end }}
Hugo now has:

{{ with .Repo -}}
* {{ .Stars }}+ [stars](https://github.com/spf13/hugo/stargazers)
* {{ len .Contributors }}+ [contributors](https://github.com/spf13/hugo/graphs/contributors)
{{- end -}}
{{ with .ThemeCount }}
* 156+ [themes](http://themes.gohugo.io/)
{{- end }}
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
{{- $docsChanges := index . "docsChanges" -}}
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
{{- with $docsChanges -}}
### Docs
{{ template "change-section"  . }}
{{- end -}}
{{- with $otherChanges -}}
### Other
{{ template "change-section"  . }}
{{- end -}}
{{ end }}


{{ define "change-section" }}
{{ range . }}
{{- if .GitHubCommit -}}
* {{ .Subject }} {{ . | commitURL }} {{ . | authorURL }} {{ range .Issues }}{{ . | issue }} {{ end }}
{{ else -}}
* {{ .Subject }} {{ range .Issues }}{{ . | issue }} {{ end }}
{{ end -}}
{{- end }}
{{ end }}
`
)

var templateFuncs = template.FuncMap{
	"isPatch": func(c changeLog) bool {
		return strings.Count(c.Version, ".") > 1
	},
	"issue": func(id int) string {
		return fmt.Sprintf(issueLinkTemplate, id, id)
	},
	"commitURL": func(info gitInfo) string {
		if info.GitHubCommit.HtmlURL == "" {
			return ""
		}
		return fmt.Sprintf(linkTemplate, info.Hash, info.GitHubCommit.HtmlURL)
	},
	"authorURL": func(info gitInfo) string {
		if info.GitHubCommit.Author.Login == "" {
			return ""
		}
		return fmt.Sprintf(linkTemplate, "@"+info.GitHubCommit.Author.Login, info.GitHubCommit.Author.HtmlURL)
	},
}

func writeReleaseNotes(version string, infos gitInfos, to io.Writer) error {
	changes := gitInfosToChangeLog(infos)
	changes.Version = version
	repo, err := fetchRepo()
	if err == nil {
		changes.Repo = &repo
	}
	themeCount, err := fetchThemeCount()
	if err == nil {
		changes.ThemeCount = themeCount
	}

	tmpl, err := template.New("").Funcs(templateFuncs).Parse(releaseNotesMarkdownTemplate)
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
	resp, err := http.Get("https://github.com/spf13/hugoThemes/blob/master/.gitmodules")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	return bytes.Count(b, []byte("submodule")), nil
}

func writeReleaseNotesToTmpFile(version string, infos gitInfos) (string, error) {
	f, err := ioutil.TempFile("", "hugorelease")
	if err != nil {
		return "", err
	}

	defer f.Close()

	if err := writeReleaseNotes(version, infos, f); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func getRelaseNotesDocsTempDirAndName(version string) (string, string) {
	return hugoFilepath("docs/temp"), fmt.Sprintf("%s-relnotes.md", version)
}

func getRelaseNotesDocsTempFilename(version string) string {
	return filepath.Join(getRelaseNotesDocsTempDirAndName(version))
}

func writeReleaseNotesToDocsTemp(version string, infos gitInfos) (string, error) {
	docsTempPath, name := getRelaseNotesDocsTempDirAndName(version)
	os.Mkdir(docsTempPath, os.ModePerm)

	f, err := os.Create(filepath.Join(docsTempPath, name))
	if err != nil {
		return "", err
	}

	defer f.Close()

	if err := writeReleaseNotes(version, infos, f); err != nil {
		return "", err
	}

	return f.Name(), nil

}

func writeReleaseNotesToDocs(title, sourceFilename string) (string, error) {
	targetFilename := filepath.Base(sourceFilename)
	contentDir := hugoFilepath("docs/content/release-notes")
	targetFullFilename := filepath.Join(contentDir, targetFilename)
	os.Mkdir(contentDir, os.ModePerm)

	b, err := ioutil.ReadFile(sourceFilename)
	if err != nil {
		return "", err
	}

	f, err := os.Create(targetFullFilename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf(`
---
date: %s
title: %s
---

	`, time.Now().Format("2006-01-02"), title)); err != nil {
		return "", err
	}

	if _, err := f.Write(b); err != nil {
		return "", err
	}

	return targetFullFilename, nil

}
