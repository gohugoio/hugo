
---
date: 2018-01-02
title: "Hugo 0.32.1: Two bugfixes"
description: "Fixes image processing in shortcodes."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png
---

This fixes 2 bugs from the Hugo 0.32 release.

* Fix image processing from shortcodes in non-server mode. [@bep](https://github.com/bep) [#4202](https://github.com/gohugoio/hugo/issues/4202)
* Fix broken `hugo --renderToMemory`.  Note that this is only useful for benchmark testing, as there is no easy way to actually view the result. [d36d71ed](https://github.com/gohugoio/hugo/commit/d36d71edd3b04df3b34edf4d108e3995a244c4f0) [@bep](https://github.com/bep) [#4212](https://github.com/gohugoio/hugo/issues/4212)




