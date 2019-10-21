---
title: Sitemap Template
# linktitle: Sitemap
description: Hugo ships with a built-in template file observing the v0.9 of the Sitemap Protocol, but you can override this template if needed.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [sitemap, xml, templates]
menu:
  docs:
    parent: "templates"
    weight: 160
weight: 160
sections_weight: 160
draft: false
aliases: [/layout/sitemap/,/templates/sitemap/]
toc: false
---

A single Sitemap template is used to generate the `sitemap.xml` file.
Hugo automatically comes with this template file. *No work is needed on
the users' part unless they want to customize `sitemap.xml`.*

A sitemap is a `Page` and therefore has all the [page variables][pagevars] available to use in this template along with Sitemap-specific ones:

`.Sitemap.ChangeFreq`
: The page change frequency

`.Sitemap.Priority`
: The priority of the page

`.Sitemap.Filename`
: The sitemap filename

If provided, Hugo will use `/layouts/sitemap.xml` instead of the internal `sitemap.xml` template that ships with Hugo.

## Sitemap Templates

Hugo has built-on Sitemap templates, but you can provide your own if needed, in either `layouts/sitemap.xml` or `layouts/_default/sitemap.xml`.

For multilingual sites, we also create a Sitemap index. You can provide a custom layout for that in either `layouts/sitemapindex.xml` or `layouts/_default/sitemapindex.xml`.

## Hugo’s sitemap.xml

This template respects the version 0.9 of the [Sitemap Protocol](http://www.sitemaps.org/protocol.html).

```xml
{{ printf "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>" | safeHTML }}
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xhtml="http://www.w3.org/1999/xhtml">
  {{ range .Data.Pages }}
  <url>
    <loc>{{ .Permalink }}</loc>{{ if not .Lastmod.IsZero }}
    <lastmod>{{ safeHTML ( .Lastmod.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ end }}{{ with .Sitemap.ChangeFreq }}
    <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
    <priority>{{ .Sitemap.Priority }}</priority>{{ end }}{{ if .IsTranslated }}{{ range .Translations }}
    <xhtml:link
                rel="alternate"
                hreflang="{{ .Lang }}"
                href="{{ .Permalink }}"
                />{{ end }}
    <xhtml:link
                rel="alternate"
                hreflang="{{ .Lang }}"
                href="{{ .Permalink }}"
                />{{ end }}
  </url>
  {{ end }}
</urlset>
```

## Hugo's sitemapindex.xml

This is used to create a Sitemap index in multilingual mode:

```xml
{{ printf "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>" | safeHTML }}
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	{{ range . }}
	<sitemap>
	   	<loc>{{ .SitemapAbsURL }}</loc>
		{{ if not .LastChange.IsZero }}
	   	<lastmod>{{ .LastChange.Format "2006-01-02T15:04:05-07:00" | safeHTML }}</lastmod>
		{{ end }}
	</sitemap>
	{{ end }}
</sitemapindex>
```

## Configure `sitemap.xml`

Defaults for `<changefreq>`, `<priority>` and `filename` values can be set in the site's config file, e.g.:

{{< code-toggle file="config" >}}
[sitemap]
  changefreq = "monthly"
  priority = 0.5
  filename = "sitemap.xml"
{{</ code-toggle >}}

The same fields can be specified in an individual content file's front matter in order to override the value assigned to that piece of content at render time.



[pagevars]: /variables/page/
