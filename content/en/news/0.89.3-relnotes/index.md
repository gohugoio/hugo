
---
date: 2021-11-15
title: "Hugo 0.89.3: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.89.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix release with a couple of important fixes.

* Improve error when we cannot determine content directory in "hugo new" [b8155452](https://github.com/gohugoio/hugo/commit/b8155452ac699473b6b2a34f53988dee01b4da34) [@bep](https://github.com/bep) [#9166](https://github.com/gohugoio/hugo/issues/9166)
* deps: Upgrade github.com/yuin/goldmark v1.4.3 => v1.4.4 [08552a7a](https://github.com/gohugoio/hugo/commit/08552a7a4cd1fe64efdd2f1b95142fa4295cb298) [@jmooring](https://github.com/jmooring) [#9159](https://github.com/gohugoio/hugo/issues/9159)
* commands: Make sure pollInterval is always set [fdad91fd](https://github.com/gohugoio/hugo/commit/fdad91fd96bc4636bf3a957cdddce18b66473124) [@bep](https://github.com/bep) [#9165](https://github.com/gohugoio/hugo/issues/9165)
* create: Improve archetype directory discovery and tests [5f3f6089](https://github.com/gohugoio/hugo/commit/5f3f60898cfe1c087841ec1fbd5ddc2916d0a2c6) [@bep](https://github.com/bep) [#9146](https://github.com/gohugoio/hugo/issues/9146)
* create: Add a log statement when archetype is a directory [057d02de](https://github.com/gohugoio/hugo/commit/057d02de256a3866b7044abaa4d03c69d9fedef0) [@bep](https://github.com/bep) [#9157](https://github.com/gohugoio/hugo/issues/9157)
* create: Always print "Content ... created" [43ac59da](https://github.com/gohugoio/hugo/commit/43ac59da850901cc848b35129ca7223f9f9a9b19) [@bep](https://github.com/bep) [#9157](https://github.com/gohugoio/hugo/issues/9157)
* commands: Fix missing file locking in server partial render [ab5c6990](https://github.com/gohugoio/hugo/commit/ab5c6990a55cbb11d97f857b4619b83fddda3d18) [@bep](https://github.com/bep) [#9162](https://github.com/gohugoio/hugo/issues/9162)
* modules: Improve error message [9369d13e](https://github.com/gohugoio/hugo/commit/9369d13e59ffac262944477fad3dcd2742d66288) [@davidsneighbour](https://github.com/davidsneighbour) 



