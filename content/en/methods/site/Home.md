---
title: Home
description: Returns the home Page object for the given site.
categories: []
keywords: []
action:
  related: []
  returnType: page.Page
  signatures: [SITE.Home]
---

This method is useful for obtaining a link to the home page.

Site configuration:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/docs/'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.Home.Permalink }} → https://example.org/docs/ 
{{ .Site.Home.RelPermalink }} → /docs/
```
