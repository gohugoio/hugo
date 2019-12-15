
---
date: 2019-12-11
title: "40K GitHub Stars Edition"
description: "40K stars on GitHub is a good enough reason to release a new version of Hugo!"
categories: ["Releases"]
---

This is the [40K GitHub Stars Edition](https://github.com/gohugoio/hugo/stargazers). It's mostly a bug fix release, and an important note is the deprecation of Amber and Ace as template engines. See [#6609](https://github.com/gohugoio/hugo/issues/6609) for more information.

This release represents **10 contributions by 3 contributors** to the main Hugo code base.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **5 contributions by 4 contributors**. A special thanks to [@YuriyOborozhnyi](https://github.com/YuriyOborozhnyi), [@bep](https://github.com/bep), [@Flogex](https://github.com/Flogex), and [@atishay](https://github.com/atishay) for their work on the documentation site.


Hugo now has:

* 40029+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 440+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 279+ [themes](http://themes.gohugo.io/)


## Notes

* Deprecate Ace and Amber [5f8c2818](https://github.com/gohugoio/hugo/commit/5f8c2818f120b881f58f4cec67aed876edb8bcdf) [@bep](https://github.com/bep) [#6609](https://github.com/gohugoio/hugo/issues/6609)

## Enhancements

### Templates

* Allow any key type in partialCached [0efb00c2](https://github.com/gohugoio/hugo/commit/0efb00c2a86ec3f52000a643f26f54bb2a9dfbd6) [@bep](https://github.com/bep) [#6572](https://github.com/gohugoio/hugo/issues/6572)

### Other

* Update Goldmark [3cc217a6](https://github.com/gohugoio/hugo/commit/3cc217a650546b8bc29deabb95e648aacef96fbf) [@bep](https://github.com/bep) 
* Add typographic chars from goldmark to toc [c5f2f583](https://github.com/gohugoio/hugo/commit/c5f2f5837fdf6a30c7b28e8368033623b74a30a0) [@tangiel](https://github.com/tangiel) [#6592](https://github.com/gohugoio/hugo/issues/6592)
* Reimplement pygmentsCodefencesGuessSyntax [40a092b0](https://github.com/gohugoio/hugo/commit/40a092b0687d44ecb53ef1fd53001a6299345780) [@bep](https://github.com/bep) [#6565](https://github.com/gohugoio/hugo/issues/6565)
* Update Goldmark [d534ce94](https://github.com/gohugoio/hugo/commit/d534ce9424c952800dfb26c2faff2d47e9597cad) [@bep](https://github.com/bep) [#6557](https://github.com/gohugoio/hugo/issues/6557)
* Update minify [86a5b59f](https://github.com/gohugoio/hugo/commit/86a5b59f64dd6c4d338a9e091e98cd0ad6d4824f) [@MeiK2333](https://github.com/MeiK2333) [#6475](https://github.com/gohugoio/hugo/issues/6475)
* Update Goldmark [347cfb0c](https://github.com/gohugoio/hugo/commit/347cfb0c17b08626250180e8a84b53fc4800473f) [@bep](https://github.com/bep) [#6549](https://github.com/gohugoio/hugo/issues/6549)[#6551](https://github.com/gohugoio/hugo/issues/6551)

## Fixes

### Core

* Fix timeout number parsing for YAML/JSON config [b60ae35b](https://github.com/gohugoio/hugo/commit/b60ae35b97c4f44b9b09fcf06c863c695bc3c73a) [@bep](https://github.com/bep) [#6555](https://github.com/gohugoio/hugo/issues/6555)

### Other

* Fix headless regression [bb80fff6](https://github.com/gohugoio/hugo/commit/bb80fff69ad3f2ddff23819bf6eb6f4b8512dc2a) [@bep](https://github.com/bep) [#6552](https://github.com/gohugoio/hugo/issues/6552)





