---
title: Sitemap templates
description: Hugo provides built-in sitemap templates.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 140
weight: 140
toc: true
aliases: [/layout/sitemap/,/templates/sitemap-template/]
---

## Overview

Hugo's embedded sitemap templates conform to v0.9 of the [sitemap protocol].

With a monolingual project, Hugo generates a sitemap.xml file in the root of the [`publishDir`] using the [embedded sitemap template].

With a multilingual project, Hugo generates:

- A sitemap.xml file in the root of each site (language) using the [embedded sitemap template]
- A sitemap.xml file in the root of the [`publishDir`] using the [embedded sitemapindex template]

[embedded sitemap template]: {{% eturl sitemap %}}
[embedded sitemapindex template]: {{% eturl sitemapindex %}}

## Configuration

These are the default sitemap configuration values. They apply to all pages unless overridden in front matter.

{{< code-toggle config=sitemap />}}

changefreq
: (`string`) How frequently a page is likely to change. Valid values are `always`, `hourly`, `daily`, `weekly`, `monthly`, `yearly`, and `never`. With the default value of `""` Hugo will omit this field from the sitemap. See [details](https://www.sitemaps.org/protocol.html#changefreqdef).

disable {{< new-in 0.125.0 >}}
: (`bool`) Whether to disable page inclusion. Default is `false`. Set to `true` in front matter to exclude the page.

filename
: (`string`) The name of the generated file. Default is `sitemap.xml`.

priority
: (`float`) The priority of a page relative to any other page on the site. Valid values range from 0.0 to 1.0.  With the default value of `-1` Hugo will omit this field from the sitemap. See [details](https://www.sitemaps.org/protocol.html#priority).

## Override default values

Override the default values for a given page in front matter.

{{< code-toggle file=news.md fm=true >}}
title = 'News'
[sitemap]
  changefreq = 'weekly'
  disable = true
  priority = 0.8
{{</ code-toggle >}}

## Override built-in templates

To override the built-in sitemap.xml template, create a new file in either of these locations:

- layouts/sitemap.xml
- layouts/_default/sitemap.xml

When ranging through the page collection, access the _change frequency_ and _priority_ with `.Sitemap.ChangeFreq` and `.Sitemap.Priority` respectively.

To override the built-in sitemapindex.xml template, create a new file in either of these locations:

- layouts/sitemapindex.xml
- layouts/_default/sitemapindex.xml

## Disable sitemap generation

You may disable sitemap generation in your site configuration:

{{< code-toggle file=hugo >}}
disableKinds = ['sitemap']
{{</ code-toggle >}}

[`publishDir`]: /getting-started/configuration#publishdir
[sitemap protocol]: <https://www.sitemaps.org/protocol.html>
