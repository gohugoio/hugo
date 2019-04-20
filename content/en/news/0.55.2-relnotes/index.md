
---
date: 2019-04-17
title: "Hugo 0.55.2: Some Important Bug Fixes"
description: "Fixes some more issues introduced in Hugo 0.55.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix release with a couple of important fixes.


Hugo now has:

* 34386+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 307+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Handle late transformation of templates [2957795f](https://github.com/gohugoio/hugo/commit/2957795f5276cc9bc8d438da2d7d9b61defea225) [@bep](https://github.com/bep) [#5865](https://github.com/gohugoio/hugo/issues/5865)

### Core

* Add more tests for Permalinkable [35f41834](https://github.com/gohugoio/hugo/commit/35f41834ea3a8799b9b7eda360cf8d30b1b727ba) [@bep](https://github.com/bep) [#5849](https://github.com/gohugoio/hugo/issues/5849)

## Fixes

### Core

* Fix Pages reinitialization on rebuilds [9b17cbb6](https://github.com/gohugoio/hugo/commit/9b17cbb62a056ea7e26b1146cbf3ba42f5acf805) [@bep](https://github.com/bep) [#5833](https://github.com/gohugoio/hugo/issues/5833)
* Fix shortcode namespace issue [56550d1e](https://github.com/gohugoio/hugo/commit/56550d1e449f45ebee398ac8a9e3b9818b3ee60e) [@bep](https://github.com/bep) [#5863](https://github.com/gohugoio/hugo/issues/5863)
* Fix false WARNINGs in lang prefix check [7881b096](https://github.com/gohugoio/hugo/commit/7881b0965f8b83d03379e9ed102cd0c3bce297e2) [@bep](https://github.com/bep) [#5860](https://github.com/gohugoio/hugo/issues/5860)
* Fix bundle resource publishing when multiple output formats [49d0a826](https://github.com/gohugoio/hugo/commit/49d0a82641581aa7dd66b9d5e8c7d75e23260083) [@bep](https://github.com/bep) [#5858](https://github.com/gohugoio/hugo/issues/5858)
* Fix panic for unused taxonomy content files [b799b12f](https://github.com/gohugoio/hugo/commit/b799b12f4a693dfeae8a5a362f131081a727bb8f) [@bep](https://github.com/bep) [#5847](https://github.com/gohugoio/hugo/issues/5847)
* Fix dates for sections with dates in front matter [70148672](https://github.com/gohugoio/hugo/commit/701486728e21bc0c6c78c2a8edb988abdf6116c7) [@bep](https://github.com/bep) [#5854](https://github.com/gohugoio/hugo/issues/5854)

### Other

* Fix WeightedPages in union etc. [f2795d4d](https://github.com/gohugoio/hugo/commit/f2795d4d2cef30170af43327f3ff7114923833b1) [@bep](https://github.com/bep) [#5850](https://github.com/gohugoio/hugo/issues/5850)





