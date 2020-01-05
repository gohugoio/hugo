---
date: 2020-01-01
title: "Hugo 0.62.1: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.62.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png
---

This release is mainly motivated by getting [this demo site](https://github.com/bep/portable-hugo-links) up and running. It demonstrates truly portable Markdown links and images, whether browsed on GitHub or deployed as a Hugo site.

* Support files in content mounts [ff6253bc](https://github.com/gohugoio/hugo/commit/ff6253bc7cf745e9c0127ddc9006da3c2c00c738) [@bep](https://github.com/bep) [#6684](https://github.com/gohugoio/hugo/issues/6684)[#6696](https://github.com/gohugoio/hugo/issues/6696)
* Update alpine base image in Dockerfile to 3.11 [aa4ccb8a](https://github.com/gohugoio/hugo/commit/aa4ccb8a1e9b8aa17397acf34049a2aa16b0b6cb) [@RemcodM](https://github.com/RemcodM) 
* hugolib: Fix inline shortcode regression [5509954c](https://github.com/gohugoio/hugo/commit/5509954c7e8b0ce8d5ea903b0ab639ea14b69acb) [@bep](https://github.com/bep) [#6677](https://github.com/gohugoio/hugo/issues/6677)



