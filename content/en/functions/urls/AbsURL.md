---
title: urls.AbsURL 
description: Returns an absolute URL.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [absURL]
    returnType: string
    signatures: [urls.AbsURL INPUT]
aliases: [/functions/absurl]
---

With multilingual configurations, use the [`urls.AbsLangURL`] function instead. The URL returned by this function depends on:

- Whether the input begins with a slash (`/`)
- The `baseURL` in your site configuration

## Input does not begin with a slash

If the input does not begin with a slash, the path in the resulting URL will be relative to the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ absURL "" }}          → https://example.org/
{{ absURL "articles" }}  → https://example.org/articles
{{ absURL "style.css" }} → https://example.org/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absURL "" }}          → https://example.org/docs/
{{ absURL "articles" }}  → https://example.org/docs/articles
{{ absURL "style.css" }} → https://example.org/docs/style.css
```

## Input begins with a slash

If the input begins with a slash, the path in the resulting URL will be relative to the protocol+host of the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ absURL "/" }}          → https://example.org/
{{ absURL "/articles" }}  → https://example.org/articles
{{ absURL "/style.css" }} → https://example.org/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absURL "/" }}          → https://example.org/
{{ absURL "/articles" }}  → https://example.org/articles
{{ absURL "/style.css" }} → https://example.org/style.css
```

> [!note]
> As illustrated by the previous example, using a leading slash is rarely desirable and can lead to unexpected outcomes. In nearly all cases, omit the leading slash.

[`urls.AbsLangURL`]: /functions/urls/abslangurl/
