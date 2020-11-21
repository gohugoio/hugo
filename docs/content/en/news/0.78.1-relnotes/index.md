
---
date: 2020-11-05
title: "Hugo 0.78.1: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.78.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

The main fix in this release is that of dependency resolution for package.json/node_modules in theme components. See [the documentation](https://gohugo.io/hugo-pipes/js/#include-dependencies-in-packagejson--node_modules) for more information.

* Disable NPM test on Travis on Windows [3437174c](https://github.com/gohugoio/hugo/commit/3437174c3a7b96925b82b351ac87530b4fa796a5) [@bep](https://github.com/bep) 
* travis: Install nodejs on Windows [f66302ca](https://github.com/gohugoio/hugo/commit/f66302ca0579171ffd1730eb8f33dd05af3d9a00) [@bep](https://github.com/bep) 
* js: Remove external source map option [944150ba](https://github.com/gohugoio/hugo/commit/944150bafbbb5c3e807ba3688174e70764dbdc64) [@bep](https://github.com/bep) [#7932](https://github.com/gohugoio/hugo/issues/7932)
* js: Misc fixes [bf2837a3](https://github.com/gohugoio/hugo/commit/bf2837a314eaf70135791984a423b0b09f58741d) [@bep](https://github.com/bep) [#7924](https://github.com/gohugoio/hugo/issues/7924)[#7923](https://github.com/gohugoio/hugo/issues/7923)



