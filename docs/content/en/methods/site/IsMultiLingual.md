---
title: IsMultiLingual
description: Reports whether there are two or more configured languages.
categories: []
keywords: []
action:
  related: []
  returnType: bool
  signatures: [SITE.IsMultiLingual]
---

{{% deprecated-in 0.124.0 %}}
Use [`hugo.IsMultilingual`] instead.

[`hugo.IsMultilingual`]: /functions/hugo/ismultilingual/
{{% /deprecated-in %}}

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
