---
title: IsMultiLingual
description: Reports whether the site is multilingual.
categories: []
keywords: []
action:
  related: []
  returnType: bool
  signatures: [SITE.IsMultiLingual]
---

Site configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true
[languages]
  [languages.de]
    languageCode = 'de-DE'
    languageName = 'Deutsch'
    title = 'Projekt Dokumentation'
    weight = 1
  [languages.en]
    languageCode = 'en-US'
    languageName = 'English'
    title = 'Project Documentation'
    weight = 2
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.IsMultiLingual }} → true
```
