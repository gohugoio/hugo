---
date: 2017-04-10T13:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.20 introduces the powerful and long sought after feature Custom Output Formats"
link: ""
title: "Hugo 0.20"
draft: false
author: bep
aliases: [/0-20/]
---

Hugo `0.20` introduces the powerful and long sought after feature [Custom Output Formats](http://gohugo.io/extras/output-formats/); Hugo isn’t just that “static HTML with an added RSS feed” anymore. _Say hello_ to calendars, e-book formats, Google AMP, and JSON search indexes, to name a few ( [#2828](//github.com/gohugoio/hugo/issues/2828) ).

This release represents **over 180 contributions by over 30 contributors** to the main Hugo code base. Since last release Hugo has **gained 1100 stars, 20 new contributors and 5 additional themes.**

Hugo now has:

*   16300+ stars
*   495+ contributors
*   156+ themes

[@bep](//github.com/bep) still leads the Hugo development with his witty Norwegian humor, and once again contributed a significant amount of additions. Also a big shoutout to [@digitalcraftsman](//github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition, and [@moorereason](//github.com/moorereason) and [@bogem](//github.com/bogem) for their ongoing contributions.

## Other Highlights

[@bogem](//github.com/bogem) has also contributed TOML as an alternative and much simpler format for language/i18n files ([#3200](//github.com/gohugoio/hugo/issues/3200)). A feature you will appreciate when you start to work on larger translations.

Also, there have been some important updates in the Emacs Org-mode handling: [@chaseadamsio](//github.com/chaseadamsio) has fixed the newline-handling ( [#3126](//github.com/gohugoio/hugo/issues/3126) ) and [@clockoon](//github.com/clockoon) has added basic footnote support.

Worth mentioning is also the ongoing work that [@rdwatters](//github.com/rdwatters) and [@budparr](//github.com/budparr) is doing to re-do the [gohugo.io](https://gohugo.io/) site, including a total restructuring and partial rewrite of the documentation. It is getting close to finished, and it looks fantastic!

## Notes

*   `RSS` description in the built-in template is changed from full `.Content` to `.Summary`. This is a somewhat breaking change, but is what most people expect from their RSS feeds. If you want full content, please provide your own RSS template.
*   The deprecated `.RSSlink` is now removed. Use `.RSSLink`.
*   `RSSUri` is deprecated and will be removed in a future Hugo version, replace it with an output format definition.
*   The deprecated `.Site.GetParam` is now removed, use `.Site.Param`.
*   Hugo does no longer append missing trailing slash to `baseURL` set as a command line parameter, making it consistent with how it behaves from site config. [#3262](//github.com/gohugoio/hugo/issues/3262)

## Enhancements

*   Hugo `0.20` is built with Go 1.8.1.
*   Add `.Site.Params.mainSections` that defaults to the section with the most pages. Plan is to get themes to use this instead of the hardcoded `blog` in `where` clauses. [#3206](//github.com/gohugoio/hugo/issues/3206)
*   File extension is now configurable. [#320](//github.com/gohugoio/hugo/issues/320)
*   Impove `markdownify` template function performance. [#3292](//github.com/gohugoio/hugo/issues/3292)
*   Add taxonomy terms’ pages to `.Data.Pages` [#2826](//github.com/gohugoio/hugo/issues/2826)
*   Change `RSS` description from full `.Content` to `.Summary`.
*   Ignore “.” dirs in `hugo --cleanDestinationDir` [#3202](//github.com/gohugoio/hugo/issues/3202)
*   Allow `jekyll import` to accept both `2006-01-02` and `2006-1-2` date format [#2738](//github.com/gohugoio/hugo/issues/2738)
*   Raise the default `rssLimit` [#3145](//github.com/gohugoio/hugo/issues/3145)
*   Unify section list vs single template lookup order [#3116](//github.com/gohugoio/hugo/issues/3116)
*   Allow `apply` to be used with the built-in Go template funcs `print`, `printf` and `println`. [#3139](//github.com/gohugoio/hugo/issues/3139)

## Fixes

*   Fix deadlock in `getJSON` [#3211](//github.com/gohugoio/hugo/issues/3211)
*   Make sure empty terms pages are created. [#2977](//github.com/gohugoio/hugo/issues/2977)
*   Fix base template lookup order for sections [#2995](//github.com/gohugoio/hugo/issues/2995)
*   `URL` fixes:
    *   Fix pagination URLs with `baseURL` with sub-root and `canonifyUrls=false` [#1252](//github.com/gohugoio/hugo/issues/1252)
    *   Fix pagination URL for resources with “.” in name [#2110](//github.com/gohugoio/hugo/issues/2110) [#2374](//github.com/gohugoio/hugo/issues/2374) [#1885](//github.com/gohugoio/hugo/issues/1885)
    *   Handle taxonomy names with period [#3169](//github.com/gohugoio/hugo/issues/3169)
    *   Handle `uglyURLs` ambiguity in `Permalink` [#3102](//github.com/gohugoio/hugo/issues/3102)
    *   Fix `Permalink` for language-roots wrong when `uglyURLs` is `true` [#3179](//github.com/gohugoio/hugo/issues/3179)
    *   Fix misc case issues for `URLs` [#1641](//github.com/gohugoio/hugo/issues/1641)
    *   Fix for taxonomies URLs when `uglyUrls=true` [#1989](//github.com/gohugoio/hugo/issues/1989)
    *   Fix empty `RSSLink` for list pages with content page. [#3131](//github.com/gohugoio/hugo/issues/3131)
*   Correctly identify regular pages on the form “my_index_page.md” [#3234](//github.com/gohugoio/hugo/issues/3234)
*   `Exit -1` on `ERROR` in global logger [#3239](//github.com/gohugoio/hugo/issues/3239)
*   Document hugo `help command` [#2349](//github.com/gohugoio/hugo/issues/2349)
*   Fix internal `Hugo` version handling for bug fix releases. [#3025](//github.com/gohugoio/hugo/issues/3025)
*   Only return `RSSLink` for pages that actually have a RSS feed. [#1302](//github.com/gohugoio/hugo/issues/1302)