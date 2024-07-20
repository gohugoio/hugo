---
title: urls.RelLangURL
description: Returns a relative URL with a language prefix, if any.
categories: []
keywords: []
action:
  aliases: [relLangURL]
  related:
    - functions/urls/AbsLangURL
    - functions/urls/AbsURL 
    - functions/urls/RelURL
  returnType: string
  signatures: [urls.RelLangURL INPUT]
aliases: [/functions/rellangurl]
---

Use this function with both monolingual and multilingual configurations. The URL returned by this function depends on:

- Whether the input begins with a slash
- The `baseURL` in your site configuration
- The language prefix, if any

In examples that follow, the project is multilingual with content in both English (`en`) and Spanish (`es`). The returned values are from the English site.

### Input does not begin with a slash

If the input does not begin with a slash, the resulting URL will be relative to the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ relLangURL "" }}                        → /en/
{{ relLangURL "articles" }}                → /en/articles
{{ relLangURL "style.css" }}               → /en/style.css
{{ relLangURL "https://example.org/foo" }} → /en/foo
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ relLangURL "" }}                             → /docs/en/
{{ relLangURL "articles" }}                     → /docs/en/articles
{{ relLangURL "style.css" }}                    → /docs/en/style.css
{{ relLangURL "https://example.org/docs/foo" }} → /docs/en/foo
```

#### Input begins with a slash

If the input begins with a slash, the resulting URL will be relative to the protocol+host of the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ relLangURL "/" }}          → /en/
{{ relLangURL "/articles" }}  → /en/articles
{{ relLangURL "/style.css" }} → /en/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ relLangURL "/" }}          → /en/
{{ relLangURL "/articles" }}  → /en/articles
{{ relLangURL "/style.css" }} → /en/style.css
```
