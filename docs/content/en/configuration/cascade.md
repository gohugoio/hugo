---
title: Configure cascade
linkTitle: Cascade
description: Configure cascade.
categories: []
keywords: []
---

You can configure your site to cascade front matter values to the home page and any of its descendants. However, this cascading will be prevented if the descendant already defines the field, or if a closer ancestor [branch](g) has already cascaded a value for the same field through its front matter's `cascade` key.

> [!NOTE]
> You can also configure cascading behavior within a page's front matter. See [details][].

For example, to cascade the `color` page parameter to all pages:

{{< code-toggle file=hugo >}}
[cascade.params]
color = 'red'
{{< /code-toggle >}}

## Target

<!-- TODO
We deprecated the `_target` front matter key in favor of `target` in v0.156.0 on 2026-02-17. Remove footnote #1 somewhere after v0.171.0, 15 minor releases
after deprecation.
-->

The `target` key accepts a [page matcher](g) to limit cascaded values to a subset of pages.[^1] If a target is omitted, values cascade to all pages.

{{% include "/_common/configuration/page-matcher.md" %}}

For example, to cascade the `color` page parameter to the `articles` section and its descendants, but only for the English (`en`) and German (`de`) language sites:

{{< code-toggle file=hugo >}}
[cascade.params]
color = 'red'
[cascade.target]
path = '{/articles,/articles/**}'
[cascade.target.sites.matrix]
languages = '{en,de}'
{{< /code-toggle >}}

## Array

Define an array of cascade maps to apply different values to different targets. For example:

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

[^1]: The `_target` alias for `target` is deprecated and will be removed in a future release.

[details]: /content-management/front-matter/#cascade-1
