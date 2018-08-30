---
date: 2017-06-16T17:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.23 is mainly a release that handles all the small changes needed to get Hugo moved to a GitHub organisation"
link: ""
title: "Hugo 0.23"
draft: false
author: bep
aliases: [/0-23/]
---

Hugo `0.23` is mainly a release that handles all the small changes needed to get Hugo moved to a GitHub organisation: [gohugoio](https://github.com/gohugoio), but it also contains a couple of important fixes that makes this an update worth-while for all.

Hugo now has:

* 17739&#43; [stars](https://github.com/gohugoio/hugo/stargazers)
* 494&#43; [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 165&#43; [themes](http://themes.gohugo.io/)

## Fixes

* Fix handling of duplicate footnotes [a9e551a1](https://github.com/gohugoio/hugo/commit/a9e551a100e60a603210ee083103dd73369d6a98) [@bep](https://github.com/bep) [#1912](https://github.com/gohugoio/hugo/issues/1912) 
*  Add support for spaces in project folder for `GitInfo` #3533 #3552

## GitHub organisation related changes

* Update layout references to gohugoio/hugo [66d4850b](https://github.com/gohugoio/hugo/commit/66d4850b89db293dc58e828de784037f06c6c8dc) [@bep](https://github.com/bep) 
* Update content references to gohugoio/hugo [715ff1f8](https://github.com/gohugoio/hugo/commit/715ff1f87406edf27738c8c0f52fe185fa974ee8) [@bep](https://github.com/bep) 
* Add note on updates for rpm-based distros [52a0cea6](https://github.com/gohugoio/hugo/commit/52a0cea65de7b75ae1662abe3dec36fca3604617) [@daftaupe](https://github.com/daftaupe) 
* Update logo link in README [ccb8300d](https://github.com/gohugoio/hugo/commit/ccb8300d380636d75a39f4133284eb0109e836c3) [@bep](https://github.com/bep) 
* Remove docs building from CI builds [214dbdfb](https://github.com/gohugoio/hugo/commit/214dbdfb6f016d21415bc1ed511a37a084238878) [@bep](https://github.com/bep) 
* Adjust docs path [729be807](https://github.com/gohugoio/hugo/commit/729be8074bddb58c9111f32c55cc769e49cd0d5a) [@bep](https://github.com/bep) 
* Add docs as submodule [6cee0dfe](https://github.com/gohugoio/hugo/commit/6cee0dfe53899d433afc3c173a87d56265904cb0) [@bep](https://github.com/bep) 
* Update Gitter link in README [fbb25014](https://github.com/gohugoio/hugo/commit/fbb25014e1306ce7127d53e5fc4fc49867790336) [@bep](https://github.com/bep) 
* Change Windows build badge link, take #3 [86543d6a](https://github.com/gohugoio/hugo/commit/86543d6a50251b40540ebd0b851d45eb99d017c7) [@bep](https://github.com/bep) 
* Update Windows build link [e6ae32a0](https://github.com/gohugoio/hugo/commit/e6ae32a0ba75b9894418227e87391defbb1b3b49) [@bep](https://github.com/bep) 
* Update links in CONTRIBUTING.md due to the org transition [95386544](https://github.com/gohugoio/hugo/commit/95386544e858949a2baa414f395f30aaf66a6257) [@digitalcraftsman](https://github.com/digitalcraftsman) 
* Update source path in Dockerfile due to the org transition [7b99fb9f](https://github.com/gohugoio/hugo/commit/7b99fb9f1ca8381457afe9d8e953a388b8ada182) [@digitalcraftsman](https://github.com/digitalcraftsman) 
* Update clone folder in appveyor.yml due to the org transition [d531d17b](https://github.com/gohugoio/hugo/commit/d531d17b3be0b14faf4934611e01ac3289e37835) [@digitalcraftsman](https://github.com/digitalcraftsman) 
* Update import path in snapcraft.yaml due to the org transition [9266bf9d](https://github.com/gohugoio/hugo/commit/9266bf9d4c24592b875a7f6b92f761b4cea40879) [@digitalcraftsman](https://github.com/digitalcraftsman) 
* Run gofmt to get imports in line vs gohugoio/hugo [873a6f18](https://github.com/gohugoio/hugo/commit/873a6f18851bcda79d562ff6c02e1109e8e31a88) [@bep](https://github.com/bep) 
* Update Makefile vs gohugoio/hugo [f503d76a](https://github.com/gohugoio/hugo/commit/f503d76a3b2719bbb65ab9df5595d0dbc871fae9) [@bep](https://github.com/bep) 
* Update README to point to gohugoio/hugo [93643860](https://github.com/gohugoio/hugo/commit/93643860c9db10c6c32176b17cc83f1c317279bd) [@bep](https://github.com/bep) 
* Update examples to point to gohugoio/hugo [db46bcf8](https://github.com/gohugoio/hugo/commit/db46bcf82d060656d4bc731550e63ec9cf8576f2) [@bep](https://github.com/bep) 
* Update textual references in Go source to point to gohugoio/hugo [c17ad675](https://github.com/gohugoio/hugo/commit/c17ad675e8fcdb2db40fc50816b8f016bc14294c) [@bep](https://github.com/bep) 
* Update import paths to gohugoio/hugo [d8717cd4](https://github.com/gohugoio/hugo/commit/d8717cd4c74e80ea8e20adead9321412a2d76022) [@bep](https://github.com/bep) 