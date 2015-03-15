---
aliases:
- /layout/sitemap/
date: 2014-05-07
linktitle: Sitemap
menu:
  main:
    parent: layout
next: /templates/404
notoc: true
prev: /templates/rss
title: Sitemap Template
weight: 95
---

A single Sitemap template is used to generate the `sitemap.xml` file.
Hugo automatically comes with this template file. **No work is needed on
the users' part unless they want to customize `sitemap.xml`.**

This page is of the type "node" and have all the [node
variables](/layout/variables/) available to use in this template
along with Sitemap-specific ones:

**.Sitemap.ChangeFreq** The page change frequency<br>
**.Sitemap.Priority** The priority of the page<br>

In addition to the standard node variables, the homepage has access to all
site pages through `.Data.Pages`.

If provided, Hugo will use `/layouts/sitemap.xml` instead of the internal
one.

## Hugo’s sitemap.xml

This template respects the version 0.9 of the [Sitemap
Protocol](http://www.sitemaps.org/protocol.html).

    <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
      {{ range .Data.Pages }}
      <url>
        <loc>{{ .Permalink }}</loc>
        <lastmod>{{ safeHtml ( .Date.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ with .Sitemap.ChangeFreq }}
        <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
        <priority>{{ .Sitemap.Priority }}</priority>{{ end }}
      </url>
      {{ end }}
    </urlset>

***Important:** Hugo will automatically add the following header line to this file
on render. Please don't include this in the template as it's not valid HTML.*

    <?xml version="1.0" encoding="utf-8" standalone="yes" ?>
