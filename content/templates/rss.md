---
aliases:
- /layout/rss/
date: 2013-07-01
linktitle: RSS
menu:
  main:
    parent: layout
next: /templates/sitemap
notoc: one
prev: /templates/partials
title: RSS (feed) Templates
weight: 90
---

Like all other templates, you can use a single RSS template to generate
all of your RSS feeds, or you can create a specific template for each
individual feed. Unlike other templates, *Hugo ships with its own
[RSS 2.0 template](#the-embedded-rss-xml:eceb479b7b3b2077408a2878a29e1320).
In most cases this will be sufficient, and an RSS
template will not need to be provided by the user.*

RSS pages are of the type "node" and have all the [node
variables](/layout/variables/) available to use in the templates.


## Which Template will be rendered?
Hugo uses a set of rules to figure out which template to use when
rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present,
then the next one in the list will be used. This enables you to craft
specific layouts when you want to without creating more templates
than necessary. For most sites only the \_default file at the end of
the list will be needed.

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

If the following are provided in the site’s config file, then they
will be included in the RSS output. Example values are provided.

    languageCode = "en-us"
    copyright = "This work is licensed under a Creative Commons Attribution-ShareAlike 4.0 International License."

    [author]
        name = "My Name Here"


## The Embedded rss.xml
This is the RSS template that ships with Hugo. It adheres to the
[RSS 2.0 Specification][RSS 2.0].

    <rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
      <channel>
          <title>{{ .Title }} on {{ .Site.Title }} </title>
          <generator uri="https://gohugo.io">Hugo</generator>
        <link>{{ .Permalink }}</link>
        {{ with .Site.LanguageCode }}<language>{{.}}</language>{{end}}
        {{ with .Site.Author.name }}<author>{{.}}</author>{{end}}
        {{ with .Site.Copyright }}<copyright>{{.}}</copyright>{{end}}
        <updated>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</updated>
        {{ range first 15 .Data.Pages }}
        <item>
          <title>{{ .Title }}</title>
          <link>{{ .Permalink }}</link>
          <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</pubDate>
          {{with .Site.Author.name}}<author>{{.}}</author>{{end}}
          <guid>{{ .Permalink }}</guid>
          <description>{{ .Content | html }}</description>
        </item>
        {{ end }}
      </channel>
    </rss>

*Important: Hugo will automatically add the following header line to this file
on render… please don't include this in the template as it's not valid HTML.*

    <?xml version="1.0" encoding="utf-8" standalone="yes" ?>


[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
