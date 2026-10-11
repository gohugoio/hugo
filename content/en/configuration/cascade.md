---
title: Configure cascade
linkTitle: Cascade
description: Configure cascade.
categories: []
keywords: []
---

{{% glossary-term "cascade" %}}

Use the `cascade` configuration to pass values down to pages. Hugo supports two forms: a map form for a single cascade, and an array form for applying different values to different subsets of pages. Both forms support the [`target`](#target) key.

> [!NOTE]
> You can also configure cascading behavior within a page's front matter. See [details][].

## Map form

Define a single cascade map. For example, this configuration cascades the `color` page parameter to all pages:

{{< code-toggle file=hugo >}}
[cascade.params]
color = 'red'
{{< /code-toggle >}}

## Array form

Define an array of cascade maps to apply different values to different subsets of pages. For example, this configuration cascades a different `color` page parameter to each of the `articles` and `tutorials` sections and their descendants:

{{< code-toggle file=hugo >}}
[[cascade]]
[cascade.params]
color = 'red'
[cascade.target]
path = '{/articles,/articles/**}'
[[cascade]]
[cascade.params]
color = 'blue'
[cascade.target]
path = '{/tutorials,/tutorials/**}'
{{< /code-toggle >}}

## Target

In both the map form and the array form, the optional `target` key accepts a [page matcher](g) to limit cascaded values to a subset of pages. If you omit `target`, values cascade to all pages.

{{% include "/_common/configuration/page-matcher.md" %}}

For example, this configuration cascades the `color` page parameter to the `articles` section and its descendants, but only for the English (`en`) and German (`de`) language sites:

{{< code-toggle file=hugo >}}
[cascade.params]
color = 'red'
[cascade.target]
path = '{/articles,/articles/**}'
[cascade.target.sites.matrix]
languages = '{en,de}'
{{< /code-toggle >}}

[details]: /content-management/front-matter/#cascade-1
