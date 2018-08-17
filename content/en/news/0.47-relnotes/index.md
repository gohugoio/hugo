
---
date: 2018-08-17
title: "Output Minification, Live-Reload Fixes and More"
description: "Hugo 0.47: Adds minification of rendered output, but is mostly a massive bug fix release."
categories: ["Releases"]
---

Hugo `0.47` is named **Hugo Reloaded**. It adds minification support for the final rendered output (run `hugo --minify`), but it is mostly a bug fix release. And most notably, it fixes a set of issues with live-reloading/partial rebuilds when running `hugo server`. Working with bundles should now be a more pleasant experience, to pick one example.

This release represents **35 contributions by 6 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@satotake](https://github.com/satotake), [@anthonyfok](https://github.com/anthonyfok), and [@coliff](https://github.com/coliff) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday)  for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **21 contributions by 10 contributors**. A special thanks to [@bep](https://github.com/bep), [@aapeliv](https://github.com/aapeliv), [@regisphilibert](https://github.com/regisphilibert), and [@brentybh](https://github.com/brentybh) for their work on the documentation site.


Hugo now has:

* 27980+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 442+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 251+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Suppress blank lines from opengraph internal template [c09ee78f](https://github.com/gohugoio/hugo/commit/c09ee78fd235599d3fb794110cd75c024d80cfca) [@anthonyfok](https://github.com/anthonyfok) 
* Add MIME type to embedded JS [755d1ffe](https://github.com/gohugoio/hugo/commit/755d1ffe7a22d8ad83485240ff78cf25d501602f) [@bep](https://github.com/bep) [#5042](https://github.com/gohugoio/hugo/issues/5042)
* Add `os.Stat` template function [d7112085](https://github.com/gohugoio/hugo/commit/d71120852a8e14d0ea4d24de269fce041ef7b666) [@satotake](https://github.com/satotake) 

### Output

* Add support for minification of final output [789ef8c6](https://github.com/gohugoio/hugo/commit/789ef8c639e4621abd36da530bcb5942ac9297da) [@bep](https://github.com/bep) [#1251](https://github.com/gohugoio/hugo/issues/1251)

### Other

* Regenerate CLI docs [4a16b5f4](https://github.com/gohugoio/hugo/commit/4a16b5f4b0adbb31fee611c378de9d5526de2f86) [@bep](https://github.com/bep) 
* Include theme name in version mismatch error [e5052f4e](https://github.com/gohugoio/hugo/commit/e5052f4e09b6df590cddf2f8bc2c834fd3af3082) [@bep](https://github.com/bep) [#5044](https://github.com/gohugoio/hugo/issues/5044)
* Make the JS minifier matcher less specific [c81fbf46](https://github.com/gohugoio/hugo/commit/c81fbf4625ae7cc7dd3a7a526331ddfdf5237cc6) [@bep](https://github.com/bep) [#5073](https://github.com/gohugoio/hugo/issues/5073)
* Close file when done [f6ae436c](https://github.com/gohugoio/hugo/commit/f6ae436c5878bafeafa0bb2646a2c9b32c9b4380) [@bep](https://github.com/bep) [#5062](https://github.com/gohugoio/hugo/issues/5062)
* https links to 3rd party sites [c2a67413](https://github.com/gohugoio/hugo/commit/c2a6741394bc609a663522b245d3d75f0ad17da4) [@coliff](https://github.com/coliff) 
* Update alias.go [06bd0136](https://github.com/gohugoio/hugo/commit/06bd0136419ebd6727635716c7023494cc5a8257) [@coliff](https://github.com/coliff) 
* Remove test debug [fb3cb05c](https://github.com/gohugoio/hugo/commit/fb3cb05cc3dfc50370449f622fb0130ba7e0ced2) [@bep](https://github.com/bep) 
* Update dependencies [d07882df](https://github.com/gohugoio/hugo/commit/d07882dfb76a65cce79aaa6f27df71279cd30600) [@bep](https://github.com/bep) 
* Update Chroma [7f535671](https://github.com/gohugoio/hugo/commit/7f5356717d14079432365974e1424fc4ff5987c9) [@bep](https://github.com/bep) [#5025](https://github.com/gohugoio/hugo/issues/5025)
* Remove alias of os.Stat [71931b30](https://github.com/gohugoio/hugo/commit/71931b30b1813b146aaa60f5cdab16c0f9ebebdb) [@satotake](https://github.com/satotake) 
* Renmae FileStat Stat [d40116e5](https://github.com/gohugoio/hugo/commit/d40116e5f941e4734ed3bed69dce8ffe11fc76b2) [@satotake](https://github.com/satotake) 
* Reduce allocation in the benchmark itself [a6b1eb1e](https://github.com/gohugoio/hugo/commit/a6b1eb1e9150aa5c1c86fe7424cc4167d6f59a5a) [@bep](https://github.com/bep) 
* Simplify the 0 transformer case [27110133](https://github.com/gohugoio/hugo/commit/27110133ffca05feae2e11a9ff28a9a00f613350) [@bep](https://github.com/bep) 

## Fixes

### Templates

* Fix compiling Amber templates that import other templates [37438757](https://github.com/gohugoio/hugo/commit/37438757788d279c839506d54f460b2ab37db164) [@Stebalien](https://github.com/Stebalien) 
* Reimplement the ".Params tolower" template transformer [5c538491](https://github.com/gohugoio/hugo/commit/5c5384916e8f954f3ea66148ecceb3732584588e) [@bep](https://github.com/bep) [#5068](https://github.com/gohugoio/hugo/issues/5068)

### Output

* Fix Resource output in multihost setups [78f8475a](https://github.com/gohugoio/hugo/commit/78f8475a054a6277d37f13329afd240b00dc9408) [@bep](https://github.com/bep) [#5058](https://github.com/gohugoio/hugo/issues/5058)

### Core

* Force render of any changed page, even in Fast Render Mode [22475460](https://github.com/gohugoio/hugo/commit/2247546017c00201d2ce1232dd5303295451f1cc) [@bep](https://github.com/bep) [#5083](https://github.com/gohugoio/hugo/issues/5083)
* Add configFile(s) back to the watch list after REMOVE event [abc54080](https://github.com/gohugoio/hugo/commit/abc54080ec8c43e8989c081d934b59f0c9570c0b) [@anthonyfok](https://github.com/anthonyfok) [#4701](https://github.com/gohugoio/hugo/issues/4701)
* Gracefully handle typos in server config when running the server [a655e00d](https://github.com/gohugoio/hugo/commit/a655e00d702dbc20b3961b131b33ab21841b043d) [@bep](https://github.com/bep) [#5081](https://github.com/gohugoio/hugo/issues/5081)
* Fix shortcode output wrapped in p [78c99463](https://github.com/gohugoio/hugo/commit/78c99463fdd45c91af9933528d12d36a86dc6482) [@gllera](https://github.com/gllera) [#1642](https://github.com/gohugoio/hugo/issues/1642)
* Adjust tests for shortcode p-issue [baa62d0a](https://github.com/gohugoio/hugo/commit/baa62d0abbbf24a17d0aa800a4bb217f026c49ad) [@bep](https://github.com/bep) [#1642](https://github.com/gohugoio/hugo/issues/1642)
* Fix image cache-clearing for sub-languages [9d973004](https://github.com/gohugoio/hugo/commit/9d973004f5379cff2adda489566fe40683553c4c) [@bep](https://github.com/bep) [#5084](https://github.com/gohugoio/hugo/issues/5084)
* Fix error when deleting a bundle in server mode [0a88741f](https://github.com/gohugoio/hugo/commit/0a88741fe85f4f7aedc02ed748dfeb8ccc073dbf) [@bep](https://github.com/bep) [#5077](https://github.com/gohugoio/hugo/issues/5077)
* Fix Related when called from shortcode [0dd06bda](https://github.com/gohugoio/hugo/commit/0dd06bdac008aa81ec2e8f29ad8110dac0227011) [@bep](https://github.com/bep) [#5071](https://github.com/gohugoio/hugo/issues/5071)
* Use the interface value when doing Related search [a6f199f7](https://github.com/gohugoio/hugo/commit/a6f199f7a640161333608b4a843d701f7e182829) [@bep](https://github.com/bep) [#5071](https://github.com/gohugoio/hugo/issues/5071)
* Fix GitInfo when multiple content dirs [2182ecfd](https://github.com/gohugoio/hugo/commit/2182ecfd34a24521bf0e3c939627a55327eb1e19) [@bep](https://github.com/bep) [#5054](https://github.com/gohugoio/hugo/issues/5054)
* Add multiple content dirs to GitInfo test site [e85833d8](https://github.com/gohugoio/hugo/commit/e85833d868a902840c5ed1c90713256153b2548b) [@bep](https://github.com/bep) [#5054](https://github.com/gohugoio/hugo/issues/5054)
* Fix "adding a bundle" in server mode [d139a037](https://github.com/gohugoio/hugo/commit/d139a037d98e4b388687eecb7831758412247c58) [@bep](https://github.com/bep) [#5075](https://github.com/gohugoio/hugo/issues/5075)
* Fix typo [c362634b](https://github.com/gohugoio/hugo/commit/c362634b7d8802ea81b0b4341c800a9f78f7cd7c) [@satotake](https://github.com/satotake) 
