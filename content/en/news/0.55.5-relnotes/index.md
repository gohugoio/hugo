
---
date: 2019-05-02
title: "Hugo 0.55.5: Take Five!"
description: "We round up this 0.55 release with a final batch of bug fixes!"
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix release with a couple of important fixes.


Hugo now has:

* 34743+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 314+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Provide more detailed errors in Where [f76e5011](https://github.com/gohugoio/hugo/commit/f76e50118b8b5dd8989d068db35222bfa0a242d8) [@moorereason](https://github.com/moorereason) 

### Other

* Regenerate docs helper [cee181c3](https://github.com/gohugoio/hugo/commit/cee181c3a67fe04b8e0c9f2807c5aa2871df474e) [@bep](https://github.com/bep) 
* Update blackfriday to v1.5.2 [1cbb501b](https://github.com/gohugoio/hugo/commit/1cbb501be8aa83b08865fbb6ad5aee254946712f) [@dbirks](https://github.com/dbirks) 

## Fixes

### Templates

* Fix nil compare in eq/ne for interface values [66b143a0](https://github.com/gohugoio/hugo/commit/66b143a01d1c192619839b732ce188923ab15d60) [@bep](https://github.com/bep) [#5905](https://github.com/gohugoio/hugo/issues/5905)
* Fix hugo package name and add godocs [4f93f8c6](https://github.com/gohugoio/hugo/commit/4f93f8c670b26258dc7e3a613c38dbc86d8eda76) [@moorereason](https://github.com/moorereason) 

### Output

* Fix permalink in sitemap etc. when multiple permalinkable output formats [6b76841b](https://github.com/gohugoio/hugo/commit/6b76841b052b97625b8995f326d758b89f5c2349) [@bep](https://github.com/bep) [#5910](https://github.com/gohugoio/hugo/issues/5910)

### Core

* Fix PrevInSection/NextInSection for nested sections [bcbed4eb](https://github.com/gohugoio/hugo/commit/bcbed4ebdaf55b67abc521d69bba456c041a7e7d) [@bep](https://github.com/bep) [#5883](https://github.com/gohugoio/hugo/issues/5883)

### Other

* Fix concurrent initialization order [009076e5](https://github.com/gohugoio/hugo/commit/009076e5ee88fc46c95a9afd34f82f9386aa282a) [@bep](https://github.com/bep) [#5901](https://github.com/gohugoio/hugo/issues/5901)





