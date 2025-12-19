---
title: Sitemap templates
description: Hugo provides built-in sitemap templates.
categories: []
keywords: []
weight: 130
aliases: [/layout/sitemap/,/templates/sitemap-template/]
---

## Overview

Hugo's embedded sitemap templates conform to v0.9 of the [sitemap protocol].

With a monolingual project, Hugo generates a sitemap.xml file in the root of the [`publishDir`] using the [embedded sitemap template].

With a multilingual project, Hugo generates:

- A sitemap.xml file in the root of each site (language) using the [embedded sitemap template]
- A sitemap.xml file in the root of the [`publishDir`] using the [embedded sitemapindex template]

## Configuration

See [configure sitemap](/configuration/sitemap).

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

To override the built-in sitemap.xml template, create a new `layouts/sitemap.xml` file. When ranging through the page collection, access the _change frequency_ and _priority_ with `.Sitemap.ChangeFreq` and `.Sitemap.Priority` respectively.

To override the built-in sitemapindex.xml template, create a new `layouts/sitemapindex.xml` file.

## Disable sitemap generation

You may disable sitemap generation in your site configuration:

{{< code-toggle file=hugo >}}
disableKinds = ['sitemap']
{{</ code-toggle >}}

[`publishDir`]: /configuration/all/#publishdir
[embedded sitemap template]: <{{% eturl sitemap %}}>
[embedded sitemapindex template]: <{{% eturl sitemapindex %}}>
[sitemap protocol]: https://www.sitemaps.org/protocol.html
