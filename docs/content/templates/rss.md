---
title: RSS Templates
linktitle: RSS Templates
description: Hugo ships with its own RSS 2.0 template that requires almost no configuration, or you can create your own RSS templates.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [rss, xml]
categories: [templates]
menu:
  docs:
    parent: "templates"
    weight: 150
weight: 150
sections_weight: 150
draft: false
aliases: [/templates/rss/]
toc: true
---

## RSS Template Lookup Order

You can use a single RSS template to generate all of your RSS feeds or create a specific template for each individual feed.

1. `/layouts/section/<section>.rss.xml`
2. `/layouts/_default/rss.xml`
3. `/themes/<theme>/layouts/section/<section>.rss.xml`
4. `/themes/<theme>/layouts/_default/rss.xml`

{{% note "Hugo Ships with an RSS Template" %}}
Hugo ships with its own [RSS 2.0 template](#the-embedded-rss-xml). The embedded template will be sufficient for most use cases.
{{% /note %}}

RSS pages are of the type `Page` and have all the [page variables](/variables/page/) available to use in the templates.

### Section RSS

A [section’s][section] RSS will be rendered at `/<SECTION>/index.xml` (e.g., http://spf13.com/project/index.xml).

Hugo provides the ability for you to define any RSS type you wish and can have different RSS files for each section and taxonomy.

## Lookup Order for RSS Templates

### Main RSS

1. `/layouts/rss.xml`
2. `/layouts/_default/rss.xml`
3.  Embedded rss.xml

### Section RSS

1. `/layouts/section/<SECTION>.rss.xml`
2. `/layouts/_default/rss.xml`
3. `/themes/<THEME>/layouts/section/<SECTION>.rss.xml`
4. `/themes/<THEME>/layouts/_default/rss.xml`
5. Embedded rss.xml

### Taxonomy RSS

1. `/layouts/taxonomy/<SINGULAR>.rss.xml`
2. `/layouts/_default/rss.xml`
3. `/themes/<THEME>/layouts/taxonomy/<SINGULAR>.rss.xml`
4. `/themes/<THEME>/layouts/_default/rss.xml`
5. Embedded rss.xml

## Configure RSS

By default, Hugo will create an unlimited number of RSS entries. You can limit the number of articles included in the built-in RSS templates by assigning a numeric value to `rssLimit:` field in your project's [`config` file][config].

The following values will also be included in the RSS output if specified in your site’s configuration:

```
languageCode = "en-us"
copyright = "This work is licensed under a Creative Commons Attribution-ShareAlike 4.0 International License."

[author]
    name = "My Name Here"
```

## The Embedded rss.xml

This is the default RSS template that ships with Hugo. It adheres to the [RSS 2.0 Specification][RSS 2.0].

```
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ if eq  .Title  .Site.Title }}{{ .Site.Title }}{{ else }}{{ with .Title }}{{.}} on {{ end }}{{ .Site.Title }}{{ end }}</title>
    <link>{{ .Permalink }}</link>
    <description>Recent content {{ if ne  .Title  .Site.Title }}{{ with .Title }}in {{.}} {{ end }}{{ end }}on {{ .Site.Title }}</description>
    <generator>Hugo -- gohugo.io</generator>{{ with .Site.LanguageCode }}
    <language>{{.}}</language>{{end}}{{ with .Site.Author.email }}
    <managingEditor>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</managingEditor>{{end}}{{ with .Site.Author.email }}
    <webMaster>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</webMaster>{{end}}{{ with .Site.Copyright }}
    <copyright>{{.}}</copyright>{{end}}{{ if not .Date.IsZero }}
    <lastBuildDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</lastBuildDate>{{ end }}
    {{ with .OutputFormats.Get "RSS" }}
        {{ printf "<atom:link href=%q rel=\"self\" type=%q />" .Permalink .MediaType | safeHTML }}
    {{ end }}
    {{ range .Data.Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</pubDate>
      {{ with .Site.Author.email }}<author>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</author>{{end}}
      <guid>{{ .Permalink }}</guid>
      <description>{{ .Summary | html }}</description>
    </item>
    {{ end }}
  </channel>
</rss>
```

{{% warning "XML Header" %}}
Hugo will automatically add the following header line to this file on render. Please do *not* include this in the template as it's not valid HTML.
```
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
```
{{% /warning %}}

## Reference your RSS Feed in `<head>`

In your `header.html` template, you can specify your RSS feed in your `<head></head>` tag like this:

```
{{ if .RSSLink }}
  <link href="{{ .RSSLink }}" rel="alternate" type="application/rss+xml" title="{{ .Site.Title }}" />
  <link href="{{ .RSSLink }}" rel="feed" type="application/rss+xml" title="{{ .Site.Title }}" />
{{ end }}
```

...with the auto-discovery link specified by the line with `rel="alternate"`.

The `.RSSLink` will render the appropriate RSS feed URL for the section, whether it's everything, posts in a section, or a taxonomy.

If you reference your RSS link, be sure to specify the MIME type with `type="application/rss+xml"`.

```
<a href="{{ .URL }}" type="application/rss+xml" target="_blank">{{ .SomeText }}</a>
```

[config]: /getting-started/configuration/
[embedded]: #the-embedded-rss-xml
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[section]: /content-management/sections/
