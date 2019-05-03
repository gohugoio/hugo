
---
date: 2019-04-25
title: "Hugo 0.55.4: Some Bug Fixes"
description: "A couple of more bug fixes."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix release with a couple of important fixes.


Hugo now has:

* 34558+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 310+ [themes](http://themes.gohugo.io/)

## Enhancements

### Core

* Avoid recloning of shortcode templates [69a56420](https://github.com/gohugoio/hugo/commit/69a56420aec5bf5abb846701d4a5ec67fe060d96) [@bep](https://github.com/bep) [#5890](https://github.com/gohugoio/hugo/issues/5890)
* No links for bundled pages [0775c98e](https://github.com/gohugoio/hugo/commit/0775c98e6c5b700e46adaaf190fc3f693a6ab002) [@bep](https://github.com/bep) [#5882](https://github.com/gohugoio/hugo/issues/5882)

### Other

* Avoid rebuilding the Translations map for every lookup [4756ec3c](https://github.com/gohugoio/hugo/commit/4756ec3cd8ef998f889619fe11be70cc900e2b75) [@bep](https://github.com/bep) [#5892](https://github.com/gohugoio/hugo/issues/5892)
* Init mem profile at the end [4c3c5120](https://github.com/gohugoio/hugo/commit/4c3c5120389cc95edc63b8f18a0eee786aa0c5e2) [@bep](https://github.com/bep) 

## Fixes

### Core

* Fix shortcode version=1 logic [33c73811](https://github.com/gohugoio/hugo/commit/33c738116c26e2ac37f4bd48159e8e3197fd7b39) [@bep](https://github.com/bep) [#5831](https://github.com/gohugoio/hugo/issues/5831)





