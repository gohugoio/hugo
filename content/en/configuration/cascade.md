---
title: Configure cascade
linkTitle: Cascade
description: Configure cascade.
categories: []
keywords: []
---

You can configure your site to cascade front matter values to the home page and any of its descendants. However, this cascading will be prevented if the descendant already defines the field, or if a closer ancestor [node](g) has already cascaded a value for the same field through its front matter's `cascade` key.

> [!note]
> You can also configure cascading behavior within a page's front matter. See&nbsp;[details].

For example, to cascade a "color" parameter to the home page and all its descendants:

{{< code-toggle file=hugo >}}
title = 'Home'
[cascade.params]
color = 'red'
{{< /code-toggle >}}

## Target

<!-- TODO
Update the <version> and <date> below when we actually get around to deprecating _target.

We deprecated the `_target` front matter key in favor of `target` in <version> on <date>. Remove footnote #1 on or after 2026-03-10 (15 months after deprecation).
-->

The `target`[^1] keyword allows you to target specific pages or [environments](g). For example, to cascade a "color" parameter to pages within the "articles" section, including the "articles" section page itself:

[^1]: The `_target` alias for `target` is deprecated and will be removed in a future release.

{{< code-toggle file=hugo >}}
[cascade.params]
color = 'red'
[cascade.target]
path = '{/articles,/articles/**}'
{{< /code-toggle >}}

Use any combination of these keywords to target pages and/or environments:

environment
: (`string`) A [glob](g) pattern matching the build [environment](g). For example: `{staging,production}`.

kind
: (`string`) A [glob](g) pattern matching the [page kind](g). For example: ` {taxonomy,term}`.

lang
: (`string`) A [glob](g) pattern matching the [page language]. For example: `{en,de}`.

path
: (`string`) A [glob](g) pattern matching the page's [logical path](g). For example: `{/books,/books/**}`.

## Array

Define an array of cascade parameters to apply different values to different targets. For example:

{{< code-toggle file=hugo >}}
[[cascade]]
[cascade.params]
color = 'red'
[cascade.target]
path = '{/books/**}'
kind = 'page'
lang = '{en,de}'
[[cascade]]
[cascade.params]
color = 'blue'
[cascade.target]
path = '{/films/**}'
kind = 'page'
environment = 'production'
{{< /code-toggle >}}

[details]: /content-management/front-matter/#cascade-1
[page language]: /methods/page/language/
