
---
date: 2020-03-21
title: "Minify config and more!"
description: "Hugo 0.68.0 brings minify configuration and fully navigable headless sections."
categories: ["Releases"]
---

	
This release (finally) brings minify configuration, a big thanks to [@satotake](https://github.com/satotake) for that contribution. See [Configure Minify](https://gohugo.io/getting-started/configuration/#configure-minify) for details.

We have also extended the [Page Build Options](https://gohugo.io/content-management/build-options/) to allow fully navigable headless sections.

This release represents **17 contributions by 6 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@evankanderson](https://github.com/evankanderson), [@QuLogic](https://github.com/QuLogic), and [@le0tan](https://github.com/le0tan) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **3 contributions by 3 contributors**.

Hugo now has:

* 42462+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 300+ [themes](http://themes.gohugo.io/)

## Enhancements

### Other

* Some minify configuration adjustments [7204b354](https://github.com/gohugoio/hugo/commit/7204b354a9f46778f068a4712447d6d4fefbefd8) [@bep](https://github.com/bep) 
* Add minify config [574c2959](https://github.com/gohugoio/hugo/commit/574c2959b8d3338764fa1db102a5e0fd6ed322d9) [@satotake](https://github.com/satotake) [#6750](https://github.com/gohugoio/hugo/issues/6750)[#6892](https://github.com/gohugoio/hugo/issues/6892)
* Allow headless bundles to list pages via $page.Pages and $page.RegularPages [99958f90](https://github.com/gohugoio/hugo/commit/99958f90fedec11d749a1397300860aa8e8459c2) [@bep](https://github.com/bep) [#7075](https://github.com/gohugoio/hugo/issues/7075)
* Update to Go 1.14.1 and 1.13.9 [1d91d8e1](https://github.com/gohugoio/hugo/commit/1d91d8e14b13bd135dc4d4a901fc936c9649b219) [@bep](https://github.com/bep) [#7078](https://github.com/gohugoio/hugo/issues/7078)
* Pass directory name to filters in LstatIfPossible in the same way as Readdir [cc2a5d52](https://github.com/gohugoio/hugo/commit/cc2a5d52a4ad188d93aeb2d51d5c19c7661e098d) [@evankanderson](https://github.com/evankanderson) 
* Update to goldmark 1.1.25. [52c159c4](https://github.com/gohugoio/hugo/commit/52c159c452ab7f48369b5cc9ecc57ecc8dc91654) [@QuLogic](https://github.com/QuLogic) 
* Add workaround for regular CSS imports in SCSS [1a8af7d4](https://github.com/gohugoio/hugo/commit/1a8af7d4f087256710ae0bdf504ed53c0c24a211) [@bep](https://github.com/bep) [#7059](https://github.com/gohugoio/hugo/issues/7059)
* Add .RegularPagesRecursive [03b93bb9](https://github.com/gohugoio/hugo/commit/03b93bb9884ea479c855c2699e8c7b039dce6224) [@bep](https://github.com/bep) [#6411](https://github.com/gohugoio/hugo/issues/6411)
* Add data context to the key in ExecuteAsTemplate [18cb21ff](https://github.com/gohugoio/hugo/commit/18cb21ff2e4a60e7094908e4d6113a9d5a086316) [@bep](https://github.com/bep) [#7046](https://github.com/gohugoio/hugo/issues/7046)
* Improve Tailwind/PostCSS error messages [df298558](https://github.com/gohugoio/hugo/commit/df298558a5a5b747288d9656402af85e0ac75a43) [@bep](https://github.com/bep) [#7041](https://github.com/gohugoio/hugo/issues/7041)[#7042](https://github.com/gohugoio/hugo/issues/7042)
* Update Blackfriday [b1106f87](https://github.com/gohugoio/hugo/commit/b1106f8715cac3544b8ea662b969336fe56fa047) [@bep](https://github.com/bep) [#7039](https://github.com/gohugoio/hugo/issues/7039)
* Add languageDirection to language configuration [5914f91b](https://github.com/gohugoio/hugo/commit/5914f91b6c980e42693661d5fd5640e237691df6) [@le0tan](https://github.com/le0tan) [#6550](https://github.com/gohugoio/hugo/issues/6550)

## Fixes

### Other

* Fix Go build version [2ebb9f54](https://github.com/gohugoio/hugo/commit/2ebb9f5484162062c74698237bcdaa31cb8666b9) [@bep](https://github.com/bep) 
* Fix GetTerms nil pointer [95f49211](https://github.com/gohugoio/hugo/commit/95f492114e33fc6e4d9dcfd2b7c1eca5c50d755f) [@carlmjohnson](https://github.com/carlmjohnson) [#7061](https://github.com/gohugoio/hugo/issues/7061)
* Fix scss vs css import regexp [c7b6d74e](https://github.com/gohugoio/hugo/commit/c7b6d74e898c78da9f5e272e528ff9654206576e) [@bep](https://github.com/bep) [#7063](https://github.com/gohugoio/hugo/issues/7063)
* Fix --templateMetricsHints [5eadc4c0](https://github.com/gohugoio/hugo/commit/5eadc4c0a8e5c51e72670591c4b7877e79c15e3c) [@bep](https://github.com/bep) [#7048](https://github.com/gohugoio/hugo/issues/7048)
* Try to fix a Go 1.15 go vet error [c0177fe2](https://github.com/gohugoio/hugo/commit/c0177fe2b28eb09d1534e62370849c3f1d70b40f) [@bep](https://github.com/bep) 





