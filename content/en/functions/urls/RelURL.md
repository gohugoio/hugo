---
title: urls.RelURL
description: Returns a relative URL.
categories: []
keywords: []
action:
  aliases: [relURL]
  related:
    - functions/urls/AbsLangURL
    - functions/urls/AbsURL 
    - functions/urls/RelLangURL
  returnType: string
  signatures: [urls.RelURL INPUT]
aliases: [/functions/relurl]
---

With multilingual configurations, use the [`urls.RelLangURL`] function instead. The URL returned by this function depends on:

- Whether the input begins with a slash
- The `baseURL` in your site configuration

### Input does not begin with a slash

If the input does not begin with a slash, the resulting URL will be relative to the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ relURL "" }}                        → /
{{ relURL "articles" }}                → /articles
{{ relURL "style.css" }}               → /style.css
{{ relURL "https://example.org/foo" }} → /foo
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ relURL "" }}                             → /docs/
{{ relURL "articles" }}                     → /docs/articles
{{ relURL "style.css" }}                    → /docs/style.css
{{ relURL "https://example.org/docs/foo" }} → /docs/foo
```

#### Input begins with a slash

If the input begins with a slash, the resulting URL will be relative to the protocol+host of the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ relURL "/" }}          → /
{{ relURL "/articles" }}  → /articles
{{ relURL "/style.css" }} → /style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ relURL "/" }}          → /
{{ relURL "/articles" }}  → /articles
{{ relURL "/style.css" }} → /style.css
```

[`urls.RelLangURL`]: /functions/urls/rellangurl/
