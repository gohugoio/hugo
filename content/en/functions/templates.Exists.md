---
title: templates.Exists
linktitle: ""
description: "Checks whether a template file exists under the given path relative to the `layouts` directory."
godocref: ""
date: 2018-11-01
publishdate: 2018-11-01
lastmod: 2018-11-01
categories: [functions]
tags: []
menu:
  docs:
    parent: "functions"
ns: ""
keywords: ["templates", "template", "layouts"]
signature: ["templates.Exists PATH"]
workson: []
hugoversion: "0.46"
aliases: []
relatedfuncs: []
toc: false
deprecated: false
---

A template file is any file living below the `layouts` directories of either the project or any of its theme components incudling partials and shortcodes.

The function is particularly handy with dynamic path. The following example ensures the build will not break on a `.Type` missing its dedicated `header` partial.

```go-html-template
{{ $partialPath := printf "headers/%s.html" .Type }}
{{ if templates.Exists ( printf "partials/%s" $partialPath ) }}
  {{ partial $partialPath . }}
{{ else }}
  {{ partial "headers/default.html" . }}
{{ end }}

```
