---
title: Sites
description: Returns a collection of all Site objects, one for each language, ordered by language weight.
categories: []
keywords: []
action:
  related:
    - methods/page/Site
  returnType: page.Sites
  signatures: [PAGE.Sites]
---

This is a convenience method to access `.Site.Sites`.

With this site configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = false

[languages.de]
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
title = 'Projekt Dokumentation'
weight = 1

[languages.en]
languageCode = 'en-US'
languageDirection = 'ltr'
languageName = 'English'
title = 'Project Documentation'
weight = 2
{{< /code-toggle >}}

This template:

```go-html-template
<ul>
  {{ range .Sites }}
    <li><a href="{{ .Home.Permalink }}">{{ .Title }}</a></li>
  {{ end }}
</ul>
```

Produces a list of links to each home page:

```html
<ul>
  <li><a href="https://example.org/de/">Projekt Dokumentation</a></li>
  <li><a href="https://example.org/en/">Project Documentation</a></li>
</ul>
```

To render a link to the home page of the site corresponding to the default content language:

```go-html-template
{{ with .Sites.Default }}
  <a href="{{ .Home.Permalink }}">{{ .Title }}</a>
{{ end }}
```

This is equivalent to:

```go-html-template
{{ with index .Sites 0 }}
  <a href="{{ .Home.Permalink }}">{{ .Title }}</a>
{{ end }}
```
