
---
date: 2019-02-01
title: "0.54.0:  Mostly Bugfixes"
description: "0.54.0 is mostly a bugfix-release, but also some nice improvements."
categories: ["Releases"]

---

This release represents **27 contributions by 7 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@tryzniak](https://github.com/tryzniak), [@anthonyfok](https://github.com/anthonyfok), and [@mywaiting](https://github.com/mywaiting) for their ongoing contributions. And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), which has received **38 contributions by 17 contributors**. A special thanks to [@bep](https://github.com/bep), [@kaushalmodi](https://github.com/kaushalmodi), [@onedrawingperday](https://github.com/onedrawingperday), and [@peaceiris](https://github.com/peaceiris) for their work on the documentation site.

Hugo now has:

* 32265+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 289+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Adjust tests [ddc6d4e3](https://github.com/gohugoio/hugo/commit/ddc6d4e30f282f23b703a3b42da552886062c8c8) [@bep](https://github.com/bep) [#5643](https://github.com/gohugoio/hugo/issues/5643)
* Prevent getJSON and getCSV fetch failure from aborting build [6a2bfcbe](https://github.com/gohugoio/hugo/commit/6a2bfcbec8df14b1741dbe9b5ead08158bf7adb9) [@anthonyfok](https://github.com/anthonyfok) [#5643](https://github.com/gohugoio/hugo/issues/5643)

### Core

* Expand TestPageWithEmoji to cover '+', '-' and '_' too [2a9060a8](https://github.com/gohugoio/hugo/commit/2a9060a85ce430b28f5ec47e1438c6ef1b8e13fa) [@anthonyfok](https://github.com/anthonyfok) [#5635](https://github.com/gohugoio/hugo/issues/5635)
* Restore 0.48 slash handling in taxonomies [40ffb048](https://github.com/gohugoio/hugo/commit/40ffb0484b96b7b77fb66202b33073b241807199) [@bep](https://github.com/bep) [#5571](https://github.com/gohugoio/hugo/issues/5571)

### Other

* Use official semver even for main releases [fab41f42](https://github.com/gohugoio/hugo/commit/fab41f42d3e23c11651ab75413b01d97e5d37c30) [@bep](https://github.com/bep) [#5639](https://github.com/gohugoio/hugo/issues/5639)
* Add test for --configDir [59d87044](https://github.com/gohugoio/hugo/commit/59d87044a4146f578b92b3d67b46660212940912) [@bep](https://github.com/bep) [#5662](https://github.com/gohugoio/hugo/issues/5662)
* Ignore unknown config files in config dir [3244cb3b](https://github.com/gohugoio/hugo/commit/3244cb3b31f8f8c39d9dfa82bc01fb2d6db59257) [@bep](https://github.com/bep) [#5646](https://github.com/gohugoio/hugo/issues/5646)
* Store supported config formats in a variable [d9282cf9](https://github.com/gohugoio/hugo/commit/d9282cf98a346fcf98f363d9c353e4920ca85fc7) [@tryzniak](https://github.com/tryzniak) 
* Bump to Go 1.11.5 [8ed2a1ca](https://github.com/gohugoio/hugo/commit/8ed2a1caa9e0892d5bf97ed1b7279befa159f764) [@bep](https://github.com/bep) [#5654](https://github.com/gohugoio/hugo/issues/5654)
* Update Afero [e8596139](https://github.com/gohugoio/hugo/commit/e85961390a050cd4f2e6ce4f2666012bc83bb449) [@bep](https://github.com/bep) [#5650](https://github.com/gohugoio/hugo/issues/5650)
* Accept hyphen and plus sign in emoji detection [3038464e](https://github.com/gohugoio/hugo/commit/3038464ea6f931c8a08ee49d47f1eaec99ba4817) [@anthonyfok](https://github.com/anthonyfok) [#5635](https://github.com/gohugoio/hugo/issues/5635)
* Support numeric sort in ByParam [26f75edb](https://github.com/gohugoio/hugo/commit/26f75edb7a76c816349749a05edf98fb36dc338a) [@tryzniak](https://github.com/tryzniak) [#5305](https://github.com/gohugoio/hugo/issues/5305)
* Make hugo server -t work again [db3c49d0](https://github.com/gohugoio/hugo/commit/db3c49d049193e0fc225fe4bdb95712c311d6615) [@tryzniak](https://github.com/tryzniak) [#5569](https://github.com/gohugoio/hugo/issues/5569)[#5061](https://github.com/gohugoio/hugo/issues/5061)[#4868](https://github.com/gohugoio/hugo/issues/4868)
* Add configFile(s) back to the watch list after RENAME event too [e3cb8e6c](https://github.com/gohugoio/hugo/commit/e3cb8e6c7874d7dfe1d4d1c7f5c9765b681fb647) [@anthonyfok](https://github.com/anthonyfok) [#5205](https://github.com/gohugoio/hugo/issues/5205)
* Remove historical rssURI config [55251aa8](https://github.com/gohugoio/hugo/commit/55251aa89099358c040d38f3af48e3699d67bab2) [@mywaiting](https://github.com/mywaiting) 
* Use subtests with server_test.go [843fcd19](https://github.com/gohugoio/hugo/commit/843fcd19d4d97bac979410a4e0abed72586a0aa0) [@tryzniak](https://github.com/tryzniak) 
* Move resource interfaces into its own package [ce8a09a4](https://github.com/gohugoio/hugo/commit/ce8a09a4c0661dece931ab1173e4f09e8e04aa38) [@bep](https://github.com/bep) 
* Move resource processors into sub-packages [669ada43](https://github.com/gohugoio/hugo/commit/669ada436787311cc5d02dae5b88e60a09adda58) [@bep](https://github.com/bep) 
* Update _index.md [50745122](https://github.com/gohugoio/hugo/commit/507451229c2255788d72b757a85ad5bb3ba00f4f) [@vrMarc](https://github.com/vrMarc) 
* Update go.sum [0584432b](https://github.com/gohugoio/hugo/commit/0584432b078f1e3a488ad4f27f39edac0557e042) [@bep](https://github.com/bep) 
* Update Chroma [cc351958](https://github.com/gohugoio/hugo/commit/cc351958e12d4dc83f664a1d51be76a447fea9b8) [@bep](https://github.com/bep) [#4993](https://github.com/gohugoio/hugo/issues/4993)
* Make docshelper run again [c24f3ae2](https://github.com/gohugoio/hugo/commit/c24f3ae22b27dfe5339662277f8183596a6d148d) [@bep](https://github.com/bep) [#5568](https://github.com/gohugoio/hugo/issues/5568)

## Fixes

### Templates

* Fix reflect [9e4f9e0b](https://github.com/gohugoio/hugo/commit/9e4f9e0bb69276e9bca0dfbdbc7aefbf5f6fc9e5) [@moorereason](https://github.com/moorereason) [#5564](https://github.com/gohugoio/hugo/issues/5564)

### Other

* Fix some inline shortcode issues [c52045bb](https://github.com/gohugoio/hugo/commit/c52045bbb38cbf64b9cb39352230060aa122cc9f) [@bep](https://github.com/bep) [#5645](https://github.com/gohugoio/hugo/issues/5645)[#5653](https://github.com/gohugoio/hugo/issues/5653)
* Fix OpenGraph image fallback to site params [526b5b1c](https://github.com/gohugoio/hugo/commit/526b5b1c4986d43d6184671b02f45ca40f041b65) [@statik](https://github.com/statik) 
* Fix Params case handling in the new site global [e1a66c73](https://github.com/gohugoio/hugo/commit/e1a66c7343db9d232749255dd9e3a58d94b86997) [@bep](https://github.com/bep) [#5615](https://github.com/gohugoio/hugo/issues/5615)
* cache/namedmemcache: Fix data race [3f3187de](https://github.com/gohugoio/hugo/commit/3f3187de0f62107da19d9341aebd1d8414bff0ea) [@bep](https://github.com/bep) 





