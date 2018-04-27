
---
date: 2018-04-09
title: "Hugo 0.38.2: Two Bugfixes"
description: "0.38.2 fixes `--contentDir` flag handling and \".\" in content filenames."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix release with a couple of important fixes:


* Fix handling of the `--contentDir` and possibly other related flags [080302eb](https://github.com/gohugoio/hugo/commit/080302eb8757fd94ccbd6bf99103432cd98e716c) [@bep](https://github.com/bep) [#4589](https://github.com/gohugoio/hugo/issues/4589)
* Fix handling of content files with "." in them [2817e842](https://github.com/gohugoio/hugo/commit/2817e842407c8dcbfc738297ab634392fcb41ce1) [@bep](https://github.com/bep) [#4559](https://github.com/gohugoio/hugo/issues/4559)


Also in this release:

* Set .Parent in bundled pages to its owner [6792d86a](https://github.com/gohugoio/hugo/commit/6792d86ad028571c684a776c5f00e0107838c955) [@bep](https://github.com/bep) [#4582](https://github.com/gohugoio/hugo/issues/4582)



