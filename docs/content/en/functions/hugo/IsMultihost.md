---
title: hugo.IsMultihost
description: Reports whether each configured language has a unique base URL.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [hugo.IsMultihost]
---

Project configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true
[languages]
  [languages.de]
    baseURL = 'https://de.example.org/'
    label = 'Deutsch'
    locale = 'de-DE'
    title = 'Projekt Dokumentation'
    weight = 1
  [languages.en]
    baseURL = 'https://en.example.org/'
    label = 'English'
    locale = 'en-US'
    title = 'Project Documentation'
    weight = 2
{{< /code-toggle >}}

Template:

```go-html-template
{{ hugo.IsMultihost }} → true
```
