---
title: hugo.IsMultilingual
description: Reports whether there are two or more configured languages.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [] 
    returnType: bool
    signatures: [hugo.IsMultilingual]
---

Project configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true
[languages]
  [languages.de]
    label = 'Deutsch'
    locale = 'de-DE'
    title = 'Projekt Dokumentation'
    weight = 1
  [languages.en]
    label = 'English'
    locale = 'en-US'
    title = 'Project Documentation'
    weight = 2
{{< /code-toggle >}}

Template:

```go-html-template
{{ hugo.IsMultilingual }} → true
```
