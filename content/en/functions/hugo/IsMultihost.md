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

The `hugo.IsMultihost` function reports whether each configured language has a unique `baseURL`.

{{< code-toggle file=hugo >}}
[languages]
  [languages.en]
    baseURL = 'https://en.example.org/'
    languageName = 'English'
    title = 'In English'
    weight = 2
  [languages.fr]
    baseURL = 'https://fr.example.org'
    languageName = 'Français'
    title = 'En Français'
    weight = 1
{{< /code-toggle >}}

```go-html-template
{{ hugo.IsMultihost }} → true
```
