---
title: "RSS (feed) Templates"
date: "2013-07-01"
weight: 40
notoc: "one"
menu:
  main:
    parent: 'layout'
---

A single RSS template is used to generate all of the RSS content for the entire
site.

RSS pages are of the type "node" and have all the [node
variables](/layout/variables/) available to use in the templates.

In addition to the standard node variables, the homepage has access to
all site content accessible from .Data.Pages

    ▾ layouts/
        rss.xml

## rss.xml
This rss template is used for [spf13.com](http://spf13.com). It adheres to the
ATOM 2.0 Spec.

    <rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
      <channel>
          <title>{{ .Title }} on {{ .Site.Title }} </title>
        <link>{{ .Permalink }}</link>
        <language>en-us</language>
        <author>Steve Francia</author>
        <rights>Copyright (c) 2008 - 2013, Steve Francia; all rights reserved.</rights>
        <updated>{{ .Date }}</updated>
        {{ range .Data.Pages }}
        <item>
          <title>{{ .Title }}</title>
          <link>{{ .Permalink }}</link>
          <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</pubDate>
          <author>Steve Francia</author>
          <guid>{{ .Permalink }}</guid>
          <description>{{ .Content | html }}</description>
        </item>
        {{ end }}
      </channel>
    </rss>

*Important: Hugo will automatically add the following header line to this file
on render...please don't include this in the template as it's not valid HTML.*

    <?xml version="1.0" encoding="utf-8" standalone="yes" ?>
