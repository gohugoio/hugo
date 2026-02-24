---
title: hugo.Sites
description: Returns a collection of all sites for all dimensions.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: page.Sites
    signatures: [hugo.Sites]
---

{{< new-in 0.156.0 />}}

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
  {{ range hugo.Sites }}
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
{{ with hugo.Sites.Default }}
  <a href="{{ .Home.RelPermalink }}">{{ .Title }}</a>
{{ end }}
```

This is equivalent to:

```go-html-template
{{ with index hugo.Sites 0 }}
  <a href="{{ .Home.RelPermalink }}">{{ .Title }}</a>
{{ end }}
```
