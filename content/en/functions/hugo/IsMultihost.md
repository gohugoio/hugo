---
title: hugo.IsMultihost
description: Reports whether each configured language has a unique base URL.
categories: []
keywords: []
action:
  aliases: []
  related:
    - /functions/hugo/IsMultilingual
  returnType: bool
  signatures: [hugo.IsMultihost]
---

{{< new-in v0.124.0 >}}

Site configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true
[languages]
  [languages.de]
    baseURL = 'https://de.example.org/'
    languageCode = 'de-DE'
    languageName = 'Deutsch'
    title = 'Projekt Dokumentation'
    weight = 1
  [languages.en]
    baseURL = 'https://en.example.org/'
    languageCode = 'en-US'
    languageName = 'English'
    title = 'Project Documentation'
    weight = 2
{{< /code-toggle >}}

Template:

```go-html-template
{{ hugo.IsMultihost }} → true
```
