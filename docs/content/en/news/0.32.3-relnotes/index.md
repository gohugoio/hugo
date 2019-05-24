
---
date: 2018-01-08
title: "Hugo 0.32.3: Some important bug fixes"
description: "Fixes multilingual resource (images etc.) handling etc."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

Hugo `0.32` was a big and [really cool](https://gohugo.io/news/0.32-relnotes/) release, and the [Hugo Forum](https://discourse.gohugo.io/) has been filled with questions from people wanting to upgrade their Hugo sites to be able to use the new image processing feature etc.

And with that we have discovered some issues, which this release should fix, mostly releated to multilingual sites:

* Fix multihost detection for sites without language definition [8969331f](https://github.com/gohugoio/hugo/commit/8969331f5be352939883074034adac6b7086ddc8) [@bep](https://github.com/bep) [#4221](https://github.com/gohugoio/hugo/issues/4221)
* Fix hugo benchmark --renderToMemory [059e8458](https://github.com/gohugoio/hugo/commit/059e8458d690dbb9fcd3ebd58cfc61b062d3138e) [@bep](https://github.com/bep) [#4218](https://github.com/gohugoio/hugo/issues/4218)
* Fix URLs for bundle resources in multihost mode [ab82a27d](https://github.com/gohugoio/hugo/commit/ab82a27d055c3aa177821d81a45a5c6e972aa29e) [@bep](https://github.com/bep) [#4217](https://github.com/gohugoio/hugo/issues/4217)
* Fix sub-folder baseURL handling for Page resources [f25d8a9e](https://github.com/gohugoio/hugo/commit/f25d8a9e17fb65fa41dafdcbf0358853d68eaf45) [@bep](https://github.com/bep) [#4228](https://github.com/gohugoio/hugo/issues/4228)
* Avoid processing and storing same image for each language [4b04db0f](https://github.com/gohugoio/hugo/commit/4b04db0f0855a1f54895d6c93c52dcea4b1ce3ca) [@bep](https://github.com/bep) [#4231](https://github.com/gohugoio/hugo/issues/4231)
* Resources.ByType should return Resources [97c1866e](https://github.com/gohugoio/hugo/commit/97c1866e322284dec46db6f3d235807507f5b69f) [@bep](https://github.com/bep) [#4234](https://github.com/gohugoio/hugo/issues/4234)
* Report build time on config.toml change [6feb1387](https://github.com/gohugoio/hugo/commit/6feb138785eeb9e813428d0df30010d9b5fb1059) [@bep](https://github.com/bep) [#4232](https://github.com/gohugoio/hugo/issues/4232)[#4224](https://github.com/gohugoio/hugo/issues/4224)
* Fix handling of mixed-case taxonomy folders with content file [2d3189b2](https://github.com/gohugoio/hugo/commit/2d3189b22760e0a8995dae082a6bc5480f770bfe) [@bep](https://github.com/bep) [#4238](https://github.com/gohugoio/hugo/issues/4238)





