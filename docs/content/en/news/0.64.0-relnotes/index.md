
---
date: 2020-02-04
title: "Mostly bugfixes …"
description: "Hugo 0.64.0 is the version you really really want …"
categories: ["Releases"]
---

Hugo **0.64.0** is mostly a bugfix-release, but well worth the download. The main reason this release comes so soon after the previous is my (me being [@bep](https://github.com/bep)) ongoing work on getting solid support for third-party libraries in [Hugo Modules](https://gohugo.io/hugo-modules/). In particular, this release makes the Hugo server's live-reload work with [Turbolinks](https://github.com/bep/hugo-alpine-test/blob/27927832630be588eab0be2197cc8c0cb5725540/config.toml#L11) and similar. Also worth mentioning is that `hugo mod get -u` (without any path) now correctly updates every module imported in `config.toml` even with Go 1.13.

This release represents **16 contributions by 2 contributors** to the main Hugo code base.
Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **6 contributions by 4 contributors**. A special thanks to [@bep](https://github.com/bep), [@peterkappus](https://github.com/peterkappus), [@kc0bfv](https://github.com/kc0bfv), and [@inwardmovement](https://github.com/inwardmovement) for their work on the documentation site.


Hugo now has:

* 41348+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 289+ [themes](http://themes.gohugo.io/)

## Enhancements

### Output

* Do not render alias paginator pages for non-HTML outputs [2d159e9c](https://github.com/gohugoio/hugo/commit/2d159e9cc7a25832e4b0cad226b149f7c4624708) [@bep](https://github.com/bep) [#6797](https://github.com/gohugoio/hugo/issues/6797)

### Other

* Mention a "no CGO rule" [29973101](https://github.com/gohugoio/hugo/commit/299731012441378bb9c057ceb0a3c277108aaf01) [@bep](https://github.com/bep) [#6842](https://github.com/gohugoio/hugo/issues/6842)
* Update to Go 1.13.7 and Go 1.12.16 [0792cfa9](https://github.com/gohugoio/hugo/commit/0792cfa9fae94a06a31e393a46fed3b1dd73b66a) [@bep](https://github.com/bep) [#6830](https://github.com/gohugoio/hugo/issues/6830)
* Add defer to livereload script tag [b3f0674b](https://github.com/gohugoio/hugo/commit/b3f0674b80a32425aeb4412f318c720391bbf773) [@bep](https://github.com/bep) 
* Don't use document.write to inject livereload [ef78a0d1](https://github.com/gohugoio/hugo/commit/ef78a0d18a13098bcea1ff2b2d45d7388b8d41a0) [@bep](https://github.com/bep) [#6507](https://github.com/gohugoio/hugo/issues/6507)
* Add a render hook whitespace test [58595864](https://github.com/gohugoio/hugo/commit/585958645372e6219239247dbac02e447d2b355b) [@bep](https://github.com/bep) [#6832](https://github.com/gohugoio/hugo/issues/6832)
* Inject livereload script right after head if possible [8f08cdd0](https://github.com/gohugoio/hugo/commit/8f08cdd0ac6a2decd5aa5c9c12c0b2c264f9a989) [@bep](https://github.com/bep) [#6821](https://github.com/gohugoio/hugo/issues/6821)
* Update goldmark to v1.1.22 [281abb18](https://github.com/gohugoio/hugo/commit/281abb18ee39fa2b5d4782b64f27cffcbf4e0240) [@bhavin192](https://github.com/bhavin192) 
* Make the build flags shared between sites [0df7bd62](https://github.com/gohugoio/hugo/commit/0df7bd62df460a49544845d5332f33b2020b48a1) [@bep](https://github.com/bep) [#6789](https://github.com/gohugoio/hugo/issues/6789)

## Fixes

### Other

* Fix module mount in sub folder [80dd6ddd](https://github.com/gohugoio/hugo/commit/80dd6ddde27ce36f5432fb780e94d4974b5277c7) [@bep](https://github.com/bep) [#6730](https://github.com/gohugoio/hugo/issues/6730)
* Fix config environment handling [2bbc865f](https://github.com/gohugoio/hugo/commit/2bbc865f7bb713b2d0d2dbb02b90ae2621ad5367) [@bep](https://github.com/bep) [#6503](https://github.com/gohugoio/hugo/issues/6503)[#6824](https://github.com/gohugoio/hugo/issues/6824)
* Fix base template handling with preceding comments [f45cb317](https://github.com/gohugoio/hugo/commit/f45cb3172862140883cfa08bd401c17e1ada5b39) [@bep](https://github.com/bep) [#6816](https://github.com/gohugoio/hugo/issues/6816)
* Fix "hugo mod get -u" with no arguments [49ef6472](https://github.com/gohugoio/hugo/commit/49ef6472039ede7d485242eba511207a8274495a) [@bep](https://github.com/bep) [#6826](https://github.com/gohugoio/hugo/issues/6826)[#6825](https://github.com/gohugoio/hugo/issues/6825)
* And now finally fix the 404 templates [74b6c4e5](https://github.com/gohugoio/hugo/commit/74b6c4e5ff5ee16f0e6b352a26c1e58b90a25dc6) [@bep](https://github.com/bep) [#6795](https://github.com/gohugoio/hugo/issues/6795)
* Fix 404 with base template regression [8df5d76e](https://github.com/gohugoio/hugo/commit/8df5d76e708238563185bac84809b34a4d395734) [@bep](https://github.com/bep) [#6795](https://github.com/gohugoio/hugo/issues/6795)
* Fix baseof with regular define regression [f441f675](https://github.com/gohugoio/hugo/commit/f441f675126ef1123d9f94429872dd683b40e011) [@bep](https://github.com/bep) [#6790](https://github.com/gohugoio/hugo/issues/6790)





