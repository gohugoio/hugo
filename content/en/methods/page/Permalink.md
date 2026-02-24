---
title: Permalink
description: Returns the permalink of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [PAGE.Permalink]
---

Project configuration:

{{< code-toggle file=hugo >}}
title = 'Documentation'
baseURL = 'https://example.org/docs/'
{{< /code-toggle >}}

Template:

```go-html-template
{{ $page := .Site.GetPage "/about" }}
{{ $page.Permalink }} → https://example.org/docs/about/
```
