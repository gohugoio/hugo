
---
date: 2020-05-06
title: "JavaScript Transpiler"
description: "Hugo 0.70.0 adds a new pipe function that uses Babel to transpile JavaScript."
categories: ["Releases"]
---

This is a small release, and the main new feature is that you can now use [Babel](https://gohugo.io/hugo-pipes/babel/) to transpile JavaScript.

This release represents **22 contributions by 12 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@BurtonQin](https://github.com/BurtonQin), [@tekezo](https://github.com/tekezo), and [@sensimevanidus](https://github.com/sensimevanidus) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **6 contributions by 4 contributors**. A special thanks to [@bep](https://github.com/bep), [@MJ2097](https://github.com/MJ2097), [@jeremyzilar](https://github.com/jeremyzilar), and [@larryclaman](https://github.com/larryclaman) for their work on the documentation site.


Hugo now has:

* 43734+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 437+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 316+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Change defer RLock to RUnlock [5146dc61](https://github.com/gohugoio/hugo/commit/5146dc614fc45df698ebf890af06421dea988c96) [@BurtonQin](https://github.com/BurtonQin) 

### Output

* Modify gen chromastyles to output all CSS classes [102ec2da](https://github.com/gohugoio/hugo/commit/102ec2da7adcc4afb7050b17989f0486f8379679) [@acahir](https://github.com/acahir) [#7167](https://github.com/gohugoio/hugo/issues/7167)

### Core

* Add Unlock before panic [736f84b2](https://github.com/gohugoio/hugo/commit/736f84b2d539857f7fdd0e42353af80b4dccfe8d) [@BurtonQin](https://github.com/BurtonQin) 

### Other

* Update minify to v2.6.2 [01befcce](https://github.com/gohugoio/hugo/commit/01befcce35ec992d195ce1b9a6a1eeda693cb5a8) [@pperzyna](https://github.com/pperzyna) [#6699](https://github.com/gohugoio/hugo/issues/6699)
* Add support for sort by boolean [04b1a6d9](https://github.com/gohugoio/hugo/commit/04b1a6d997e72d9abada28db22650d38ccbcbb39) [@Mipsters](https://github.com/Mipsters) 
* Update to Libsass 3.6.4 [dd31e800](https://github.com/gohugoio/hugo/commit/dd31e800075eebd78f921df8b4865c238006e7a7) [@bep](https://github.com/bep) [#7231](https://github.com/gohugoio/hugo/issues/7231)
* Rename transpileJS to babel [6add6d77](https://github.com/gohugoio/hugo/commit/6add6d77b48cf0aab8b39d7a2bddedb1aa2a52b8) [@bep](https://github.com/bep) [#5764](https://github.com/gohugoio/hugo/issues/5764)
* Add JavaScript transpiling solution [2a171ff1](https://github.com/gohugoio/hugo/commit/2a171ff1c5d9b1603fe78c67d2d894bb2efccc8b) [@hmmmmniek](https://github.com/hmmmmniek) [#5764](https://github.com/gohugoio/hugo/issues/5764)
* Disable a test locally [67f92041](https://github.com/gohugoio/hugo/commit/67f920419a53c7ff11e01c4286dca23e92110a12) [@bep](https://github.com/bep) 
* Add diagnostic hints to init timeout message [fe60b7d9](https://github.com/gohugoio/hugo/commit/fe60b7d9e4c12dbc428f992c05969bc14c7fe7a2) [@mtlynch](https://github.com/mtlynch) 
* Update goldmark-highlighting [5c41f41a](https://github.com/gohugoio/hugo/commit/5c41f41ad4b14e48aea64687a7600f5ad231e879) [@satotake](https://github.com/satotake) [#7027](https://github.com/gohugoio/hugo/issues/7027)[#6596](https://github.com/gohugoio/hugo/issues/6596)
* Update go-org to v1.1.0 [2b28e5a9](https://github.com/gohugoio/hugo/commit/2b28e5a9cb79af2a8d70c80036f52bcf5399b9df) [@niklasfasching](https://github.com/niklasfasching) 
* Update to goldmark v1.1.28 [feaa582c](https://github.com/gohugoio/hugo/commit/feaa582cbe950e82969da5e99e3fb9a3947025df) [@bep](https://github.com/bep) [#7113](https://github.com/gohugoio/hugo/issues/7113)

## Fixes

### Other

* Fix some missing JS class collector cases [c03ea2b6](https://github.com/gohugoio/hugo/commit/c03ea2b66010d2996d652903cb8fa41e983e787f) [@bep](https://github.com/bep) [#7216](https://github.com/gohugoio/hugo/issues/7216)
* Fix IsAncestor and IsDescendant when the same page is passed [8d5766d4](https://github.com/gohugoio/hugo/commit/8d5766d417d6564a1aa1cbe8f9a29ab9bba22371) [@tekezo](https://github.com/tekezo) 
* Fix IsAncestor and IsDescendant under subsection [27a4c441](https://github.com/gohugoio/hugo/commit/27a4c4410cd9592249925fb14b32605fb961c597) [@tekezo](https://github.com/tekezo) 
* Fix typo in test suite [49e6c8cb](https://github.com/gohugoio/hugo/commit/49e6c8cb4ed83e20f1e0ac164e91c38854177b99) [@panakour](https://github.com/panakour) 
* Fix class collector when running with --minify [f37e77f2](https://github.com/gohugoio/hugo/commit/f37e77f2d338cf876cfa637a662acd76f0f2009b) [@bep](https://github.com/bep) [#7161](https://github.com/gohugoio/hugo/issues/7161)
* Fix toLower [27af5a33](https://github.com/gohugoio/hugo/commit/27af5a339a4d3c5712b5ed946a636a8c21916039) [@bep](https://github.com/bep) [#7198](https://github.com/gohugoio/hugo/issues/7198)
* Fix broken test [b3c82575](https://github.com/gohugoio/hugo/commit/b3c825756f3251f8b26e53262f9d6f484aecf750) [@bep](https://github.com/bep) 
* Fix typo in Hugo's Security Model [cd4d8202](https://github.com/gohugoio/hugo/commit/cd4d8202016bd3eb5ed9144c8945edaba73c8cf4) [@sensimevanidus](https://github.com/sensimevanidus) 
* Fix query parameter handling in server fast render mode [ee67dbef](https://github.com/gohugoio/hugo/commit/ee67dbeff5bae6941facaaa39cb995a1ee6def03) [@bep](https://github.com/bep) [#7163](https://github.com/gohugoio/hugo/issues/7163)





