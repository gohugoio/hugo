---
title: Sitemap
description: Returns the sitemap settings for the given page as defined in front matter, falling back to the sitemap settings as defined in the site configuration.
categories: []
keywords: []
action:
  related: []
  returnType: config.SitemapConfig
  signatures: [PAGE.Sitemap]
toc: true
---

Access to the `Sitemap` method on a `Page` object is restricted to [sitemap templates].

## Methods

ChangeFreq
: (`string`) How frequently a page is likely to change. Valid values are `always`, `hourly`, `daily`, `weekly`, `monthly`, `yearly`, and `never`. Default is "" (change frequency omitted from rendered sitemap).

```go-html-template
{{ .Sitemap.ChangeFreq }}
```

Priority
: (`float`) The priority of a page relative to any other page on the site. Valid values range from 0.0 to 1.0. Default is -1 (priority omitted from rendered sitemap).

```go-html-template
{{ .Sitemap.Priority }}
```

## Example

With this site configuration:

{{< code-toggle file=hugo >}}
[sitemap]
changeFreq = 'monthly'
{{< /code-toggle >}}

And this content:

{{< code-toggle file=content/news.md fm=true >}}
title = 'News'
[sitemap]
changeFreq = 'hourly'
{{< /code-toggle >}}

And this simplistic sitemap template:

{{< code file=layouts/_default/sitemap.xml >}}
{{ printf "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>" | safeHTML }}
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xhtml="http://www.w3.org/1999/xhtml">
  {{ range .Pages }}
    <url>
      <loc>{{ .Permalink }}</loc>
      {{ if not .Lastmod.IsZero }}
        <lastmod>{{ .Lastmod.Format "2006-01-02T15:04:05-07:00" | safeHTML }}</lastmod>
      {{ end }}
      {{ with .Sitemap.ChangeFreq }}
        <changefreq>{{ . }}</changefreq>
      {{ end }}
    </url>
  {{ end }}
</urlset>
{{< /code >}}

The change frequency will be `hourly` for the news page, and `monthly` for other pages.

[sitemap templates]: /templates/sitemap-template/
