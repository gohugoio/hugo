---
title: Hugo-specific Variables
linktitle: Hugo Variables
description: The `.Hugo` variable provides easy access to Hugo-related data.
date: 2017-03-12
publishdate: 2017-03-12
lastmod: 2017-03-12
categories: [variables and params]
keywords: [hugo,generator]
draft: false
menu:
  docs:
    parent: "variables"
    weight: 60
weight: 60
sections_weight: 60
aliases: []
toc: false
wip: false
---

{{% warning "Deprecated" %}}
Page's `.Hugo` is deprecated and will be removed in a future release. Use the global `hugo` function.  
For example: `hugo.Generator`.
{{% /warning %}}

It contains the following fields:

.Hugo.Generator
: `<meta>` tag for the version of Hugo that generated the site. `.Hugo.Generator` outputs a *complete* HTML tag; e.g. `<meta name="generator" content="Hugo 0.18" />`

.Hugo.Version
: the current version of the Hugo binary you are using e.g. `0.13-DEV`<br>

.Hugo.Environment
: the current running environment as defined through the `--environment` cli tag.

.Hugo.CommitHash
: the git commit hash of the current Hugo binary e.g. `0e8bed9ccffba0df554728b46c5bbf6d78ae5247`

.Hugo.BuildDate
: the compile date of the current Hugo binary formatted with RFC 3339 e.g. `2002-10-02T10:00:00-05:00`<br>



{{% note "Use the Hugo Generator Tag" %}}
We highly recommend using `.Hugo.Generator` in your website's `<head>`. `.Hugo.Generator` is included by default in all themes hosted on [themes.gohugo.io](https://themes.gohugo.io). The generator tag allows the Hugo team to track the usage and popularity of Hugo.
{{% /note %}}

