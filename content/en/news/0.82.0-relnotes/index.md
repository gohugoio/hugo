
---
date: 2021-03-21
title: "Hugo 0.82: Mostly bugfixes"
description: "Mostly bug fixes, but some useful improvements with Markdown attributes."
categories: ["Releases"]
---

This is a small release, mostly a maintainance/bugfix release. But also notable is that you can now add custom Markdown attributes (e.g. CSS classes) to code fences ([aed7df62](https://github.com/gohugoio/hugo/commit/aed7df62a811b07b73ec5cbbf03e69e4bbf00919) [@bep](https://github.com/bep) [#8278](https://github.com/gohugoio/hugo/issues/8278)) and that you can use the attribute lists in title render hooks (`.Attributes`; see [cd0c5d7e](https://github.com/gohugoio/hugo/commit/cd0c5d7ef32cbd570af00c50ce760452381df64e) [@bep](https://github.com/bep) [#8270](https://github.com/gohugoio/hugo/issues/8270)).

This release represents **28 contributions by 8 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), and [@gzagatti](https://github.com/gzagatti) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **20 contributions by 12 contributors**. A special thanks to [@bep](https://github.com/bep), [@jmooring](https://github.com/jmooring), [@rootkea](https://github.com/rootkea), and [@PaulPineda](https://github.com/PaulPineda) for their work on the documentation site.


Hugo now has:

* 50763+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 435+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 370+ [themes](http://themes.gohugo.io/)

## Notes

* We have made `.MediaType` comparable [ba1d0051](https://github.com/gohugoio/hugo/commit/ba1d0051b44fdd242b20899e195e37ab26501516) [@bep](https://github.com/bep) [#8317](https://github.com/gohugoio/hugo/issues/8317)[#8324](https://github.com/gohugoio/hugo/issues/8324). This also means that the old `MediaType.Suffix` and `MediaType.FullSuffix` is moved to `MediaType.FirstSuffix.Suffix` and `MediaType.FirstSuffix.FullSuffix`, which also better describes what they represent.

## Enhancements

### Templates

* Add method mappings for strings.Contains, strings.ContainsAny [7f853003](https://github.com/gohugoio/hugo/commit/7f8530039aa018f23bad9d58e97777705a6d19ac) [@bep](https://github.com/bep) 

### Output

* Make Type comparable [ba1d0051](https://github.com/gohugoio/hugo/commit/ba1d0051b44fdd242b20899e195e37ab26501516) [@bep](https://github.com/bep) [#8317](https://github.com/gohugoio/hugo/issues/8317)[#8324](https://github.com/gohugoio/hugo/issues/8324)
* Add a basic benchmark [4d24e2a3](https://github.com/gohugoio/hugo/commit/4d24e2a3261d8c7dc0395db3ac4de89ebb0974a5) [@bep](https://github.com/bep) 

### Other

* Regenerate docs helper [86b4fd35](https://github.com/gohugoio/hugo/commit/86b4fd35e78f545510f19b49246a3ccf5487831b) [@bep](https://github.com/bep) 
* Regen CLI docs [195d108d](https://github.com/gohugoio/hugo/commit/195d108da75c9e5b9ef790bc4a5879c1e913964b) [@bep](https://github.com/bep) 
* Simplify some config loading code [df8bb881](https://github.com/gohugoio/hugo/commit/df8bb8812f466bce563cdba297db3cd3f954a799) [@bep](https://github.com/bep) 
* Update github.com/evanw/esbuild v0.9.0 => v0.9.6 [57d8d208](https://github.com/gohugoio/hugo/commit/57d8d208ed2245858c6439f19803bf2749f9377f) [@bep](https://github.com/bep) 
* Apply OS env overrides twice [fc06e850](https://github.com/gohugoio/hugo/commit/fc06e85082b63a54d9403e57c8d01a7d5a62fc04) [@bep](https://github.com/bep) 
* Attributes for code fences should be placed after the lang indicator only [b725253f](https://github.com/gohugoio/hugo/commit/b725253f9e3033e18bd45096c0622e6fb7b1ff79) [@bep](https://github.com/bep) [#8313](https://github.com/gohugoio/hugo/issues/8313)
* Bump github.com/tdewolff/minify/v2 v2.9.15 [35dedf15](https://github.com/gohugoio/hugo/commit/35dedf15c04a605df4d4a09263b0b299e5161f86) [@bep](https://github.com/bep) [#8332](https://github.com/gohugoio/hugo/issues/8332)
* More explicit support link to discourse [137d2dab](https://github.com/gohugoio/hugo/commit/137d2dab3285e9b0f8fe4dcc65ab6ecf8bb09002) [@davidsneighbour](https://github.com/davidsneighbour) 
* Update to esbuild v0.9.0 [1b1dcf58](https://github.com/gohugoio/hugo/commit/1b1dcf586e220c3a8ad5ecfa8e4c3dac97f0ab44) [@bep](https://github.com/bep) 
* Allow more spacing characters in strings [0a2ab3f8](https://github.com/gohugoio/hugo/commit/0a2ab3f8feb961f8394b1f9964fab36bfa468027) [@moorereason](https://github.com/moorereason) [#8079](https://github.com/gohugoio/hugo/issues/8079)[#8079](https://github.com/gohugoio/hugo/issues/8079)
* Rename a test [35bfb662](https://github.com/gohugoio/hugo/commit/35bfb662229226d5f3cc3077ca74323f0aa88b7d) [@bep](https://github.com/bep) 
* Add a debug helper [6d21559f](https://github.com/gohugoio/hugo/commit/6d21559fb55cda39c7b92bb61fd8e65a84465fe5) [@bep](https://github.com/bep) 
* Add support for Google Analytics v4 [ba16a14c](https://github.com/gohugoio/hugo/commit/ba16a14c6e884e309380610331aff78213f84751) [@djatwood](https://github.com/djatwood) 
* Bump go.mod to Go 1.16 [782c79ae](https://github.com/gohugoio/hugo/commit/782c79ae61a5ec30746ce3729933d6b4d31e0540) [@bep](https://github.com/bep) [#8294](https://github.com/gohugoio/hugo/issues/8294)
* #8210 Upgrade golang version for Dockerfile [5afcae7e](https://github.com/gohugoio/hugo/commit/5afcae7e0b4c08bc37db6e34ab4cf960558f4b6e) [@systemkern](https://github.com/systemkern) 
* Update CONTRIBUTING.md [60469f42](https://github.com/gohugoio/hugo/commit/60469f429e227631d76d951f2ed92986f0bd92e9) [@bep](https://github.com/bep) 
* Handle attribute lists in code fences [aed7df62](https://github.com/gohugoio/hugo/commit/aed7df62a811b07b73ec5cbbf03e69e4bbf00919) [@bep](https://github.com/bep) [#8278](https://github.com/gohugoio/hugo/issues/8278)
* Allow markdown attribute lists to be used in title render hooks [cd0c5d7e](https://github.com/gohugoio/hugo/commit/cd0c5d7ef32cbd570af00c50ce760452381df64e) [@bep](https://github.com/bep) [#8270](https://github.com/gohugoio/hugo/issues/8270)
* bump github.com/kyokomi/emoji/v2 from 2.2.7 to 2.2.8 [88a85dce](https://github.com/gohugoio/hugo/commit/88a85dcea951b0b5622cf02b167ec9299d93118b) [@dependabot[bot]](https://github.com/apps/dependabot) 

## Fixes

### Output

* Fix output format handling for render hooks [18074d0c](https://github.com/gohugoio/hugo/commit/18074d0c2375cc4bf4d7933dd4206cb878a23d1c) [@bep](https://github.com/bep) [#8176](https://github.com/gohugoio/hugo/issues/8176)

### Other

* Fix OS env override for nested config param only available in theme [7ed56c69](https://github.com/gohugoio/hugo/commit/7ed56c6941edfdfa42eef2b779020b5d46ca194a) [@bep](https://github.com/bep) [#8346](https://github.com/gohugoio/hugo/issues/8346)
* Fix `new theme` command description [24c716ca](https://github.com/gohugoio/hugo/commit/24c716cac35b0c5476944108e545058749c43e61) [@rootkea](https://github.com/rootkea) 
* Fix handling of utf8 runes in nullString() [f6612d8b](https://github.com/gohugoio/hugo/commit/f6612d8bd8c4c3bb498178d14f45d3acdf86aa7c) [@moorereason](https://github.com/moorereason) 
* Fixes #7698. [01dd7c16](https://github.com/gohugoio/hugo/commit/01dd7c16af6204d18d530f9d3018689215482170) [@gzagatti](https://github.com/gzagatti) 
* Fix autocomplete docs [c8f45d1d](https://github.com/gohugoio/hugo/commit/c8f45d1d861f596821afc068bd12eb1213aba5ce) [@bep](https://github.com/bep) 





