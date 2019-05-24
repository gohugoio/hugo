---
date: 2017-06-12T17:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.22 brings nested sections, by popular demand and a long sought after feature"
link: ""
title: "Hugo 0.22"
draft: false
author: bep
aliases: [/0-22/]
---


Hugo `0.22` brings **nested sections**, by popular demand and a long sought after feature ([#465](https://github.com/gohugoio/hugo/issues/465)).  We are still low on documentation for this great feature, but [@bep](https://github.com/bep)  has been kind enough to accompany his implementation with a [demo site](http://hugotest.bep.is/).

This release represents **58 contributions by 10 contributors** to the main Hugo code base. Since last release Hugo has **gained 420 stars and 2 additional themes.**

[@bep](https://github.com/bep) still leads the Hugo development with his witty Norwegian humor, and once again contributed a significant amount of additions. But also a big shoutout to [@bogem](https://github.com/bogem), [@moorereason](https://github.com/moorereason), and [@onedrawingperday](https://github.com/onedrawingperday) for their ongoing contributions. And as always big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Hugo now has:

* 17576&#43; [stars](https://github.com/gohugoio/hugo/stargazers)
* 455&#43; [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 165&#43; [themes](http://themes.gohugo.io/)

## Other Highlights

`.Site.GetPage` can now also be used to get regular pages ([#2844](https://github.com/gohugoio/hugo/issues/2844)):

```
{{ (.Site.GetPage "page" "blog" "mypost.md" ).Title }}
```

Also, considerable work has been put into writing automated benchmark tests for the site builds, and we&#39;re happy to report that although this release comes with fundamental structural changes, this version is -- in general -- even faster than the previous. It’s quite a challenge to consistently add significant new functionality and simultaneously maintain the stellar performance Hugo is famous for. 


 
## Notes

`.Site.Sections` is replaced. We have reworked how sections work in Hugo, they can now be nested and are no longer taxonomies. If you use the old collection, you should get detailed upgrade instructions in the log when you run `hugo`. For more information, see this [demo site](http://hugotest.bep.is/). 

## Enhancements

### Templates

* Add `uint` support to `In` [b82cd82f](https://github.com/gohugoio/hugo/commit/b82cd82f1198a371ed94bda7faafe22813f4cb29) [@moorereason](https://github.com/moorereason) 
* Support interfaces in `union` [204c3a9e](https://github.com/gohugoio/hugo/commit/204c3a9e32fcf6617ede978e35d3e2e89a5b491c) [@moorereason](https://github.com/moorereason) [#3411](https://github.com/gohugoio/hugo/issues/3411) 
* Add `uniq` function [e28d9aa4](https://github.com/gohugoio/hugo/commit/e28d9aa42c3429d22fe254e69e4605aaf1e684f3) [@adiabatic](https://github.com/adiabatic) 
* Handle `template.HTML` and friends in `ToInt` [4113693a](https://github.com/gohugoio/hugo/commit/4113693ac1b275f3a40aa5c248269340ef9b57f6) [@moorereason](https://github.com/moorereason) [#3308](https://github.com/gohugoio/hugo/issues/3308) 


### Core

* Make the `RSS feed` use the date for the node it represents [f1da5a15](https://github.com/gohugoio/hugo/commit/f1da5a15a37666ee59350d6600a8c14c1383f5bc) [@bep](https://github.com/bep) [#2708](https://github.com/gohugoio/hugo/issues/2708) 
* Enable `nested sections` [b3968939](https://github.com/gohugoio/hugo/commit/b39689393ccb8434d9a57658a64b77568c718e99) [@bep](https://github.com/bep) [#465](https://github.com/gohugoio/hugo/issues/465) 
* Add test for &#34;no 404&#34; in `sitemap` [8aaec644](https://github.com/gohugoio/hugo/commit/8aaec644a90d09bd7f079d35d382f76bb4ed35db) [@bep](https://github.com/bep) [#3563](https://github.com/gohugoio/hugo/issues/3563) 
* Support regular pages in `.Site.GetPage` [e0c2e798](https://github.com/gohugoio/hugo/commit/e0c2e798201f75ae6e9a81a7442355288c2d141b) [@bep](https://github.com/bep) [#2844](https://github.com/gohugoio/hugo/issues/2844) 
[#3082](https://github.com/gohugoio/hugo/issues/3082) 

### Performance
* Add site building benchmarks [8930e259](https://github.com/gohugoio/hugo/commit/8930e259d78cba4041b550cc51a7f40bc91d7c20) [@bep](https://github.com/bep) [#3535](https://github.com/gohugoio/hugo/issues/3535) 
* Add a cache to `GetPage` which makes it much faster [50d11138](https://github.com/gohugoio/hugo/commit/50d11138f3e18b545c15fadf52f7b0b744bf3e7c) [@bep](https://github.com/bep) 
* Speed up `GetPage` [fbb78b89](https://github.com/gohugoio/hugo/commit/fbb78b89df8ccef8f0ab26af00aa45d35c1ee2cf) [@bep](https://github.com/bep) [#3503](https://github.com/gohugoio/hugo/issues/3503) 
* Add BenchmarkFrontmatterTags [3d9c4f51](https://github.com/gohugoio/hugo/commit/3d9c4f513b0443648d7e88995e351df1739646d2) [@bep](https://github.com/bep) [#3464](https://github.com/gohugoio/hugo/issues/3464) 
* Add `benchSite.sh` to make it easy to run Hugo performance benchmarks [d74452cf](https://github.com/gohugoio/hugo/commit/d74452cfe8f69a85ec83e05481e16bebf199a5cb) [@bep](https://github.com/bep) 
* Cache language config [4aff2b6e](https://github.com/gohugoio/hugo/commit/4aff2b6e7409a308f30cff1825fec02991e0d56a) [@bep](https://github.com/bep) 
* Temporarily revert to BurntSushi for `TOML` front matter handling; it is currently much faster [0907a5c1](https://github.com/gohugoio/hugo/commit/0907a5c1c293755e6bf297246f07888448d81f8b) [@bep](https://github.com/bep) [#3541](https://github.com/gohugoio/hugo/issues/3541) [#3464](https://github.com/gohugoio/hugo/issues/3464) 
* Add a simple partitioned lazy cache [87203139](https://github.com/gohugoio/hugo/commit/87203139c38e0b992c96d7b8a23c7730649c68e5) [@bep](https://github.com/bep) 

### Other

* Add `noindex` tag to HTML generated by Hugo aliases [d5ab7f08](https://github.com/gohugoio/hugo/commit/d5ab7f087d967b30e7de7d789e6ad3091b42f1f7) [@onedrawingperday](https://github.com/onedrawingperday) 
* Update Go versions [bde807bd](https://github.com/gohugoio/hugo/commit/bde807bd1e560fb4cc765c0fc22132db7f8a0801) [@bep](https://github.com/bep) 
* Remove the `rlimit` tweaking on `macOS` [bcd32f10](https://github.com/gohugoio/hugo/commit/bcd32f1086c8c604fb22a7496924e41cc46b1605) [@bep](https://github.com/bep) [#3512](https://github.com/gohugoio/hugo/issues/3512) 

### Docs
* Rewrite “Archetypes” article [@davidturnbull](https://github.com/davidturnbull) [#3543](https://github.com/gohugoio/hugo/pull/3543/) 
* Remove Unmaintaned Frontends from Tools. [f41f7282](https://github.com/gohugoio/hugo/commit/f41f72822251c9a31031fd5b3dda585c57c8b028) [@onedrawingperday](https://github.com/onedrawingperday) 

## Fixes

### Core
* Improve `live-reload` on directory structure changes making removal of directories or pasting new content directories into  `/content` just work [fe901b81](https://github.com/gohugoio/hugo/commit/fe901b81191860b60e6fcb29f8ebf87baef2ee79) [@bep](https://github.com/bep) [#3570](https://github.com/gohugoio/hugo/issues/3570) 
* Respect `disableKinds=[&#34;sitemap&#34;]` [69d92dc4](https://github.com/gohugoio/hugo/commit/69d92dc49cb8ab9276ab013d427ba2d9aaf9135d) [@bep](https://github.com/bep) [#3544](https://github.com/gohugoio/hugo/issues/3544) 
* Fix `disablePathToLower` regression [5be04486](https://github.com/gohugoio/hugo/commit/5be0448635fdf5fe6b1ee673e869f2b9baf1a5c6) [@bep](https://github.com/bep) [#3374](https://github.com/gohugoio/hugo/issues/3374) 
* Fix `ref`/`relref` issue with duplicate base filenames [612f6e3a](https://github.com/gohugoio/hugo/commit/612f6e3afe0510c31f70f3621f3dc8ba609dade4) [@bep](https://github.com/bep) [#2507](https://github.com/gohugoio/hugo/issues/2507) 

### Docs

* Fix parameter name in `YouTube` shortcode section [37e37877](https://github.com/gohugoio/hugo/commit/37e378773fbc127863f2b7a389d5ce3a14674c73) [@zivbk1](https://github.com/zivbk1) 