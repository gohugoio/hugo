---
title: Languages
description: Returns a collection of language objects for all sites, ordered by language weight.
categories: []
keywords: []
action:
  related:
    - methods/site/Language
  returnType: langs.Languages
  signatures: [SITE.Languages]
---

The `Languages` method on a `Site` object returns a collection of language objects for all sites, ordered by language weight. Each language object points to its language definition in the site configuration.

To inspect the data structure:

```go-html-template
<pre>{{ debug.Dump .Site.Languages }}</pre>
```

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
  {{ range .Site.Languages }}
    <li>{{ .Title }} ({{ .LanguageName }})</li>
  {{ end }}
</ul>
```

Is rendered to:

```html
<ul>
  <li>Projekt Dokumentation (Deutsch)</li>
  <li>Project Documentation (English)</li>
</ul>
```
