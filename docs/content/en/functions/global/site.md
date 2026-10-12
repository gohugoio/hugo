---
title: site
description: Returns the Site object of the current site, accessible from any context.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: page.siteWrapper
    signatures: [site]
---

Use the `site` function to return the `Site` object regardless of current context.

```go-html-template
{{ site.Params.foo }}
```

When the `Site` object is in context you can use the `Site` property:

```go-html-template
<!-- current context -->
{{ .Site.Params.foo }}
<!-- template context -->
{{ $.Site.Params.foo }}
```

> [!NOTE]
> To simplify your templates, use the global `site` function regardless of whether the `Site` object is in context.
