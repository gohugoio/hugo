---
title: urls.RelURL
description: Returns a relative URL.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [relURL]
    returnType: string
    signatures: [urls.RelURL INPUT]
aliases: [/functions/relurl]
---

With multilingual configurations, use the [`urls.RelLangURL`] function instead. The URL returned by this function depends on:

- Whether the input begins with a slash (`/`)
- The `baseURL` in your site configuration

## Input does not begin with a slash

If the input does not begin with a slash, the resulting URL will be relative to the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ relURL "" }}                         → /
{{ relURL "articles" }}                 → /articles
{{ relURL "style.css" }}                → /style.css
{{ relURL "https://example.org" }}      → https://example.org
{{ relURL "https://example.org/" }}     → /
{{ relURL "https://www.example.org" }}  → https://www.example.org
{{ relURL "https://www.example.org/" }} → https://www.example.org/
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ relURL "" }}                           → /docs/
{{ relURL "articles" }}                   → /docs/articles
{{ relURL "style.css" }}                  → /docs/style.css
{{ relURL "https://example.org" }}        → https://example.org
{{ relURL "https://example.org/" }}       → https://example.org/
{{ relURL "https://example.org/docs" }}   → https://example.org/docs
{{ relURL "https://example.org/docs/" }}  → /docs
{{ relURL "https://www.example.org" }}    → https://www.example.org
{{ relURL "https://www.example.org/" }}   → https://www.example.org/
```

## Input begins with a slash

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

> [!note]
> As illustrated by the previous example, using a leading slash is rarely desirable and can lead to unexpected outcomes. In nearly all cases, omit the leading slash.

[`urls.RelLangURL`]: /functions/urls/rellangurl/
