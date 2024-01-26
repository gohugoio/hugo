---
title: RelPermalink
description: Returns the relative permalink of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Permalink
  returnType: string
  signatures: [PAGE.RelPermalink]
---

Site configuration:

{{< code-toggle file=hugo >}}
title = 'Documentation'
baseURL = 'https://example.org/docs/'
{{< /code-toggle >}}

Template:

```go-html-template
{{ $page := .Site.GetPage "/about" }}
{{ $page.RelPermalink }} → /docs/about/
```
