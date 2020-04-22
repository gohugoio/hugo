
---
date: 2020-04-22
title: "Hugo 0.69.1: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.69.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix release with a couple of important fixes.

*  hugolib/filesystems: Fix typo in test suite [49e6c8cb](https://github.com/gohugoio/hugo/commit/49e6c8cb4ed83e20f1e0ac164e91c38854177b99) [@panakour](https://github.com/panakour) 
* Fix class collector when running with --minify [f37e77f2](https://github.com/gohugoio/hugo/commit/f37e77f2d338cf876cfa637a662acd76f0f2009b) [@bep](https://github.com/bep) [#7161](https://github.com/gohugoio/hugo/issues/7161)
* related: Fix toLower [27af5a33](https://github.com/gohugoio/hugo/commit/27af5a339a4d3c5712b5ed946a636a8c21916039) [@bep](https://github.com/bep) [#7198](https://github.com/gohugoio/hugo/issues/7198)
* Fix broken test [b3c82575](https://github.com/gohugoio/hugo/commit/b3c825756f3251f8b26e53262f9d6f484aecf750) [@bep](https://github.com/bep) 
* tpl/tmplimpl/template: Change defer RLock to RUnlock [5146dc61](https://github.com/gohugoio/hugo/commit/5146dc614fc45df698ebf890af06421dea988c96) [@BurtonQin](https://github.com/BurtonQin) 
* hugolib: Add Unlock before panic [736f84b2](https://github.com/gohugoio/hugo/commit/736f84b2d539857f7fdd0e42353af80b4dccfe8d) [@BurtonQin](https://github.com/BurtonQin) 
* docs: Fix typo in Hugo's Security Model [cd4d8202](https://github.com/gohugoio/hugo/commit/cd4d8202016bd3eb5ed9144c8945edaba73c8cf4) [@sensimevanidus](https://github.com/sensimevanidus) 
* deps: Update go-org to v1.1.0 [2b28e5a9](https://github.com/gohugoio/hugo/commit/2b28e5a9cb79af2a8d70c80036f52bcf5399b9df) [@niklasfasching](https://github.com/niklasfasching) 
* commands: Modify gen chromastyles to output all CSS classes [102ec2da](https://github.com/gohugoio/hugo/commit/102ec2da7adcc4afb7050b17989f0486f8379679) [@acahir](https://github.com/acahir) [#7167](https://github.com/gohugoio/hugo/issues/7167)
* deps: Update to goldmark v1.1.28 [feaa582c](https://github.com/gohugoio/hugo/commit/feaa582cbe950e82969da5e99e3fb9a3947025df) [@bep](https://github.com/bep) [#7113](https://github.com/gohugoio/hugo/issues/7113)
* Fix query parameter handling in server fast render mode [ee67dbef](https://github.com/gohugoio/hugo/commit/ee67dbeff5bae6941facaaa39cb995a1ee6def03) [@bep](https://github.com/bep) [#7163](https://github.com/gohugoio/hugo/issues/7163)



