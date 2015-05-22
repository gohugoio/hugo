---
aliases:
- /layout/rss/
date: 2015-05-19
linktitle: RSS
menu:
  main:
    parent: layout
next: /templates/sitemap
notoc: one
prev: /templates/partials
title: RSS (feed) Templates
weight: 90
toc: true
---

Like all other templates, you can use a single RSS template to generate all of your RSS feeds, or you can create a specific template for each individual feed. 

*Unlike other Hugo templates*, Hugo ships with its own [RSS 2.0 template](#the-embedded-rss-xml:eceb479b7b3b2077408a2878a29e1320). In most cases this will be sufficient, and an RSS template will not need to be provided by the user. But you can provide an rss template if you like, as you can see in the next section. 

RSS pages are of the **type "node"** and have all the [node variables](/layout/variables/) available to use in the templates.

## Which Template will be rendered?
Hugo uses a set of rules to figure out which template to use when rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present, then the next one in the list will be used. This enables you to craft specific layouts when you want to without creating more templates than necessary. For most sites only the `\_default` file at the end of the list will be needed.

### Main RSS

* /layouts/rss.xml
* /layouts/\_default/rss.xml
* \__internal/rss.xml

### Section RSS

* /layouts/section/`SECTION`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/section/`SECTION`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml
* \__internal/rss.xml

### Taxonomy RSS

* /layouts/taxonomy/`SINGULAR`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml
* \__internal/rss.xml


## Configuring RSS

If the following values are specified in the site’s config file (`config.toml`), then they will be included in the RSS output. Example values are provided.

    languageCode = "en-us"
    copyright = "This work is licensed under a Creative Commons Attribution-ShareAlike 4.0 International License."

    [author]
        name = "My Name Here"


## The Embedded rss.xml
This is the default RSS template that ships with Hugo. It adheres to the [RSS 2.0 Specification][RSS 2.0].

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

**Important**: _Hugo will automatically add the following header line to this file on render… please don't include this in the template as it's not valid HTML._

~~~css
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
~~~

## Referencing your RSS Feed in `<head>`

In your `header.html` template, you can specify your RSS feed in your `<head></head>` tag like this: 

~~~html
{{ if .RSSlink }}
  <link href="{{ .RSSlink }}" rel="alternate" type="application/rss+xml" title="{{ .Site.Title }}" />
  <link href="{{ .RSSlink }}" rel="feed" type="application/rss+xml" title="{{ .Site.Title }}" />
{{ end }}
~~~

... with the autodiscovery link specified by the line with `rel="alternate"`. 

The `.RSSlink` will render the appropriate RSS feed URL for the section, whether it's everything, posts in a section, or a taxonomy. 

**N.b.**, if you reference your RSS link, be sure to specify the mime type with `type="application/rss+xml"`. 

~~~html
<a href="{{ .URL }}" type="application/rss+xml" target="_blank">{{ .SomeText }}</a>
~~~

[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
