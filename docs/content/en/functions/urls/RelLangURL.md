---
title: urls.RelLangURL
description: Returns a relative URL with a language prefix, if any.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [relLangURL]
    returnType: string
    signatures: [urls.RelLangURL INPUT]
aliases: [/functions/rellangurl]
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

If the input does not begin with a slash, the resulting URL will be relative to the `baseURL` in your site configuration.

When rendering the `en` site with `baseURL = https://example.org/`

```go-html-template
{{ relLangURL "" }}                         → /en/
{{ relLangURL "articles" }}                 → /en/articles
{{ relLangURL "style.css" }}                → /en/style.css
{{ relLangURL "https://example.org" }}      → https://example.org
{{ relLangURL "https://example.org/" }}     → /en
{{ relLangURL "https://www.example.org" }}  → https://www.example.org
{{ relLangURL "https://www.example.org/" }} → https://www.example.org/
```

When rendering the `en` site with `baseURL = https://example.org/docs/`

```go-html-template
{{ relLangURL "" }}                           → /docs/en/
{{ relLangURL "articles" }}                   → /docs/en/articles
{{ relLangURL "style.css" }}                  → /docs/en/style.css
{{ relLangURL "https://example.org" }}        → https://example.org
{{ relLangURL "https://example.org/" }}       → https://example.org/
{{ relLangURL "https://example.org/docs" }}   → https://example.org/docs
{{ relLangURL "https://example.org/docs/" }}  → /docs/en
{{ relLangURL "https://www.example.org" }}    → https://www.example.org
{{ relLangURL "https://www.example.org/" }}   → https://www.example.org/
```

## Input begins with a slash

If the input begins with a slash, the resulting URL will be relative to the protocol+host of the `baseURL` in your site configuration.

When rendering the `en` site with `baseURL = https://example.org/`

```go-html-template
{{ relLangURL "/" }}          → /en/
{{ relLangURL "/articles" }}  → /en/articles
{{ relLangURL "/style.css" }} → /en/style.css
```

When rendering the `en` site with `baseURL = https://example.org/docs/`

```go-html-template
{{ relLangURL "/" }}          → /en/
{{ relLangURL "/articles" }}  → /en/articles
{{ relLangURL "/style.css" }} → /en/style.css
```

> [!note]
> As illustrated by the previous example, using a leading slash is rarely desirable and can lead to unexpected outcomes. In nearly all cases, omit the leading slash.
