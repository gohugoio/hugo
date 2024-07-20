---
title: LanguagePrefix
description: Returns the URL language prefix, if any, for the given site.
categories: []
keywords: []
action:
  related:
    - functions/urls/AbsLangURL
    - functions/urls/RelLangURL
  returnType: string
  signatures: [SITE.LanguagePrefix]
---

Consider this site configuration:

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

When visiting the German language site:

```go-html-template
{{ .Site.LanguagePrefix }} → ""
```

When visiting the English language site:

```go-html-template
{{ .Site.LanguagePrefix }} → /en
```

If you change `defaultContentLanguageInSubdir` to `true`, when visiting the German language site:

```go-html-template
{{ .Site.LanguagePrefix }} → /de
```

You may use the `LanguagePrefix` method with both monolingual and multilingual sites.
