---
author: bep
categories: ["Releases"]
date: 2016-10-07T13:54:06-04:00
description: "Hugo now supports multilingual sites with the most simple and elegant experience."
link: ""
title: "0.17: Hugo is going global"
draft: false
aliases: [/0-17/]
---
Hugo is going global with our 0.17 release.  We put a lot of thought into how we could extend Hugo
to support multilingual websites with the most simple and elegant experience. Hugo's multilingual
capabilities rival the best web and documentation software, but Hugo's experience is unmatched.
If you have a single language website, the simple Hugo experience you already love is unchanged.
Adding additional languages to your website is simple and straightforward. Hugo has been completely
internally rewritten to be multilingual aware with translation and internationalization features
embedded throughout Hugo.

Hugo continues its trend of each release being faster than the last. It's quite a challenge to consistently add
significant new functionality and simultaneously dramatically improve performance. {{<gh "@bep">}} has made it
his personal mission to apply the Go mantra of "Enable more. Do less" to Hugo. Hugo's consistent improvement
is a testament to his brilliance and his dedication to his craft. Hugo additionally benefits from the
performance improvements from the Go team in the Go 1.7 release.

This release represents **over 300 contributions by over 70 contributors** to
the main Hugo code base. Since last release Hugo has **gained 2000 stars, 50 new
contributors and 20 additional themes.**

Hugo now has:

* 12,000 stars on GitHub
* 370+ contributors
* 110+ themes

{{<gh "@bep" >}} continues to lead the project with the lionshare of contributions
and reviews. A special thanks to {{<gh "@bep" >}} and {{<gh "@abourget" >}} for their
considerable work on multilingual support.

A big welcome to newcomers {{<gh "@MarkDBlackwell" >}}, {{<gh "@bogem" >}} and
{{<gh "@g3wanghc" >}} for their critical contributions.

### Highlights

**Multilingual Support:**
Hugo now supports multiple languages side-by-side. A single site can now have multiple languages rendered with
full support for translation and i18n.

**Performance:**
Hugo is faster than ever! Hugo 0.17 is not only our fastest release, it's also the most efficient.
Hugo 0.17 is **nearly twice as fast as Hugo 0.16** and uses about 10% less memory.
This means that the same site will build in nearly half the time it took with Hugo 0.16.
For the first time Hugo sites are averaging well under 1ms per rendered content.

**Docs overhaul:**
This release really focused on improving the documentation. [Gohugo.io](http://gohugo.io) is
more accurate and complete than ever.

**Support for macOS Sierra**

### New Features
* Multilingual support {{<gh 2303>}}
* Allow content expiration {{<gh 2137 >}}
* New templates functions:
  * `querify` function to generate query strings inside templates {{<gh 2257>}}
  * `htmlEscape` and `htmlUnescape` template functions {{<gh 2287>}}
  * `time` converts a timestamp string into a time.Time structure {{<gh 2329>}}

### Enhancements

* Render the shortcodes as late as possible {{<gh 0xed0985404db4630d1b9d3ad0b7e41fb186ae0112>}}
* Remove unneeded casts in page.getParam {{<gh 2186 >}}
* Automatic page date fallback {{<gh 2239>}}
* Enable safeHTMLAttr {{<gh 2234>}}
* Add TODO list support for markdown {{<gh 2296>}}
* Make absURL and relURL accept any type {{<gh 2352>}}
* Suppress 'missing static' error {{<gh 2344>}}
* Make summary, wordcount etc. more efficient {{<gh 2378>}}
* Better error reporting in `hugo convert` {{<gh 2440>}}
* Reproducible builds thanks to govendor {{<gh 2461>}}

### Fixes

* Fix shortcode in markdown headers {{<gh 2210 >}}
* Explicitly bind livereload to hugo server port {{<gh 2205>}}
* Fix Emojify for certain text patterns {{<gh 2198>}}
* Normalize file name to NFC {{<gh 2259>}}
* Ignore emacs temp files {{<gh 2266>}}
* Handle symlink change event {{<gh 2273>}}
* Fix panic when using URLize {{<gh 2274>}}
* `hugo import jekyll`: Fixed target path location check {{<gh 2293>}}
* Return all errors from casting in templates {{<gh 2356>}}
* Fix paginator counter on x86-32 {{<gh 2420>}}
* Fix half-broken self-closing shortcodes {{<gh 2499>}}
