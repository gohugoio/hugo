---
date: 2017-02-27T13:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.19 brings native Emacs Org-mode content support, and Hugo has its own Twitter account"
link: ""
title: "Hugo 0.19"
draft: false
author: budparr
aliases: [/0-19/]
---

We're happy to announce the first release of Hugo in 2017.

This release represents **over 180 contributions by over 50 contributors** to the main Hugo code base. Since last release Hugo has **gained 1450 stars, 35 new contributors, and 15 additional themes.**

Hugo now has:

* 15200+ stars
* 470+ contributors
* 151+ themes

Furthermore, Hugo has its own Twitter account ([@gohugoio](https://twitter.com/gohugoio)) where we share bite-sized news and themes from the Hugo community.

{{< gh "@bep" >}} leads the Hugo development and once again contributed a significant amount of additions. Also a big shoutout to  {{< gh "@chaseadamsio" >}} for the Emacs Org-mode support, {{< gh "@digitalcraftsman" >}} for his relentless work on keeping the documentation and the themes site in pristine condition, {{< gh "@fj" >}}for his work on revising the `params` handling in Hugo, and {{< gh "@moorereason" >}} and {{< gh "@bogem" >}} for their ongoing contributions.

### Highlights

Hugo `0.19` brings native Emacs Org-mode content support ({{<gh 1483>}}), big thanks to {{< gh "@chaseadamsio" >}}.

Also, a considerably amount of work have been put into cleaning up the Hugo source code, in an issue titled [Refactor the globals out of site build](https://github.com/gohugoio/hugo/issues/2701). This is not immediately visible to the Hugo end user, but will speed up future development.

Hugo `0.18` was bringing full-parallel page rendering, so workarounds depending on rendering order did not work anymore, and pages with duplicate target paths (common examples would be `/index.md` or `/about/index.md`) would now conflict with the home page or the section listing.

With Hugo `0.19`, you can control this behaviour by turning off page types you do not want ({{<gh 2534 >}}). In its most extreme case, if you put the below setting in your [`config.toml`](/getting-started/configuration/), you will get **nothing!**:

```
disableKinds = ["page", "home", "section", "taxonomy", "taxonomyTerm", "RSS", "sitemap", "robotsTXT", "404"]
```

### Other New Features

* Add ability to sort pages by front matter parameters, enabling easy custom "top 10" page lists. {{<gh 3022 >}}
* Add `truncate` template function {{<gh 2882 >}}
* Add `now` function, which replaces the now deprecated `.Now` {{<gh 2859 >}}
* Make RSS item limit configurable {{<gh 3035 >}}

### Enhancements

*  Enhance `.Param` to permit arbitrarily nested parameter references {{<gh 2598 >}}
* Use `Page.Params` more consistently when adding metadata {{<gh 3033 >}}
* The `sectionPagesMenu` feature ("Section menu for the lazy blogger") is now integrated with the section content pages. {{<gh 2974 >}}
* Hugo `0.19` is compiled with Go 1.8!
* Make template funcs like `findRE` and friends more liberal in what argument types they accept {{<gh 3018 >}} {{<gh 2822 >}}
* Improve generation of OpenGraph date tags {{<gh 2979 >}}

### Notes

* `sourceRelativeLinks` is now deprecated and will be removed in Hugo `0.21` if  no one is stepping up to the plate and fixes and maintains this feature. {{<gh 3028 >}}

### Fixes

* Fix `.Site.LastChange` on sites where the default sort order is not chronological. {{<gh 2909 >}}
* Fix regression of `.Truncated` evaluation in manual summaries. {{<gh 2989 >}}
* Fix `preserveTaxonomyNames` regression {{<gh 3070 >}}
* Fix issue with taxonomies when only some have content page {{<gh 2992 >}}
* Fix instagram shortcode panic on invalid ID {{<gh 3048 >}}
* Fix subtle data race in `getJSON` {{<gh 3045 >}}
* Fix deadlock in cached partials {{<gh 2935 >}}
* Avoid double-encoding of paginator URLs {{<gh 2177 >}}
* Allow tilde in URLs {{<gh 2177 >}}
* Fix `.Site.Pages` handling on live reloads {{<gh 2869 >}}
* `UniqueID` now correctly uses the fill file path from the content root to calculate the hash, and is finally ... unique!
* Discard current language based on `.Lang()`, go get translations correct for paginated pages. {{<gh 2972 >}}
* Fix infinite loop in template AST handling for recursive templates  {{<gh 2927 >}}
* Fix issue with watching when config loading fails {{<gh 2603 >}}
* Correctly flush the imageConfig on live-reload {{<gh 3016 >}}
* Fix parsing of TOML arrays in front matter {{<gh 2752 >}}

### Docs

* Add tutorial "How to use Google Firebase to host a Hugo site" {{<gh 3007 >}}
* Improve documentation for menu rendering {{<gh 3056 >}}
* Revise GitHub Pages deployment tutorial {{<gh 2930 >}}
