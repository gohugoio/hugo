---
title: Page
description: Returns the Page object of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Page
    signatures: [PAGE.Page]
---

This is a convenience method, useful within _partial_ templates that are called from both _shortcode_ and other template types.

```go-html-template {file="layouts/_shortcodes/foo.html"}
{{ partial "my-partial.html" . }}
```

When the _shortcode_ template calls the _partial_ template, it passes the current [context](g) (the dot). The context includes identifiers such as `Page`, `Params`, `Inner`, and `Name`.

```go-html-template {file="layouts/page.html"}
{{ partial "my-partial.html" . }}
```

When the _page_ template calls the _partial_ template, it also passes the current context (the dot). But in this case, the dot _is_ the `Page` object.

```go-html-template {file="layouts/_partials/my-partial.html"}
The page title is: {{ .Page.Title }}
```

To handle both scenarios, the _partial_ template must be able to access the `Page` object with `Page.Page`.

> [!note]
> And yes, that means you can do `.Page.Page.Page.Page.Title` too.
>
> But don't.
