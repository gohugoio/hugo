
---
date: 2019-04-20
title: "Hugo 0.55.3: A Few More Bug Fixes!"
description: "To wrap up this Easter, here is one more patch release with some important fixes."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

This is a bug-fix release with a couple of important fixes.


Hugo now has:

* 34468+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 308+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Return error on invalid input in in [7fbfedf0](https://github.com/gohugoio/hugo/commit/7fbfedf01367ff076c3c875b183789b769b99241) [@bep](https://github.com/bep) [#5875](https://github.com/gohugoio/hugo/issues/5875)
* Make Pages etc. work with the in func [06f56fc9](https://github.com/gohugoio/hugo/commit/06f56fc983d460506d39b3a6f638b1632af07073) [@bep](https://github.com/bep) [#5875](https://github.com/gohugoio/hugo/issues/5875)
* Make Pages etc. work in uniq [d7a67dcb](https://github.com/gohugoio/hugo/commit/d7a67dcb51829b12d492d3f2ee4f6e2a3834da63) [@bep](https://github.com/bep) [#5852](https://github.com/gohugoio/hugo/issues/5852)

### Core

* Add some OutputFormats.Get tests [7aeeb60d](https://github.com/gohugoio/hugo/commit/7aeeb60d7ee71690461df92ff41cb8b2f7f5aa61) [@bep](https://github.com/bep) [#5877](https://github.com/gohugoio/hugo/issues/5877)
* Add some integration tests for in/uniq using Pages [6c80acbd](https://github.com/gohugoio/hugo/commit/6c80acbd5e314dd92fc075551ffabafaae01dca7) [@bep](https://github.com/bep) [#5875](https://github.com/gohugoio/hugo/issues/5875)[#5852](https://github.com/gohugoio/hugo/issues/5852)

### Other

* Regenerate docs helper [75b16e30](https://github.com/gohugoio/hugo/commit/75b16e30ec55e82a8024cc4d27880d9b79e0fa41) [@bep](https://github.com/bep) 
* Replace IsDraft with Draft in list command [3e421bd4](https://github.com/gohugoio/hugo/commit/3e421bd47cd35061df89c1c127ec8fa4ae368449) [@bep](https://github.com/bep) [#5873](https://github.com/gohugoio/hugo/issues/5873)

## Fixes

### Output

* Fix links for non-HTML output formats [c7dd66bf](https://github.com/gohugoio/hugo/commit/c7dd66bfe2e32430f9b1a3126c67014e40d8405e) [@bep](https://github.com/bep) [#5877](https://github.com/gohugoio/hugo/issues/5877)
* Fix menu URL when multiple permalinkable output formats [ea529c84](https://github.com/gohugoio/hugo/commit/ea529c847ebc0267c6d0426cc8f77d5c76c73fe4) [@bep](https://github.com/bep) [#5849](https://github.com/gohugoio/hugo/issues/5849)





