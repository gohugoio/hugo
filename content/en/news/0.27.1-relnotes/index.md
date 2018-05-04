
---
date: 2017-09-13
title: "Hugo 0.27.1: One bugfix"
description: "Hugo 0.27.1 fixes an issue introduced in Go 1.9 with HTML escaping of shortcodes in multi output sites."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png
---

	

This fixes a regression introduced in [Go 1.9](https://github.com/golang/go/issues/21844) which lead to HTML in shortcodes in multi output Hugo sites being wrongly escaped in some cases.

* Fix escaped HTML Go 1.9 multioutput issue (#3880) [2d613dd9](https://github.com/gohugoio/hugo/commit/2d613dd905bb8eeb8af57e30ddd749a0f04fbd3c) [@bep](https://github.com/bep) [#3876](https://github.com/gohugoio/hugo/issues/3876)

* Bump to Go 1.9 in the Snap build [642ba6ca](https://github.com/gohugoio/hugo/commit/642ba6cab24c558b16378178fe829cbc45845424) [@bep](https://github.com/bep) 




