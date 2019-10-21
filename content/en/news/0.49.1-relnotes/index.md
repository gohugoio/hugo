
---
date: 2018-10-10
title: "Hugo 0.49.1: Bug Fix"
description: "This release fixes an issue where resources.Concat would sometimes fail."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

This is a bug-fix release with 2 related fixes. This was introduced in Hugo 0.49. The most notable error situation was that `resources.Concat` could fail in some situations.


* Fix handling of different interface types in Slice [e2201ef1](https://github.com/gohugoio/hugo/commit/e2201ef15fdefe257ad284b2df4ccc8f8c38fac2) [@bep](https://github.com/bep) [#5269](https://github.com/gohugoio/hugo/issues/5269)

* Improve append in Scratch [23f48c30](https://github.com/gohugoio/hugo/commit/23f48c300cb5ffe0fe43c88464f38c68831a17ad) [@bep](https://github.com/bep) [#5275](https://github.com/gohugoio/hugo/issues/5275)





