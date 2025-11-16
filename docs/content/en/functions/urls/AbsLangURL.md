---
title: urls.AbsLangURL
description: Returns an absolute URL with a language prefix, if any.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [absLangURL]
    returnType: string
    signatures: [urls.AbsLangURL INPUT]
aliases: [/functions/abslangurl]
---

Use this function with both monolingual and multilingual configurations. The URL returned by this function depends on:

- Whether the input begins with a slash (`/`)
- The `baseURL` in your site configuration
- The language prefix, if any

This is the site configuration for the examples that follow:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[languages.en]
weight = 1
[languages.es]
weight = 2
{{< /code-toggle >}}

## Input does not begin with a slash

If the input does not begin with a slash, the path in the resulting URL will be relative to the `baseURL` in your site configuration.

When rendering the `en` site with `baseURL = https://example.org/`

```go-html-template
{{ absLangURL "" }}           → https://example.org/en/
{{ absLangURL "articles" }}   → https://example.org/en/articles
{{ absLangURL "style.css" }}  → https://example.org/en/style.css
```

When rendering the `en` site with `baseURL = https://example.org/docs/`

```go-html-template
{{ absLangURL "" }}           → https://example.org/docs/en/
{{ absLangURL "articles" }}   → https://example.org/docs/en/articles
{{ absLangURL "style.css" }}  → https://example.org/docs/en/style.css
```

## Input begins with a slash

If the input begins with a slash, the path in the resulting URL will be relative to the protocol+host of the `baseURL` in your site configuration.

When rendering the `en` site with `baseURL = https://example.org/`

```go-html-template
{{ absLangURL "/" }}          → https://example.org/en/
{{ absLangURL "/articles" }}  → https://example.org/en/articles
{{ absLangURL "/style.css" }} → https://example.org/en/style.css
```

When rendering the `en` site with `baseURL = https://example.org/docs/`

```go-html-template
{{ absLangURL "/" }}          → https://example.org/en/
{{ absLangURL "/articles" }}  → https://example.org/en/articles
{{ absLangURL "/style.css" }} → https://example.org/en/style.css
```

> [!note]
> As illustrated by the previous example, using a leading slash is rarely desirable and can lead to unexpected outcomes. In nearly all cases, omit the leading slash.
