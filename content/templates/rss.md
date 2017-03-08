---
title: RSS Templates
linktitle: RSS Templates
description: Hugo ships with its own RSS 2.0 template that requires almost no configuration, or you can create your own RSS templates.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [rss, xml]
categories: [templates]
weight: 150
draft: false
aliases: [/templates/rss/]
toc: true
wip: true
---

## RSS Template Lookup Order

Like all other templates, you can use a single RSS template to generate all of your RSS feeds, or you can create a specific template for each individual feed.

1. `/layouts/section/<section>.rss.xml`
2. `/layouts/\_default/rss.xml`
3. `/themes/<theme>/layouts/section/<section>.rss.xml`
4. `/themes/<theme>/layouts/\_default/rss.xml`

{{% note "Hugo Ships with an RSS Template" %}}
Unlike other Hugo templates, Hugo ships with its own [RSS 2.0 template](#the-embedded-rss-xml). The embedded template will be sufficient in most cases, and an RSS template will not need to be provided by the user. But you can provide an RSS template, as you can see in the next section.
{{% /note %}}

RSS pages are of the type `Page` and have all the [page variables](/layout/variables/) available to use in the templates.

### Section RSS

A [section’s][section] RSS will be rendered at `/<SECTION>/index.xml` (e.g., http://spf13.com/project/index.xml).

Hugo provides the ability for you to define any RSS type you wish and can have different RSS files for each section and taxonomy.

## Which Template will be Rendered?

Hugo uses a set of rules to figure out which template to use when rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present, then the next one in the list will be used. This enables you to craft specific layouts when you want to without creating more templates than necessary. For most sites only the `\_default` file at the end of the list will be needed.

### Main RSS

1. `/layouts/rss.xml`
2. `/layouts/\_default/rss.xml`
3.  [Embedded rss.xml][embedded]

### Section RSS

1. `/layouts/section/<SECTION>.rss.xml`
2. `/layouts/\_default/rss.xml`
3. `/themes/<THEME>/layouts/section/<SECTION>.rss.xml`
4. `/themes/<THEME>/layouts/\_default/rss.xml`
5. [Embedded rss.xml][embedded]

### Taxonomy RSS

1. `/layouts/taxonomy/<SINGULAR>.rss.xml`
2. `/layouts/\_default/rss.xml`
3. `/themes/<THEME>/layouts/taxonomy/<SINGULAR>.rss.xml`
4. `/themes/<THEME>/layouts/\_default/rss.xml`
5. [Embedded rss.xml][embedded]

## Configuring RSS

The following values will be included in the RSS output if specified in your site’s [`config` file][config]. Example values are provided.

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

...with the autodiscovery link specified by the line with `rel="alternate"`.

The `.RSSLink` will render the appropriate RSS feed URL for the section, whether it's everything, posts in a section, or a taxonomy.

{{% note %}}
If you reference your RSS link, be sure to specify the MIME type with `type="application/rss+xml"`.
{{% /note %}}

```html
<a href="{{ .URL }}" type="application/rss+xml" target="_blank">{{ .SomeText }}</a>
```

[config]: /getting-started/configuration/
[embedded]: #the-embedded-rss-xml
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[section]: /content-management/sections/