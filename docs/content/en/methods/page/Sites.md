---
title: Sites
description: Returns a collection of all sites for all dimensions.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Sites
    signatures: [PAGE.Sites]
expiryDate: '2028-02-18' # deprecated 2026-02-18 in v0.156.0
---

{{< deprecated-in 0.156.0 >}}
Use [`hugo.Sites`] instead.

[`hugo.Sites`]: /functions/hugo/sites/
{{< /deprecated-in >}}

{{% include "/_common/functions/hugo/sites-collection.md" %}}

With this project configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true
defaultContentVersionInSubdir = true

[languages.de]
contentDir = 'content/de'
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
title = 'Projekt Dokumentation'
weight = 1

[languages.en]
contentDir = 'content/en'
languageCode = 'en-US'
languageDirection = 'ltr'
languageName = 'English'
title = 'Project Documentation'
weight = 2

[versions.'v1.0.0']
[versions.'v2.0.0']
[versions.'v3.0.0']
{{< /code-toggle >}}

This template:

```go-html-template
<ul>
  {{ range .Sites }}
    <li><a href="{{ .Home.RelPermalink }}">{{ .Title }} {{ .Version.Name }}</a></li>
  {{ end }}
</ul>
```

Produces a list of links to each home page:

```html
<ul>
  <li><a href="/v3.0.0/de/">Projekt Dokumentation v3.0.0</a></li>
  <li><a href="/v2.0.0/de/">Projekt Dokumentation v2.0.0</a></li>
  <li><a href="/v1.0.0/de/">Projekt Dokumentation v1.0.0</a></li>
  <li><a href="/v3.0.0/en/">Project Documentation v3.0.0</a></li>
  <li><a href="/v2.0.0/en/">Project Documentation v2.0.0</a></li>
  <li><a href="/v1.0.0/en/">Project Documentation v1.0.0</a></li>
</ul>
```

To render a link to the home page of the [default site](g):

```go-html-template
{{ with .Sites.Default }}
  <a href="{{ .Home.RelPermalink }}">{{ .Title }}</a>
{{ end }}
```

This is equivalent to:

```go-html-template
{{ with index .Sites 0 }}
  <a href="{{ .Home.RelPermalink }}">{{ .Title }}</a>
{{ end }}
```
