---
title: urls.AbsLangURL
description: Returns an absolute URL with a language prefix, if any.
categories: []
keywords: []
action:
  aliases: [absLangURL]
  related:
    - functions/urls/AbsURL 
    - functions/urls/RelLangURL
    - functions/urls/RelURL
  returnType: string
  signatures: [urls.AbsLangURL INPUT]
aliases: [/functions/abslangurl]
---

Use this function with both monolingual and multilingual configurations. The URL returned by this function depends on:

- Whether the input begins with a slash
- The `baseURL` in site configuration
- The language prefix, if any

In examples that follow, the project is multilingual with content in both English (`en`) and Spanish (`es`). The returned values are from the English site.

### Input does not begin with a slash

If the input does not begin with a slash, the path in the resulting URL will be relative to the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ absLangURL "" }}          → https://example.org/en/
{{ absLangURL "articles" }}  → https://example.org/en/articles
{{ absLangURL "style.css" }} → https://example.org/en/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absLangURL "" }}          → https://example.org/docs/en/
{{ absLangURL "articles" }}  → https://example.org/docs/en/articles
{{ absLangURL "style.css" }} → https://example.org/docs/en/style.css
```

### Input begins with a slash

If the input begins with a slash, the path in the resulting URL will be relative to the protocol+host of the `baseURL` in your site configuration.

With `baseURL = https://example.org/`

```go-html-template
{{ absLangURL "/" }}          → https://example.org/en/
{{ absLangURL "/articles" }}  → https://example.org/en/articles
{{ absLangURL "/style.css" }} → https://example.org/en/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absLangURL "/" }}          → https://example.org/en/
{{ absLangURL "/articles" }}  → https://example.org/en/articles
{{ absLangURL "/style.css" }} → https://example.org/en/style.css
```
