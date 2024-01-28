---
title: Permalink
description: Returns the permalink of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/RelPermalink
  returnType: string
  signatures: [PAGE.Permalink]
---

Site configuration:

{{< code-toggle file=hugo >}}
title = 'Documentation'
baseURL = 'https://example.org/docs/'
{{< /code-toggle >}}

Template:

```go-html-template
{{ $page := .Site.GetPage "/about" }}
{{ $page.Permalink }} → https://example.org/docs/about/
```
