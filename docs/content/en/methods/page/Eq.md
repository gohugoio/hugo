---
title: Eq
description: Reports whether two Page objects are equal.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: bool
    signatures: [PAGE1.Eq PAGE2]
---

In this contrived example we list all pages in the current section except for the current page.

```go-html-template {file="layouts/page.html"}
{{ $currentPage := . }}
{{ range .CurrentSection.Pages }}
  {{ if not (.Eq $currentPage) }}
    <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
  {{ end }}
{{ end }}
```
