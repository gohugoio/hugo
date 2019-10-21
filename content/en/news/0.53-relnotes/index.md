
---
date: 2018-12-24
title: "Hugo Christmas Edition"
description: "Hugo 0.53: Faster, config dir support, new unmarshal func, global site var, and more ..."
categories: ["Releases"]
---

From all of us to all of you, a very Merry Christmas -- and Hugo `0.53`!

The main new features in this release are:

**Config Dir:** You can now split your configuration sections into directories per environment. Hugo did support multiple configuration files before this release, but it was hard to manage for bigger sites, especially those with multiple languages. With this we have also formalized the concept of an `environment`; the defaults are `production` (when running `hugo`) or `development` (when running `hugo server`) but you can create any environment you like.  We will update the documentation, but all the details are in [this issue](https://github.com/gohugoio/hugo/pull/5501#issue-236237630). Also, see [this PR](https://github.com/gohugoio/hugoDocs/pull/683) for how the refactored configuration for the Hugo website looks like.

**Unmarshal JSON, TOML, YAML or CSV:** `transform.Unmarshal` (see the [documentation](https://gohugo.io//functions/transform.unmarshal/) is a new and powerful template function that can turn `Resource` objects or strings with JSON, TOML, YAML or CSV into maps/arrays.

**Global site and hugo var:** Two new global variables in `site` and `hugo`. `hugo` gives you version info etc. (`{{ hugo.Version }}`, `{{ hugo.Environment }}`), but the `site` is probably more useful, as it allows you to access the current [site's variables](https://gohugo.io/variables/site/) (e.g. `{{ site.RegularPages }}`) without any context (or ".").

This version is also the fastest to date. A site building benchmark shows around 10% faster, but that depends on the site. The important part here is that we're not getting slower.  It’s quite a challenge to consistently add significant new functionality and simultaneously improve performance. It's like not gaining weight during Christmas. We also had a small performance boost in version `0.50`. A user then reported that his big and complicated site had a 30% reduction in build time. This is important to us, one of the core features. It's in the slogan: "The world’s fastest framework for building websites."

This release represents **37 contributions by 5 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@coliff](https://github.com/coliff), and [@jfyuen](https://github.com/jfyuen) for their ongoing contributions. And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **19 contributions by 8 contributors**. A special thanks to [@bep](https://github.com/bep), [@kaushalmodi](https://github.com/kaushalmodi), [@peaceiris](https://github.com/peaceiris), and [@moorereason](https://github.com/moorereason) for their work on the documentation site.


Hugo now has:

* 31174+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 279+ [themes](http://themes.gohugo.io/)


## Notes

* The `hugo benchmark` command is removed
* We now do not publish transformed inline resources, e.g. minified CSS only accessed via `.Content.`, e.g. `{{ ($css | minify).Content }}`. Before this version, the minified CSS in that example would be copied to `/public`, which was never the intention. If you want that, you need to access either `.RelPermalink` or `.Permalink`.


## Enhancements

### Templates

* Include options in cache key [be58c7b9](https://github.com/gohugoio/hugo/commit/be58c7b9c88116094ca2b424c77210ddcccfff8e) [@bep](https://github.com/bep) [#5555](https://github.com/gohugoio/hugo/issues/5555)
* Simplify transform.Unmarshal func [094709e1](https://github.com/gohugoio/hugo/commit/094709e105d48547bf5297adc0ad0c777678b0a6) [@bep](https://github.com/bep) [#5428](https://github.com/gohugoio/hugo/issues/5428)
* Add transform.Unmarshal func [822dc627](https://github.com/gohugoio/hugo/commit/822dc627a1cfdf1f97882f27761675ac6ace7669) [@bep](https://github.com/bep) [#5428](https://github.com/gohugoio/hugo/issues/5428)
* Remove "double layout" lookup [d5a0b6bb](https://github.com/gohugoio/hugo/commit/d5a0b6bbbc83a3e274c62ed397a293f04ee8d241) [@bep](https://github.com/bep) [#5390](https://github.com/gohugoio/hugo/issues/5390)
* Add reflect namespace [c84f506f](https://github.com/gohugoio/hugo/commit/c84f506f8ef1f2ca94ab96718a22ba6e290235ac) [@moorereason](https://github.com/moorereason) [#4081](https://github.com/gohugoio/hugo/issues/4081)
* Use the correct Hugo var [931a1324](https://github.com/gohugoio/hugo/commit/931a1324503a4414e38d26efe82e1add811a8d29) [@bep](https://github.com/bep) [#5467](https://github.com/gohugoio/hugo/issues/5467)
* Add tpl/site and tpl/hugo [831d23cb](https://github.com/gohugoio/hugo/commit/831d23cb4d1ca99cdc15ed31c8ee1f981497be8f) [@bep](https://github.com/bep) [#5470](https://github.com/gohugoio/hugo/issues/5470)[#5467](https://github.com/gohugoio/hugo/issues/5467)[#5503](https://github.com/gohugoio/hugo/issues/5503)
* Add godoc packages comments [30a7c9ea](https://github.com/gohugoio/hugo/commit/30a7c9ea37a0f36451946f8688a3f807618a7eff) [@moorereason](https://github.com/moorereason) 

### Core

* Adjust test [25ddbb09](https://github.com/gohugoio/hugo/commit/25ddbb09fea7794edbbafa2ffce4e361cdc9bacf) [@bep](https://github.com/bep) [#5544](https://github.com/gohugoio/hugo/issues/5544)
* Add .Name as a shortcode variable [10217144](https://github.com/gohugoio/hugo/commit/1021714449a05ef85b2fdfaf65b354cbdee44f23) [@bep](https://github.com/bep) [#5546](https://github.com/gohugoio/hugo/issues/5546)
* Improve logic of output path trimming [0483299b](https://github.com/gohugoio/hugo/commit/0483299bc06a742d40528e0d675e42e149910853) [@moorereason](https://github.com/moorereason) [#4666](https://github.com/gohugoio/hugo/issues/4666)
* Enable Emoji in site benchmark [4d93aca2](https://github.com/gohugoio/hugo/commit/4d93aca27dfdebc9e06948ccf37a7922dac09d65) [@bep](https://github.com/bep) 
* Restore taxonomy term path separation [9ce0a1fb](https://github.com/gohugoio/hugo/commit/9ce0a1fb7011bd75eb0e2262e35354c49ce98ac5) [@bep](https://github.com/bep) [#5513](https://github.com/gohugoio/hugo/issues/5513)
* Add .Site.Sites [83783588](https://github.com/gohugoio/hugo/commit/8378358857d852458d01c667d59d13baa59a719c) [@bep](https://github.com/bep) [#5504](https://github.com/gohugoio/hugo/issues/5504)

### Other

* Adjust CSV example [62d031ae](https://github.com/gohugoio/hugo/commit/62d031aedfc128729b460241bd99d97b5d902e62) [@bep](https://github.com/bep) [#5555](https://github.com/gohugoio/hugo/issues/5555)
* Rename CSV option from comma to delimiter [ce06bdb1](https://github.com/gohugoio/hugo/commit/ce06bdb16a64dd39a8ebbb2e5a53b33520b00bb1) [@bep](https://github.com/bep) [#5555](https://github.com/gohugoio/hugo/issues/5555)
* Document transform.Unmarshal [2efc1a64](https://github.com/gohugoio/hugo/commit/2efc1a64c391420b1007f6e94b6ff616fb136635) [@bep](https://github.com/bep) [#5556](https://github.com/gohugoio/hugo/issues/5556)
* Regenerate CLI docs [e691c48a](https://github.com/gohugoio/hugo/commit/e691c48a5a9b4db5aa5383de6b83352fc18cc633) [@bep](https://github.com/bep) [#5544](https://github.com/gohugoio/hugo/issues/5544)
* Add CSV support to transform.Unmarshal [a5744697](https://github.com/gohugoio/hugo/commit/a5744697971d296eb973e04e4259fe9e516b908f) [@bep](https://github.com/bep) [#5555](https://github.com/gohugoio/hugo/issues/5555)
* Prevent resource publishing for transformed inline resources [43f9df01](https://github.com/gohugoio/hugo/commit/43f9df0194d229805d80b13c9e38a7a0fec12cf4) [@bep](https://github.com/bep) [#4944](https://github.com/gohugoio/hugo/issues/4944)
* Remove the benchmark command [35bfca3b](https://github.com/gohugoio/hugo/commit/35bfca3b14977eaebab4003b43b5236c1888d93d) [@bep](https://github.com/bep) [#5543](https://github.com/gohugoio/hugo/issues/5543)
* Move the emoji parsing to pageparser [9cd54cab](https://github.com/gohugoio/hugo/commit/9cd54cab20a03475e34ca462bd943069111481ae) [@bep](https://github.com/bep) [#5534](https://github.com/gohugoio/hugo/issues/5534)
* Split the page lexer into some more files [a8853f1c](https://github.com/gohugoio/hugo/commit/a8853f1c5ace30ae8d256ad374bdb280c95d4228) [@bep](https://github.com/bep) [#5534](https://github.com/gohugoio/hugo/issues/5534)
* parser/pageparser: Add a benchmark [f2167de8](https://github.com/gohugoio/hugo/commit/f2167de83493f13f02dd622425364668834f8208) [@bep](https://github.com/bep) 
* Update to Go 1.11.4 [bb9c2988](https://github.com/gohugoio/hugo/commit/bb9c2988f871ca5fe6af9c8e207ec852c631c3b3) [@bep](https://github.com/bep) [#5524](https://github.com/gohugoio/hugo/issues/5524)
* Simplify implementation [f7691fe9](https://github.com/gohugoio/hugo/commit/f7691fe9652aa12b6c582dea0ae2555e772d1a5f) [@bep](https://github.com/bep) 
* Support unquoted URLs in canonifyURLs replacer [efe0b4e5](https://github.com/gohugoio/hugo/commit/efe0b4e5c0292c1e5e27b0c32fbc368062fde3e8) [@bep](https://github.com/bep) [#5529](https://github.com/gohugoio/hugo/issues/5529)
* Regenerate CLI docs [50686817](https://github.com/gohugoio/hugo/commit/50686817072c8bef947959cb2bcc7f1914c7f839) [@bep](https://github.com/bep) [#5507](https://github.com/gohugoio/hugo/issues/5507)
* Add /config dir support [78294740](https://github.com/gohugoio/hugo/commit/7829474088f835251f04caa1121d47e35fe89f7e) [@bep](https://github.com/bep) [#5422](https://github.com/gohugoio/hugo/issues/5422)
* cache/filecache: Simplify test [514e18dc](https://github.com/gohugoio/hugo/commit/514e18dc27ce37a0e9a231741d616cf29d50d610) [@bep](https://github.com/bep) [#5497](https://github.com/gohugoio/hugo/issues/5497)
* Use OS fs for test [b804a708](https://github.com/gohugoio/hugo/commit/b804a70881c7be26dc15274c4f98f1057469cbc1) [@bep](https://github.com/bep) [#5497](https://github.com/gohugoio/hugo/issues/5497)

## Fixes

### Templates

* Fix case handling in cast params [64b6b290](https://github.com/gohugoio/hugo/commit/64b6b290751df01c47ff8d8fe21a3eca7a5db283) [@bep](https://github.com/bep) [#5538](https://github.com/gohugoio/hugo/issues/5538)

### Other

* Fix "failed to create file caches from configuration: file exists" on Windows [5178cd13](https://github.com/gohugoio/hugo/commit/5178cd13a7da3c5f5ec5d3217c9e40fc0be7152a) [@bep](https://github.com/bep) [#5497](https://github.com/gohugoio/hugo/issues/5497)
* fix jekyll import highlight options [ab921476](https://github.com/gohugoio/hugo/commit/ab9214768de4ce10032d3fe7ec8c7b2932ead892) [@jfyuen](https://github.com/jfyuen) 
* Fix "always false" condition [25641891](https://github.com/gohugoio/hugo/commit/256418917c6642f7e5b3d3206ff4b6fa03b1cb28) [@Quasilyte](https://github.com/Quasilyte) 
* Fixx CSS2 color code handling [4b5f7439](https://github.com/gohugoio/hugo/commit/4b5f743959394d443c4dcaa0ccae21842b51adaf) [@bep](https://github.com/bep) [#5506](https://github.com/gohugoio/hugo/issues/5506)
* common/collections: Fix defines typo [83468481](https://github.com/gohugoio/hugo/commit/8346848109ab57cb04de87c6d86859c6b3de8ffa) [@coliff](https://github.com/coliff) 
