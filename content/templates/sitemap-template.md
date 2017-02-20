---
title: Sitemap Template
linktitle: Sitemap
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
tags: [sitemap, xml]
categories: [templates]
draft: false
slug:
aliases: [/layout/sitemap/,/templates/sitemap/]
toc: false
needsreview: true
---

A single Sitemap template is used to generate the `sitemap.xml` file.
Hugo automatically comes with this template file. **No work is needed on
the users' part unless they want to customize `sitemap.xml`.**

A sitemap is a `Page` and have all the [page variables](/layout/variables/) available to use in this template along with Sitemap-specific ones:

`.Sitemap.ChangeFreq`
: The page change frequency

`.Sitemap.Priority`
: The priority of the page

`.Sitemap.Filename`
: The sitemap filename

If provided, Hugo will use `/layouts/sitemap.xml` instead of the internal one.

## Hugo’s sitemap.xml

This template respects the version 0.9 of the [Sitemap Protocol](http://www.sitemaps.org/protocol.html).

```xml
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {{ range .Data.Pages }}
  <url>
    <loc>{{ .Permalink }}</loc>{{ if not .Lastmod.IsZero }}
    <lastmod>{{ safeHTML ( .Lastmod.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ end }}{{ with .Sitemap.ChangeFreq }}
    <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
    <priority>{{ .Sitemap.Priority }}</priority>{{ end }}
  </url>
  {{ end }}
</urlset>
```

{{% note %}}
Hugo will automatically add the following header line to this file
on render. Please don't include this in the template as it's not valid HTML.

`<?xml version="1.0" encoding="utf-8" standalone="yes" ?>`
{{% /note %}}

## Configuring `sitemap.xml`

Defaults for `<changefreq>`, `<priority>` and `filename` values can be set in the site's config file, e.g.:

```toml
[sitemap]
  changefreq = "monthly"
  priority = 0.5
  filename = "sitemap.xml"
```

The same fields can be specified in an individual page's front matter in order to override the value for that page.