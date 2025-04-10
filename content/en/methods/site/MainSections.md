---
title: MainSections
description: Returns a slice of the main section names as defined in the site configuration, falling back to the top-level section with the most pages.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: '[]string'
    signatures: [SITE.MainSections]
---

Site configuration:

{{< code-toggle file=hugo >}}
mainSections = ['books','films']
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.MainSections }} → [books films]
```

If `mainSections` is not defined in the site configuration, this method returns a slice with one element---the top-level section with the most pages.

With this content structure, the "films" section has the most pages:

```text
content/
├── books/
│   ├── book-1.md
│   └── book-2.md
├── films/
│   ├── film-1.md
│   ├── film-2.md
│   └── film-3.md
└── _index.md
```

Template:

```go-html-template
{{ .Site.MainSections }} → [films]
```

When creating a theme, instead of hardcoding section names when listing the most relevant pages on the front page, instruct site authors to set `mainSections` in their site configuration.

Then your home template can do something like this:

```go-html-template
{{ range where .Site.RegularPages "Section" "in" .Site.MainSections }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
