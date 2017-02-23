---
title: RSS Templates
linktitle: RSS Templates
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [rss, xml]
categories: [templates]
weight: 150
draft: false
aliases: [/templates/rss/]
toc: true
needsreview: true
---

## RSS Template Lookup Order

Like all other templates, you can use a single RSS template to generate all of your RSS feeds, or you can create a specific template for each individual feed.

* /layouts/section/`SECTION`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/section/`SECTION`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml

{{% note "Hugo Ships with an RSS Template" %}}
*Unlike other Hugo templates*, Hugo ships with its own [RSS 2.0 template](#the-embedded-rss-xml:eceb479b7b3b2077408a2878a29e1320). In most cases this will be sufficient, and an RSS template will not need to be provided by the user. But you can provide an rss template if you like, as you can see in the next section.
{{% /note %}}

RSS pages are of the type `Page` and have all the [page variables](/layout/variables/) available to use in the templates.

### Section RSS

A [section’s][section] RSS will be rendered at /`SECTION`/index.xml (e.g., http://spf13.com/project/index.xml)

*Hugo ships with its own [RSS 2.0][] template. In most cases this will
be sufficient, and an RSS template will not need to be provided by the
user.*

Hugo provides the ability for you to define any RSS type you wish, and
can have different RSS files for each section and taxonomy.

## Which Template will be Rendered?

Hugo uses a set of rules to figure out which template to use when rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present, then the next one in the list will be used. This enables you to craft specific layouts when you want to without creating more templates than necessary. For most sites only the `\_default` file at the end of the list will be needed.

### Main RSS

* /layouts/rss.xml
* /layouts/\_default/rss.xml
* [Embedded rss.xml](#the-embedded-rss-xml:eceb479b7b3b2077408a2878a29e1320)

### Section RSS

* /layouts/section/`SECTION`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/section/`SECTION`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml
* [Embedded rss.xml](#the-embedded-rss-xml:eceb479b7b3b2077408a2878a29e1320)

### Taxonomy RSS

* /layouts/taxonomy/`SINGULAR`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml
* [Embedded rss.xml](#the-embedded-rss-xml:eceb479b7b3b2077408a2878a29e1320)

## Configuring RSS

If the following values are specified in the site’s config file (`config.toml`), then they will be included in the RSS output. Example values are provided.

```toml
languageCode = "en-us"
copyright = "This work is licensed under a Creative Commons Attribution-ShareAlike 4.0 International License."

[author]
    name = "My Name Here"
```


## The Embedded rss.xml
This is the default RSS template that ships with Hugo. It adheres to the [RSS 2.0 Specification][RSS 2.0].

```xml
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ with .Title }}{{.}} on {{ end }}{{ .Site.Title }}</title>
    <link>{{ .Permalink }}</link>
    <description>Recent content {{ with .Title }}in {{.}} {{ end }}on {{ .Site.Title }}</description>
    <generator>Hugo -- gohugo.io</generator>{{ with .Site.LanguageCode }}
    <language>{{.}}</language>{{end}}{{ with .Site.Author.email }}
    <managingEditor>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</managingEditor>{{end}}{{ with .Site.Author.email }}
    <webMaster>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</webMaster>{{end}}{{ with .Site.Copyright }}
    <copyright>{{.}}</copyright>{{end}}{{ if not .Date.IsZero }}
    <lastBuildDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</lastBuildDate>{{ end }}
    <atom:link href="{{.URL}}" rel="self" type="application/rss+xml" />
    {{ range first 15 .Data.Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</pubDate>
      {{ with .Site.Author.email }}<author>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</author>{{end}}
      <guid>{{ .Permalink }}</guid>
      <description>{{ .Content | html }}</description>
    </item>
    {{ end }}
  </channel>
</rss>
```

{{% warning "XML Header" %}}
Hugo will automatically add the following header line to this file on render…please don't include this in the template as it's not valid HTML.
```xml
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
```
{{% /warning %}}

## Referencing your RSS Feed in `<head>`

In your `header.html` template, you can specify your RSS feed in your `<head></head>` tag like this:

```html
{{ if .RSSLink }}
  <link href="{{ .RSSLink }}" rel="alternate" type="application/rss+xml" title="{{ .Site.Title }}" />
  <link href="{{ .RSSLink }}" rel="feed" type="application/rss+xml" title="{{ .Site.Title }}" />
{{ end }}
```

... with the autodiscovery link specified by the line with `rel="alternate"`.

The `.RSSLink` will render the appropriate RSS feed URL for the section, whether it's everything, posts in a section, or a taxonomy.

**N.b.**, if you reference your RSS link, be sure to specify the mime type with `type="application/rss+xml"`.

```html
<a href="{{ .URL }}" type="application/rss+xml" target="_blank">{{ .SomeText }}</a>
```

[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[section]: /content-management/sections/