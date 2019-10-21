---
date: 2017-07-07T17:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.25 automatically opens the page you&#39;re working on in the browser"
link: ""
title: "Hugo 0.25"
draft: false
author: bep
aliases: [/0-25/]
---

Hugo `0.25` is the **Kinder Surprise**: It automatically opens the page you&#39;re working on in the browser, it adds full `AND` and `OR` support in page queries, and you can now have templates per language.

![Hugo Open on Save](https://cdn-standard5.discourse.org/uploads/gohugo/optimized/2X/6/622088d4a8eacaf62bbbaa27dab19d789e10fe09_1_690x345.gif "Hugo Open on Save")

If you start with `hugo server --navigateToChanged`, Hugo will navigate to the relevant page on save (see animated GIF). This is extremely useful for site-wide edits. Another very useful feature in this version is the added support for `AND` (`intersect`)  and `OR` (`union`)  filters when combined with `where`.

Example:

```
{{ $pages := where .Site.RegularPages "Type" "not in" (slice "page" "about") }}
{{ $pages := $pages | union (where .Site.RegularPages "Params.pinned" true) }}
{{ $pages := $pages | intersect (where .Site.RegularPages "Params.images" "!=" nil) }}
```

The above fetches regular pages not of `page` or `about` type unless they are pinned. And finally, we exclude all pages with no `images` set in Page params.

This release represents **36 contributions by 12 contributors** to the main Hugo code base. [@bep](https://github.com/bep) still leads the Hugo development with his witty Norwegian humor, and once again contributed a significant amount of additions. But also a big shoutout to [@yihui](https://github.com/yihui), [@anthonyfok](https://github.com/anthonyfok), and [@kropp](https://github.com/kropp) for their ongoing contributions. And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Hugo now has:

* 18209&#43; [stars](https://github.com/gohugoio/hugo/stargazers)
* 455&#43; [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 168&#43; [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Add `Pages` support to `intersect` (`AND`) and `union`(`OR`). This makes the `where` template func even more powerful. [ccdd08d5](https://github.com/gohugoio/hugo/commit/ccdd08d57ab64441e93d6861ae126b5faacdb92f) [@bep](https://github.com/bep) [#3174](https://github.com/gohugoio/hugo/issues/3174)
* Add `math.Log` function. This is very handy for creating tag clouds. [34c56677](https://github.com/gohugoio/hugo/commit/34c566773a1364077e1397daece85b22948dc721) [@artem-sidorenko](https://github.com/artem-sidorenko) 
* Add `WebP` images support [8431c8d3](https://github.com/gohugoio/hugo/commit/8431c8d39d878c18c6b5463d9091a953608df10b) [@bep](https://github.com/bep) [#3529](https://github.com/gohugoio/hugo/issues/3529)
* Only show post&#39;s own keywords in schema.org [da72805a](https://github.com/gohugoio/hugo/commit/da72805a4304a57362e8e79a01cc145767b027c5) [@brunoamaral](https://github.com/brunoamaral) [#2635](https://github.com/gohugoio/hugo/issues/2635)[#2646](https://github.com/gohugoio/hugo/issues/2646)
* Simplify the `Disqus` template a little bit (#3655) [eccb0647](https://github.com/gohugoio/hugo/commit/eccb0647821e9db20ba9800da1b4861807cc5205) [@yihui](https://github.com/yihui) 
* Improve the built-in Disqus template (#3639) [2e1e4934](https://github.com/gohugoio/hugo/commit/2e1e4934b60ce8081a7f3a79191ed204f3098481) [@yihui](https://github.com/yihui) 

### Output

* Support templates per site/language. This is for both regular templates and shortcode templates. [aa6b1b9b](https://github.com/gohugoio/hugo/commit/aa6b1b9be7c9d7322333893b642aaf8c7a5f2c2e) [@bep](https://github.com/bep) [#3360](https://github.com/gohugoio/hugo/issues/3360)

### Core

* Extend the sections API [a1d260b4](https://github.com/gohugoio/hugo/commit/a1d260b41a6673adef679ec4e262c5f390432cf5) [@bep](https://github.com/bep) [#3591](https://github.com/gohugoio/hugo/issues/3591)
* Make `.Site.Sections` return the top level sections [dd9b1baa](https://github.com/gohugoio/hugo/commit/dd9b1baab0cb860a3eb32fd9043bac18cab3f9f0) [@bep](https://github.com/bep) [#3591](https://github.com/gohugoio/hugo/issues/3591)
* Render `404.html` for all languages [41805dca](https://github.com/gohugoio/hugo/commit/41805dca9e40e9b0952e04d06074e6fc91140495) [@mitchchn](https://github.com/mitchchn) [#3598](https://github.com/gohugoio/hugo/issues/3598)

### Other

* Support human-readable `YAML` boolean values in `undraft` [1039356e](https://github.com/gohugoio/hugo/commit/1039356edf747f044c989a5bc0e85d792341ed5d) [@kropp](https://github.com/kropp) 
* `hugo import jekyll` support nested `_posts` directories [7ee1f25e](https://github.com/gohugoio/hugo/commit/7ee1f25e9ef3be8f99c171e8e7982f4f82c13e16) [@coderzh](https://github.com/coderzh) [#1890](https://github.com/gohugoio/hugo/issues/1890)[#1911](https://github.com/gohugoio/hugo/issues/1911)
* Update `Dockerfile` and add Docker optimizations [118f8f7c](https://github.com/gohugoio/hugo/commit/118f8f7cf22d756d8a894ff93551974a806f2155) [@ellerbrock](https://github.com/ellerbrock) 
* Add Blackfriday `joinLines` extension support (#3574) [a5440496](https://github.com/gohugoio/hugo/commit/a54404968a4b36579797f2e7ff7f5eada94866d9) [@choueric](https://github.com/choueric) 
* add `--initial-header-level=2` to rst2html (#3528) [bfce30d8](https://github.com/gohugoio/hugo/commit/bfce30d85972c27c27e8a2caac9db6315f813298) [@frankbraun](https://github.com/frankbraun) 
* Support open "current content page" in browser [c825a731](https://github.com/gohugoio/hugo/commit/c825a7312131b4afa67ee90d593640dee3525d98) [@bep](https://github.com/bep) [#3643](https://github.com/gohugoio/hugo/issues/3643)
* Make `--navigateToChanged` more robust on Windows [30e14cc3](https://github.com/gohugoio/hugo/commit/30e14cc31678ddc204b082ab362f86b6b8063881) [@anthonyfok](https://github.com/anthonyfok) [#3645](https://github.com/gohugoio/hugo/issues/3645)
* Remove the docs submodule [31393f60](https://github.com/gohugoio/hugo/commit/31393f6024416ea1b2e61d1080dfd7104df36eda) [@bep](https://github.com/bep) [#3647](https://github.com/gohugoio/hugo/issues/3647)
* Use `example.com` as homepage for new theme [aff1ac32](https://github.com/gohugoio/hugo/commit/aff1ac3235b6c075d01f7237addf44fecdd36d82) [@anthonyfok](https://github.com/anthonyfok) 

## Fixes

### Templates

* Fix `in` function for JSON arrays [d12cf5a2](https://github.com/gohugoio/hugo/commit/d12cf5a25df00fa16c59f0b2ae282187a398214c) [@bep](https://github.com/bep) [#1468](https://github.com/gohugoio/hugo/issues/1468)

### Other

* Fix handling of `JSON` front matter with escaped quotes [e10e51a0](https://github.com/gohugoio/hugo/commit/e10e51a00827b9fdc1bee51439fef05afc529831) [@bep](https://github.com/bep) [#3661](https://github.com/gohugoio/hugo/issues/3661)
* Fix typo in code comment [56d82aa0](https://github.com/gohugoio/hugo/commit/56d82aa025f4d2edb1dc6315132cd7ab52df649a) [@dvic](https://github.com/dvic) 