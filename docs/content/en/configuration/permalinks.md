---
title: Configure permalinks
linkTitle: Permalinks
description: Configure permalinks.
categories: []
keywords: []
---

Use the `permalinks` configuration to define custom URL patterns for your pages. Hugo supports two forms: a map form for simple section-based patterns, and an array form that supports [page matchers](g) for more precise targeting.

> [!NOTE]
> The [`url`][] front matter field overrides any matching permalink pattern.

## Map form

Define URL patterns for each top-level [section](g), keyed by [page kind](g). For example, to configure URL patterns for the `articles` section:

{{< code-toggle file=hugo >}}
[permalinks.page]
articles = '/blog/:year/:month/:slug/'
[permalinks.section]
articles = '/blog/'
{{< /code-toggle >}}

To configure permalinks per language, nest the `permalinks` key under the language key:

{{< code-toggle file=hugo >}}
[languages]
  [languages.de]
    label = 'Deutsch'
    locale = 'de-DE'
    weight = 1
    [languages.de.permalinks]
      [languages.de.permalinks.page]
        articles = '/artikel/:year/:month/:slug/'
      [languages.de.permalinks.section]
        articles = '/artikel/'
  [languages.en]
    label = 'English'
    locale = 'en-US'
    weight = 2
    [languages.en.permalinks]
      [languages.en.permalinks.page]
        articles = '/blog/:year/:month/:slug/'
      [languages.en.permalinks.section]
        articles = '/blog/'
{{< /code-toggle >}}

## Array form

{{< new-in 0.161.0 />}}

Define an array of permalink entries to apply different URL patterns to different subsets of pages. Each entry requires a `pattern` key. Hugo applies the first matching pattern.

The optional `target` key accepts a [page matcher](g). If `target` is omitted, the pattern applies to all pages.

{{% include "/_common/configuration/page-matcher.md" %}}

For example, to apply language-specific URL patterns to the `articles` section page and its leaf pages separately:

{{< code-toggle file=hugo >}}
[[permalinks]]
  pattern = '/artikel/'
  [permalinks.target]
    path = '{/articles}'
    [permalinks.target.sites]
      [permalinks.target.sites.matrix]
        languages = ['de']
[[permalinks]]
  pattern = '/artikel/:year/:month/:slug/'
  [permalinks.target]
    path = '{/articles/**}'
    [permalinks.target.sites]
      [permalinks.target.sites.matrix]
        languages = ['de']
[[permalinks]]
  pattern = '/blog/'
  [permalinks.target]
    path = '{/articles}'
    [permalinks.target.sites]
      [permalinks.target.sites.matrix]
        languages = ['en']
[[permalinks]]
  pattern = '/blog/:year/:month/:slug/'
  [permalinks.target]
    path = '{/articles/**}'
    [permalinks.target.sites]
      [permalinks.target.sites.matrix]
        languages = ['en']
{{< /code-toggle >}}

To define a fallback that matches any page not already matched by a preceding entry, place a pattern without a `target` key at the end:

{{< code-toggle file=hugo >}}
[[permalinks]]
pattern = '/:section/:slug/'
{{< /code-toggle >}}

## Tokens

Use these tokens when defining a URL pattern.

{{% include "/_common/permalink-tokens.md" %}}

[`url`]: /content-management/front-matter/#url
