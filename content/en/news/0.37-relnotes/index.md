
---
date: 2018-02-27
title: "Hugo 0.37: Preserve PNG Colour Palette"
description: "Reduces processed PNG file sizes. And 0.37 is built with Go 1.10!"
categories: ["Releases"]
---

The main item in Hugo `0.37` is that we now properly preserve the colour palette when processing `PNG` images. We got reports from users experiencing their `PNG` images getting bigger in file size when scaled down.   Now, if you, as an example, start out with a carefully chosen 8 bit colour palette (i.e. `PNG-8`), this is now what you will end up with. A special thanks to [@aitva](https://github.com/aitva) for doing the investigative work finding a proper fix for this issue.

This release represents **40 contributions by 5 contributors** to the main Hugo code base.

[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@vassudanagunta](https://github.com/vassudanagunta), [@kaushalmodi](https://github.com/kaushalmodi), and [@curttimson](https://github.com/curttimson) for their ongoing contributions.

And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **24 contributions by 8 contributors**. A special thanks to [@bep](https://github.com/bep), [@4RU](https://github.com/4RU), [@kaushalmodi](https://github.com/kaushalmodi), and [@mitchchn](https://github.com/mitchchn) for their work on the documentation site.

Hugo now has:

* 23649+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 447+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 197+ [themes](http://themes.gohugo.io/)

## Notes

* Hugo will now convert non-string `YAML` map keys to string.  See [#4393](https://github.com/gohugoio/hugo/issues/4393) for more information. You will get a `WARNING` in the console if you are touched by this.
* We have improved the `PNG` processing, and have incremented the version numbers on the URL for the processed `PNG` image. This will just work, but you may want to run `hugo --gc` to clean up some old stale images in the resource cache.

## Enhancements

### Templates

* Add template func for TOML/JSON/YAML docs examples conversion. This is mainly motivated by the needs of the Hugo docs site. [d382502d](https://github.com/gohugoio/hugo/commit/d382502d6dfa1c066545e215ba83e2e0a9d2c8d7) [@bep](https://github.com/bep) [#4389](https://github.com/gohugoio/hugo/issues/4389)

### Core

* Refactor tests for JSON, YAML and TOML equivalency, add coverage [82eefded](https://github.com/gohugoio/hugo/commit/82eefded1353f0198fd8fe9f7df1aa620d3d50eb) [@vassudanagunta](https://github.com/vassudanagunta) 
* Re-enable YAML data tests disabled in f554503f [e1728349](https://github.com/gohugoio/hugo/commit/e1728349930e2cc1b6580766473de98adb0f3e50) [@vassudanagunta](https://github.com/vassudanagunta) 

### Other

* Preserve color palette for PNG images [799c654b](https://github.com/gohugoio/hugo/commit/799c654b0d39ec869c2da24d41de3636eb7157f0) [@bep](https://github.com/bep) [#4416](https://github.com/gohugoio/hugo/issues/4416)
* Use `Floyd-Steinberg` dithering for PNGs [13ea1e7c](https://github.com/gohugoio/hugo/commit/13ea1e7c352852966f88ef119d9434bbb1ee62fa) [@bep](https://github.com/bep) [#4453](https://github.com/gohugoio/hugo/issues/4453)
* Make `ge`, `le` etc. work with the Hugo Version number [0602135f](https://github.com/gohugoio/hugo/commit/0602135fd44b0cfa0a51b0ec6e451ae58ac95666) [@bep](https://github.com/bep) [#4443](https://github.com/gohugoio/hugo/issues/4443)
* Update dependencies [eaf573a2](https://github.com/gohugoio/hugo/commit/eaf573a2778e79287b871b69f4959fd3082d8887) [@bep](https://github.com/bep) [#4418](https://github.com/gohugoio/hugo/issues/4418)
* Update to Go 1.10 (!) Take 2 [a3f26e56](https://github.com/gohugoio/hugo/commit/a3f26e56368c62b0900a10d83a11b7783630963b) [@bep](https://github.com/bep) 
* Update to Go 1.10 (!) [ff10c15a](https://github.com/gohugoio/hugo/commit/ff10c15a93632043f7a7f6551a30487c9ef58c50) [@bep](https://github.com/bep) 
* Add WARNING for integer YAML keys [0816a97a](https://github.com/gohugoio/hugo/commit/0816a97a469f11d8e9706143975eaa532e29639b) [@bep](https://github.com/bep) [#4393](https://github.com/gohugoio/hugo/issues/4393)
* Tune stringifyMapKeys [10a917df](https://github.com/gohugoio/hugo/commit/10a917dfdce8851666c5b89ebc02af6f6c84ab59) [@bep](https://github.com/bep) 
* Rename stringifyYAMLMapKeys to stringifyMapKeys [d4beef0d](https://github.com/gohugoio/hugo/commit/d4beef0d2bb8f6481fa80e1d938454a7d4e38814) [@bep](https://github.com/bep) 
* Add benchmarks for stringifyYAMLMapKeys [51213e0b](https://github.com/gohugoio/hugo/commit/51213e0be19fc19dbca9815afa95c73bd6d159c2) [@bep](https://github.com/bep) 
* Add support for `YAML` array data files [1fa24177](https://github.com/gohugoio/hugo/commit/1fa2417777d82b81bf37919ad02de4f5dcbf0d50) [@vassudanagunta](https://github.com/vassudanagunta) [#3890](https://github.com/gohugoio/hugo/issues/3890)
* Account for array type data in data dir merge/override logic [bb549a0d](https://github.com/gohugoio/hugo/commit/bb549a0d57505a6b8f28930bb91a9ab44cbb3288) [@vassudanagunta](https://github.com/vassudanagunta) [#4366](https://github.com/gohugoio/hugo/issues/4366)
* Add "target" and "rel" parameters to figure shortcode [2e95ec68](https://github.com/gohugoio/hugo/commit/2e95ec6844bf65a25485bdc8e2638e45788f2dcf) [@kaushalmodi](https://github.com/kaushalmodi) 
* image" property, not "twitter:image:src" [76d38d5e](https://github.com/gohugoio/hugo/commit/76d38d5e5322fc6220fb9e74f9ca0668606ebb5d) [@kaushalmodi](https://github.com/kaushalmodi) 

## Fixes

### Core

* Continue `GitInfo` lookup on error [e9750d83](https://github.com/gohugoio/hugo/commit/e9750d831f749afa928d8a099af5889d18cb2484) [@bep](https://github.com/bep) 
* Fix paginator `URL` for sections with URL in front matter [9f740b37](https://github.com/gohugoio/hugo/commit/9f740b37cfb3278e34a5d085380ccd4d619dabff) [@bep](https://github.com/bep) [#4415](https://github.com/gohugoio/hugo/issues/4415)

### Other

* Fix bug in Site.assembleSections method [00868081](https://github.com/gohugoio/hugo/commit/00868081f624928d773a7b698654766f8cd70069) [@vassudanagunta](https://github.com/vassudanagunta) [#4447](https://github.com/gohugoio/hugo/issues/4447)
* Update Blackfriday to fix footnote HTML5 validation error [492fea7c](https://github.com/gohugoio/hugo/commit/492fea7cd2bfcbdfe9f56aa0ae659cf62648833b) [@bep](https://github.com/bep) [#4433](https://github.com/gohugoio/hugo/issues/4433)
* Fix `YAML` maps key type [16a5c745](https://github.com/gohugoio/hugo/commit/16a5c74519771138023f019fe535fa5b250dc50d) [@dmgawel](https://github.com/dmgawel) [#2441](https://github.com/gohugoio/hugo/issues/2441)
* Remove `ERROR` on missing baseURL [55bd46a6](https://github.com/gohugoio/hugo/commit/55bd46a633d68f62e131457631ba839d6f876a55) [@bep](https://github.com/bep) [#4397](https://github.com/gohugoio/hugo/issues/4397)





