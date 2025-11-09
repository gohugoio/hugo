---
title: Home
description: Returns the home Page object for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Page
    signatures: [SITE.Home]
---

The `Home` method on a `Site` object is a convenient way to access the home page, and is functionally equivalent to:

```go-html-template
{{ .Site.GetPage "/" }}
```

Because it returns a `Page` object, you can use any of the available [page methods][] by chaining them. For example:

```go-html-template
{{ .Site.Home.Store.Set "greeting" "Hello" }}
```

This method is commonly used to generate a link to the home page. For example:

Site configuration:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/docs/'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.Home.Permalink }} → https://example.org/docs/ 
{{ .Site.Home.RelPermalink }} → /docs/
```

[page methods]: /methods/page/
