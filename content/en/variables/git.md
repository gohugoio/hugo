---
title: Git Info Variables
linktitle: Git Variables
description: Get the last Git revision information for every content file.
date: 2017-03-12
publishdate: 2017-03-12
lastmod: 2017-03-12
categories: [variables and params]
keywords: [git]
draft: false
menu:
  docs:
    parent: "variables"
    weight: 70
weight: 70
sections_weight: 70
aliases: [/extras/gitinfo/]
toc: false
wip: false
---

{{% note "`.GitInfo` Performance Considerations"  %}}
Hugo's Git integrations should be fairly performant but *can* increase your build time. This will depend on the size of your Git history.
{{% /note %}}

## `.GitInfo` Prerequisites

1. The Hugo site must be in a Git-enabled directory.
2. The Git executable must be installed and in your system `PATH`.
3. The `.GitInfo` feature must be enabled in your Hugo project by passing `--enableGitInfo` flag on the command line or by setting `enableGitInfo` to `true` in your [site's configuration file][configuration].

## The `.GitInfo` Object

The `GitInfo` object contains the following fields:

.AbbreviatedHash
: the abbreviated commit hash (e.g., `866cbcc`)

.AuthorName
: the author's name, respecting `.mailmap`

.AuthorEmail
: the author's email address, respecting `.mailmap`

.AuthorDate
: the author date

.Hash
: the commit hash (e.g., `866cbccdab588b9908887ffd3b4f2667e94090c3`)

.Subject
: commit message subject (e.g., `tpl: Add custom index function`)

## `.Lastmod`

If the `.GitInfo` feature is enabled, `.Lastmod` (on `Page`) is fetched from Git i.e. `.GitInfo.AuthorDate`. This behaviour can be changed by adding your own [front matter configuration for dates](/getting-started/configuration/#configure-front-matter).

[configuration]: /getting-started/configuration/
