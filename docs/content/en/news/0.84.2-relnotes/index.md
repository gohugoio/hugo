
---
date: 2021-06-28
title: "Hugo 0.84.2: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.84.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is mostly a bug fix release, but it also contains some minor modules related improvements. Most notable you now get some more information in ` hugo config mounts`, and even more so when typing ` hugo config mounts -v`.

* modules: Add module.import.noMounts config [40dfdd09](https://github.com/gohugoio/hugo/commit/40dfdd09521bcb8f56150e6791d60445198f27ab) [@bep](https://github.com/bep) [#8708](https://github.com/gohugoio/hugo/issues/8708)
* modules: Use value type for module.Time [3a6dc6d3](https://github.com/gohugoio/hugo/commit/3a6dc6d3f423c4acb79ef21b5a76e616fa2c9477) [@bep](https://github.com/bep) 
* commands: Add version time to "hugo config mounts" [6cd2110a](https://github.com/gohugoio/hugo/commit/6cd2110ab295f598907a18da91e34d31407c1d9d) [@bep](https://github.com/bep) 
* commands: Add some more info to "hugo config mounts" [6a365c27](https://github.com/gohugoio/hugo/commit/6a365c2712c7607e067e192d213b266f0c88d0f3) [@bep](https://github.com/bep) 
* Fix config handling with empty config entries after merge [19aa95fc](https://github.com/gohugoio/hugo/commit/19aa95fc7f4cd58dcc8a8ff075762cfc86d41dc3) [@bep](https://github.com/bep) [#8701](https://github.com/gohugoio/hugo/issues/8701)
* Fix config loading for "hugo mod init" [923dd9d1](https://github.com/gohugoio/hugo/commit/923dd9d1c1f649142f3f377109318b07e0f44d5d) [@bep](https://github.com/bep) [#8697](https://github.com/gohugoio/hugo/issues/8697)
* deps: Update to Minify v2.9.18 [d9bdd37d](https://github.com/gohugoio/hugo/commit/d9bdd37d35ccd436b4dd470ef99efa372a6a086b) [@bep](https://github.com/bep) [#8693](https://github.com/gohugoio/hugo/issues/8693)
* Remove credit from release notes [b2eaf4c8](https://github.com/gohugoio/hugo/commit/b2eaf4c8c2e31aa1c1bc4a2c0061f661e01d2de1) [@digitalcraftsman](https://github.com/digitalcraftsman) 



