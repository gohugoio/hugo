
---
date: 2018-04-05
title: "Hugo 0.38.1: Some Live Reload Fixes"
description: "Hugo 0.38.1 fixes some live reload issues introduced in 0.38."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix that is mainly motivated by some issues with server live reloading introduced in Hugo 0.38.

* Fix livereload for the home page bundle [f87239e4](https://github.com/gohugoio/hugo/commit/f87239e4cab958bf59ecfb1beb8cac439441a553) [@bep](https://github.com/bep) [#4576](https://github.com/gohugoio/hugo/issues/4576)
* Fix empty BuildDate in "hugo version" [294c0f80](https://github.com/gohugoio/hugo/commit/294c0f8001fe598278c1eb8015deb6b98e8de686) [@anthonyfok](https://github.com/anthonyfok) 
* Fix some livereload content regressions [a4deaeff](https://github.com/gohugoio/hugo/commit/a4deaeff0cfd70abfbefa6d40c0b86839a216f6d) [@bep](https://github.com/bep) [#4566](https://github.com/gohugoio/hugo/issues/4566)
* Update github.com/bep/gitmap to fix snap build [4d115c56](https://github.com/gohugoio/hugo/commit/4d115c56fac9060230fbac6181a05f7cc6d10b42) [@anthonyfok](https://github.com/anthonyfok) [#4538](https://github.com/gohugoio/hugo/issues/4538)
* Fix two tests that are broken on Windows [26f34fd5](https://github.com/gohugoio/hugo/commit/26f34fd59da1ce1885d4f2909c5d9ef9c1726944) [@neurocline](https://github.com/neurocline) 


This release also contains some improvements:

* Add bash completion [874159b5](https://github.com/gohugoio/hugo/commit/874159b5436bc9080aec71a9c26d35f8f62c9fd0) [@anthonyfok](https://github.com/anthonyfok) 
* Handle mass content etc. edits in server mode [730b66b6](https://github.com/gohugoio/hugo/commit/730b66b6520f263af16f555d1d7be51205a8e51d) [@bep](https://github.com/bep) [#4563](https://github.com/gohugoio/hugo/issues/4563)






