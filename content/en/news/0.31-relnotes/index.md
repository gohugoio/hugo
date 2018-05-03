
---
date: 2017-11-20
title: "Hugo 0.31: Language Multihost Edition!"
description: "Hugo 0.31: Multihost, smart union static dirs, and more ..."
categories: ["Releases"]
images:
- images/blog/hugo-31-poster.png
---

	Hugo `0.31` is the **Language Multihost Edition!**

> <img src="https://esolia.com/img/eSolia-Logo-Flat-2015.svg" alt="eSolia" width="100px" align="top" style="width:100px" />The Multihost feature  is sponsored by [eSolia](https://esolia.com/), [@rickcogley](https://github.com/rickcogley)'s company.

[Multihost](https://gohugo.io/content-management/multilingual/#configure-multilingual-multihost) means that you can have a **`baseURL` per language**, for example, `https://no.example.com` and `https://en.example.com`. This is seamlessly integrated, and the built-in web server with live reload and `navigateToChanged` etc. just works. A related enhancement in this release is the support for **as many static dirs as you need**, with intelligent language overrides, forming a big union file system. Add to that several other language related fixes and enhancements, it is safe to say that this is the version you want for multilingual Hugo sites!

This release represents **44 contributions by 7 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@kaushalmodi](https://github.com/kaushalmodi), [@natefinch](https://github.com/natefinch), and [@betaveros](https://github.com/betaveros) for their ongoing contributions.
And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **13 contributions by 9 contributors**. A special thanks to [@oncletom](https://github.com/oncletom), [@kaushalmodi](https://github.com/kaushalmodi), [@XhmikosR](https://github.com/XhmikosR), and [@digitalcraftsman](https://github.com/digitalcraftsman) for their work on the documentation site.


Hugo now has:

* 21105+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 455+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 184+ [themes](http://themes.gohugo.io/)

## Notes

* For mapping of translated content, Hugo now considers the full path of the content file, which makes it possible with translation of duplicate content filenames such as `index.md`. A specific translation key can be specified with the new `translationKey` front matter variable. See [#2699](https://github.com/gohugoio/hugo/issues/2699).


## Enhancements

### Language related

* Support unknown language codes [23ba779f](https://github.com/gohugoio/hugo/commit/23ba779fab90ce45cddd68b4f49a2515ce6d4878) [@bep](https://github.com/bep) [#3564](https://github.com/gohugoio/hugo/issues/3564)
* Fix `.IsTranslated`  with identical filenames [b3daa1f4](https://github.com/gohugoio/hugo/commit/b3daa1f4bf1b84bcc5da028257ba609be74e3ecc) [@bep](https://github.com/bep) [#2699](https://github.com/gohugoio/hugo/issues/2699)
* Fall back to unstranslated base template [0a81a6b4](https://github.com/gohugoio/hugo/commit/0a81a6b4bae3de53aa9c179b855c671a2d30eec7) [@bep](https://github.com/bep) [#3893](https://github.com/gohugoio/hugo/issues/3893)
* Add support for multiple static dirs [60dfb9a6](https://github.com/gohugoio/hugo/commit/60dfb9a6e076200ab3ca3fd30e34bb3c14e0a893) [@bep](https://github.com/bep) [#36](https://github.com/gohugoio/hugo/issues/36)[#4027](https://github.com/gohugoio/hugo/issues/4027)
* Add multilingual multihost support [2e046576](https://github.com/gohugoio/hugo/commit/2e0465764b5dacc511b977b1c9aa07324ad0ee9c) [@bep](https://github.com/bep) [#4027](https://github.com/gohugoio/hugo/issues/4027)

### Templates

* Refactor `Mod` with `cast` [76dc811c](https://github.com/gohugoio/hugo/commit/76dc811c6539b2ed8b4d3b22693e5088b9f6ecfe) [@artem-sidorenko](https://github.com/artem-sidorenko) 
* Add support for height argument to figure shortcode [488631fe](https://github.com/gohugoio/hugo/commit/488631fe0abc3667355345c7eb98ba7a2204deb5) [@kaushalmodi](https://github.com/kaushalmodi) [#4014](https://github.com/gohugoio/hugo/issues/4014)

### Core

* Use ms precision for static change logging [bb048d81](https://github.com/gohugoio/hugo/commit/bb048d811d3977adb10656335cd339cd8c945a25) [@bep](https://github.com/bep) 
* Update Chroma to get the latest SASS lexer [b32ffed6](https://github.com/gohugoio/hugo/commit/b32ffed6abc67646cad89e163846f3ffef29cec8) [@bep](https://github.com/bep) [#4069](https://github.com/gohugoio/hugo/issues/4069)
* Bump to Go 1.9.2 [9299a16c](https://github.com/gohugoio/hugo/commit/9299a16c9952a284d3ac3f31d2662f1812f77768) [@bep](https://github.com/bep) [#4064](https://github.com/gohugoio/hugo/issues/4064)
* Update Travis and snapcraft to Go 1.9.2 [77cbd001](https://github.com/gohugoio/hugo/commit/77cbd001ff6b2e0aaa48566ef2af49ca68e19af9) [@bep](https://github.com/bep) [#4064](https://github.com/gohugoio/hugo/issues/4064)
* Handle Taxonomy permalinks [d9a78b61](https://github.com/gohugoio/hugo/commit/d9a78b61adefe8e1803529f4774185874af85148) [@betaveros](https://github.com/betaveros) [#1208](https://github.com/gohugoio/hugo/issues/1208)


### Other

* Support Fast Render mode with sub-path in baseURL [31641033](https://github.com/gohugoio/hugo/commit/3164103310fbca1211cfa9ce4a5eb7437854b6ad) [@bep](https://github.com/bep) [#3981](https://github.com/gohugoio/hugo/issues/3981)
* Simplify Site benchmarks [c3c10f2c](https://github.com/gohugoio/hugo/commit/c3c10f2c7ce4ee11186f51161943efc8b37a28c9) [@bep](https://github.com/bep) 
* Replace `make` with `mage` to build Hugo [#3969](https://github.com/gohugoio/hugo/issues/3969)
* Convert to `dep` as  dependency/vendor manager for Hugo [#3988](https://github.com/gohugoio/hugo/issues/3988)
* Pre-allocate some slices [a9be687b](https://github.com/gohugoio/hugo/commit/a9be687b81df01c7343f78f0d3760042f467baa4) [@bep](https://github.com/bep) 

## Fixes

### Templates

* Make sure only one instance of a cached partial is rendered [#4086](https://github.com/gohugoio/hugo/issues/4086)

### Other

* Fix broken shortcodes for `Ace` and `Amber` [503ca6de](https://github.com/gohugoio/hugo/commit/503ca6de6ceb0b4af533f9efeff917d6f3871278) [@bep](https://github.com/bep) [#4051](https://github.com/gohugoio/hugo/issues/4051)
* Fix error handling in `mage` build [c9c19d79](https://github.com/gohugoio/hugo/commit/c9c19d794537cf76ff281788c3d6cf5f2beac54d) [@natefinch](https://github.com/natefinch) 
* Fix `hugo -w` [fa53b13c](https://github.com/gohugoio/hugo/commit/fa53b13ca0ffb1db6ed20f5353661d3f8a5fd455) [@bep](https://github.com/bep) [#3980](https://github.com/gohugoio/hugo/issues/3980)

