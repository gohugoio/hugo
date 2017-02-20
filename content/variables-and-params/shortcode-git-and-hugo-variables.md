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
aliases: [/extras/gitinfo/,/variables-and-params/other/]
toc: true
needsreview: true
notesforauthors:
---

## Shortcode Variables

[Shortcodes][shortcodes] have access to parameters delimited in the shortcode declaration via [`.Get`][getfunction], page- and site-level variables, and also the following shortcode-specific fields:

`.Parent`
: Provides access to the parent shortcode context in nested shortcodes. This can be very useful for inheritance of common shortcode parameters from the root.

`.IsNamedParams`
: Boolean that returns `true` when the shortcode in question uses [named rather than positional parameters][shortcodes]

`.Inner`
: Represents the content between the opening and closing shortcode tags when a [closing shortcode][markdownshortcode] is used

## Git Variables

Hugo provides a way to integrate Git data into your website.

{{% note "`.GitInfo` Performance Considerations"  %}}
Hugo's Git integrations should be fairly performant but *can* increase your build time. This will depend on the size of your Git history.
{{% /note %}}

### `.GitInfo` Prerequisites

1. The Hugo site must be in a Git-enabled directory.
2. The Git executable must be installed and in your system `PATH`.
3. The `.GitInfo` feature must be enabled in your Hugo project by passing `--enableGitInfo` flag on the command line or by setting `enableGitInfo` to `true` in your [site's configuration file][configuration].

### The `.GitInfo` Object

The `GitInfo` object contains the following fields:

`.AbbreviatedHash`
: The abbreviated commit hash (e.g., `866cbcc`)

`.AuthorName`
: The author's name, respecting `.mailmap`

`.AuthorEmail`
: The author's email address, respecting `.mailmap`

`.AuthorDate`
: The author date

`.Hash`
: The commit hash (e.g., `866cbccdab588b9908887ffd3b4f2667e94090c3`)

`.Subject`
: commit message subject (e.g., `tpl: Add custom index function`)

## Hugo Variables

The `.Hugo` variable provides easy access to Hugo-related data and contains the following fields:

`.Hugo.Generator`
: Meta tag for the version of Hugo that generated the site. `.Hugo.Generator` outputs a *complete* HTML tag; e.g. `<meta name="generator" content="Hugo 0.18" />`

`.Hugo.Version`
: The current version of the Hugo binary you are using e.g. `0.13-DEV`<br>

`.Hugo.CommitHash`
: The git commit hash of the current Hugo binary e.g. `0e8bed9ccffba0df554728b46c5bbf6d78ae5247`<br>

`.Hugo.BuildDate`
: The compile date of the current Hugo binary formatted with RFC 3339 e.g. `2002-10-02T10:00:00-05:00`<br>

{{% note "Use the Hugo Generator Tag" %}}
We highly recommend using `.Hugo.Generator` in your website. It is already included in all theme headers. The generator tag is significant in that it allows the Hugo team to track the usage and popularity of Hugo.
{{% /note %}}

[configuration]: /getting-started/configuration/
[getfunction]: /functions/get/
[markdownshortcode]: /content-management/shortcodes/#shortcodes-with-markdown
[shortcodes]: /templates/shortcode-templates/

