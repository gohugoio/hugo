---
title: "Source Directory Organization"
date: "2013-07-01"
aliases: ["/doc/source-directory/"]
weight: 50
notoc: true
menu:
  main:
    parent: 'getting started'
---

Hugo takes a single directory and uses it as the input for creating a complete website.

Hugo has a very small amount of configuration, while remaining highly customizable.
It accomplishes by assuming that you will only provide templates with the intent of
using them.

An example directory may look like:

    .
    ├── config.yaml
    ├── content
    |   ├── post
    |   |   ├── firstpost.md
    |   |   └── secondpost.md
    |   └── quote
    |   |   ├── first.md
    |   |   └── second.md
    ├── layouts
    |   ├── chrome
    |   |   ├── header.html
    |   |   └── footer.html
    |   ├── indexes
    |   |   ├── category.html
    |   |   ├── post.html
    |   |   ├── quote.html
    |   |   └── tag.html
    |   ├── post
    |   |   ├── li.html
    |   |   ├── single.html
    |   |   └── summary.html
    |   ├── quote
    |   |   ├── li.html
    |   |   ├── single.html
    |   |   └── summary.html
    |   ├── shortcodes
    |   |   ├── img.html
    |   |   ├── vimeo.html
    |   |   └── youtube.html
    |   ├── index.html
    |   ├── rss.xml
    |   └── sitemap.xml
    └── static

This directory structure tells us a lot about this site:

1. the website intends to have two different types of content, posts and quotes.
2. It will also apply two different indexes to that content, categories and tags.
3. It will be displaying content in 3 different views, a list, a summary and a full page view.

Included with the repository is this example site ready to be rendered.
