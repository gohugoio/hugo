---
date: 2016-12-30T13:54:02-04:00
categories: ["Releases"]
description: "The primary new feature in Hugo 0.18 is that every piece of content is now a Page."
link: ""
title: "Hugo 0.18"
draft: false
author: bep
aliases: [/0-18/]
---

Hugo 0.18.1 is a bug fix release fixing some issues introduced in Hugo 0.18:

* Fix 32-bit binaries {{<gh 2847 >}}
* Fix issues with `preserveTaxonomyNames` {{<gh 2809 >}}
* Fix `.URL` for taxonomy pages when `uglyURLs=true` {{<gh 2819 >}}
* Fix `IsTranslated` and `Translations` for node pages {{<gh 2812 >}}
* Make template error messages more verbose {{<gh 2820 >}}

## **0.18.0** December 19th 2016

Today, we're excited to release the much-anticipated Hugo 0.18!

We're heading towards the end of the year 2016, and we can look back on three releases and a steady growing community around the project.
This release includes **over 220 contributions by nearly 50 contributors** to the main codebase.
Since the last release, Hugo has **gained 1750 stars and 27 additional themes**.

Hugo now has:

- 13750+ stars
- 408+ contributors
- 137+ themes

{{< gh "@bep" >}} once again took the lead of Hugo and contributed a significant amount of additions.
Also a big shoutout to {{< gh "@digitalcraftsman" >}} for his relentless work on keeping the documentation and the themes site in pristine condition,
and also a big thanks to {{< gh "@moorereason" >}} and {{< gh "@bogem" >}} for their contributions.

We wish you all a Merry Christmas and a Happy New Year.<br>
*The Hugo team*

### Highlights

The primary new feature in Hugo 0.18 is that every piece of content is now a `Page` ({{<gh 2297>}}). This means that every page, including the homepage, can have a content file with front matter.

Not only is this a much simpler model to understand, it is also faster and paved the way for several important new features:

* Enable proper titles for Nodes {{<gh 1051>}}
* Sitemap.xml should include nodes, as well as pages {{<gh 1303>}}
* Document homepage content workaround {{<gh 2240>}}
* Allow home page to be easily authored in markdown {{<gh 720>}}
* Minimalist website with homepage as content {{<gh 330>}}

Hugo again continues its trend of each release being faster than the last. It's quite a challenge to consistently add significant new functionality and simultaneously dramatically improve performance. Running [this benchmark]( https://github.com/bep/hugo-benchmark) with [these sites](https://github.com/bep/hugo-benchmark/tree/master/sites) (renders to memory) shows about   60% reduction in time spent and 30% reduction in memory usage compared to Hugo 0.17.

### Other New Features

* Every `Page` now has a `Kind` property. Since everything is a `Page` now, the `Kind` is used to differentiate different kinds of pages.
  Possible values are `page`, `home`, `section`, `taxonomy`, and `taxonomyTerm`.
  (Internally, we also define `RSS`, `sitemap`, `robotsTXT`, and `404`, but those have no practical use for end users at the moment since they are not included in any collections.)
* Add a `GitInfo` object to `Page` if `enableGitInfo` is set. It then also sets `Lastmod` for the given `Page` to the author date provided by Git. {{<gh 2291>}}
* Implement support for alias templates  {{<gh 2533 >}}
* New template functions:
  * Add `imageConfig` function {{<gh 2677>}}
  * Add `sha256` function {{<gh 2762>}}
  * Add `partialCached` template function {{<gh 1368>}}
* Add shortcode to display Instagram images {{<gh 2690>}}
* Add `noChmod` option to disable perm sync {{<gh 2749>}}
* Add `quiet` build mode {{<gh 1218>}}


### Notices

* `.Site.Pages` will now contain *several kinds of pages*, including regular pages, sections, taxonomies, and the home page.
  If you want a specific kind of page, you can filter it with `where` and `Kind`.
  `.Site.RegularPages` is a shortcut to the page collection you have been used to getting.
* `RSSlink` is now deprecated.  Use `RSSLink` instead.
  Note that in Hugo 0.17 both of them existed, so there is a fifty-fifty chance you will not have to do anything
  (if you use a theme, the chance is close to 0), and `RSSlink` will still work for two Hugo versions.

### Fixes

* Revise the `base` template lookup logic so it now better matches the behavior of regular templates, making it easier to override the master templates from the theme {{<gh 2783>}}
* Add workaround for `block` template crash.
  Block templates are very useful, but there is a bug in Go 1.6 and 1.7 which makes the template rendering crash if you use the block template in more complex scenarios.
  This is fixed in the upcoming Go 1.8, but Hugo adds a temporary workaround in Hugo 0.18. {{<gh 2549>}}
* All the `Params` configurations are now case insensitive {{<gh 1129>}} {{<gh 2590>}} {{<gh 2615>}}
* Make RawContent raw again {{<gh 2601>}}
* Fix archetype title and date handling {{<gh 2750>}}
* Fix TOML archetype parsing in `hugo new` {{<gh 2745>}}
* Fix page sorting when weight is zero {{<gh 2673>}}
* Fix page names that contain dot {{<gh 2555>}}
* Fix RSS Title regression {{<gh 2645>}}
* Handle ToC before handling shortcodes {{<gh 2433>}}
* Only watch relevant themes dir {{<gh 2602>}}
* Hugo new content creates TOML slices with closing bracket on new line {{<gh 2800>}}

### Improvements

* Add page information to error logging in rendering {{<gh 2570>}}
* Deprecate `RSSlink` in favor of `RSSLink`
* Make benchmark command more useful {{<gh 2432>}}
* Consolidate the `Param` methods {{<gh 2590>}}
* Allow to set cache dir in config file
* Performance improvements:
  * Avoid repeated Viper loads of `sectionPagesMenu` {{<gh 2728>}}
  * Avoid reading from Viper for path and URL funcs {{<gh 2495>}}
  * Add `partialCached` template function. This can be a significant performance boost if you have complex partials that does not need to be rerendered for every page. {{<gh 1368>}}

### Documentation Updates

* Update roadmap {{<gh 2666>}}
* Update multilingual example {{<gh 2417>}}
* Add a "Deployment with rsync" tutorial page {{<gh 2658>}}
* Refactor `/docs` to use the `block` keyword {{<gh 2226>}}
