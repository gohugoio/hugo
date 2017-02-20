---
title: Shortcode, Git, and Hugo Variables
linktitle: Shortcode, Git, and Hugo Variables
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
tags: [shortcodes,git]
draft: false
weight: 50
aliases: []
toc: false
needsreview: true
notesforauthors:
---

## Shortcodes

`.Parent` (reference to nested shortcodes paragraph in /shortcodes/)
`.IsNamedParams` (reference to shortcodes)
`.Inner`

## Hugo Variables

Also available is `.Hugo` which has the following:

**.Hugo.Generator** Meta tag for the version of Hugo that generated the site. Highly recommended to be included by default in all theme headers so we can start to track the usage and popularity of Hugo. Unlike other variables it outputs a **complete** HTML tag, e.g. `<meta name="generator" content="Hugo 0.15" />`<br>
**.Hugo.Version** The current version of the Hugo binary you are using e.g. `0.13-DEV`<br>
**.Hugo.CommitHash** The git commit hash of the current Hugo binary e.g. `0e8bed9ccffba0df554728b46c5bbf6d78ae5247`<br>
**.Hugo.BuildDate** The compile date of the current Hugo binary formatted with RFC 3339 e.g. `2002-10-02T10:00:00-05:00`<br>

