---
title: Sitemap templates
description: Hugo provides built-in sitemap templates.
categories: [templates]
keywords: [sitemap,xml,templates]
menu:
  docs:
    parent: templates
    weight: 170
weight: 170
toc: true
aliases: [/layout/sitemap/,/templates/sitemap/]
---

## Overview

Hugo's built-in sitemap templates conform to v0.9 of the [sitemap protocol].

With a monolingual project, Hugo generates a sitemap.xml file in the root of the [`publishDir`] using the built-in [sitemap.xml] template.

With a multilingual project, Hugo generates:

- A sitemap.xml file in the root of each site (language) using the built-in [sitemap.xml] template
- A sitemap.xml file in the root of the [`publishDir`] using the built-in [sitemapindex.xml] template

## Configuration

Set the default values for [change frequency] and [priority], and the name of the generated file, in your site configuration.

{{< code-toggle config=sitemap />}}

changefreq
: How frequently a page is likely to change. Valid values are `always`, `hourly`, `daily`, `weekly`, `monthly`, `yearly`, and `never`. Default is `""` (change frequency omitted from rendered sitemap).

filename
: The name of the generated file. Default is `sitemap.xml`.

priority
: The priority of a page relative to any other page on the site. Valid values range from 0.0 to 1.0. Default is `-1` (priority omitted from rendered sitemap).

## Override default values

Override the default values for a given page in front matter.

{{< code-toggle file=news.md fm=true >}}
title = 'News'
[sitemap]
  changefreq = 'weekly'
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
[change frequency]: <https://www.sitemaps.org/protocol.html#changefreqdef>
[priority]: <https://www.sitemaps.org/protocol.html#priority>
[sitemap protocol]: <https://www.sitemaps.org/protocol.html>
[sitemap.xml]: <https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_default/sitemap.xml>
[sitemapindex.xml]: <https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_default/sitemapindex.xml>
