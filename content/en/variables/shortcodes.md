---
title: Shortcode Variables
linktitle: Shortcode Variables
description: Shortcodes can access page variables and also have their own specific built-in variables.
date: 2017-03-12
publishdate: 2017-03-12
lastmod: 2017-03-12
categories: [variables and params]
keywords: [shortcodes]
draft: false
menu:
  docs:
    parent: "variables"
    weight: 20
weight: 20
sections_weight: 20
aliases: []
toc: false
---

[Shortcodes][shortcodes] have access to parameters delimited in the shortcode declaration via [`.Get`][getfunction], page- and site-level variables, and also the following shortcode-specific fields:

.Name
: Shortcode name.

.Ordinal
: Zero-based ordinal in relation to its parent. If the parent is the page itself, this ordinal will represent the position of this shortcode in the page content.

.Parent
: provides access to the parent shortcode context in nested shortcodes. This can be very useful for inheritance of common shortcode parameters from the root.

.Position
: Contains [filename and position](https://godoc.org/github.com/gohugoio/hugo/common/text#Position) for the shortcode in a page. Note that this can be relatively expensive to calculate, and is meant for error reporting. See [Error Handling in Shortcodes](/templates/shortcode-templates/#error-handling-in-shortcodes).




.IsNamedParams
: boolean that returns `true` when the shortcode in question uses [named rather than positional parameters][shortcodes]

.Inner
: represents the content between the opening and closing shortcode tags when a [closing shortcode][markdownshortcode] is used

[getfunction]: /functions/get/
[markdownshortcode]: /content-management/shortcodes/#shortcodes-with-markdown
[shortcodes]: /templates/shortcode-templates/


