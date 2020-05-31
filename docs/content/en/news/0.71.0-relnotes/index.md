
---
date: 2020-05-18
title: "Markdown Render Hooks for Headings"
description: "Render hooks for headings, update to Go 1.14.3, several bug fixes etc."
categories: ["Releases"]
---

Hugo 0.71 brings [Markdown render hooks for headings](https://gohugo.io/getting-started/configuration-markup#markdown-render-hooks), a set of bug fixes and more.

This release represents **12 contributions by 7 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@apexskier](https://github.com/apexskier), and [@johnweldon](https://github.com/johnweldon) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **8 contributions by 7 contributors**. A special thanks to [@bep](https://github.com/bep), [@mikeee](https://github.com/mikeee), [@h-enk](https://github.com/h-enk), and [@tjamet](https://github.com/tjamet) for their work on the documentation site.


Hugo now has:

* 44043+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 437+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 322+ [themes](http://themes.gohugo.io/)

## Enhancements

### Other

* Use WARN log level also for the early initialization [518d1496](https://github.com/gohugoio/hugo/commit/518d149646c13fb49c296a63e61a048f5e672179) [@bep](https://github.com/bep) [#7285](https://github.com/gohugoio/hugo/issues/7285)
* Update to Go 1.14.3 and Go 1.13.11 [3cc41523](https://github.com/gohugoio/hugo/commit/3cc41523bef802d1942f3d31018547a18cc55923) [@bep](https://github.com/bep) 
* Improve error message when no Babel installed [2fd0a5a6](https://github.com/gohugoio/hugo/commit/2fd0a5a6781456e88745b370d12aaf5351a020ff) [@bep](https://github.com/bep) 
* Add test for headings render hook [6e051c05](https://github.com/gohugoio/hugo/commit/6e051c053e2b5b8419f357ed8acd177440266d07) [@apexskier](https://github.com/apexskier) 
* Add render template hooks for headings [423b8f2f](https://github.com/gohugoio/hugo/commit/423b8f2fb834139cf31514b14b1c1bf28e43b384) [@elihunter173](https://github.com/elihunter173) [#6713](https://github.com/gohugoio/hugo/issues/6713)
* Add math.Pow [99193449](https://github.com/gohugoio/hugo/commit/991934497e88dcd4134a369a213bb5072c51c139) [@jmooring](https://github.com/jmooring) [#7266](https://github.com/gohugoio/hugo/issues/7266)
* Do not suppress .well-known/ directory [558c0930](https://github.com/gohugoio/hugo/commit/558c09305e2be16953238c6c0e828f62b950e4f5) [@johnweldon](https://github.com/johnweldon) [#6691](https://github.com/gohugoio/hugo/issues/6691)
* Quote "@babel/cli" to solve build error [b69a3614](https://github.com/gohugoio/hugo/commit/b69a36140f42ec99ffa2d1e029b8b86ecf8ff929) [@anthonyfok](https://github.com/anthonyfok) 
* Remove custom x-nodejs plugin [a0103864](https://github.com/gohugoio/hugo/commit/a0103864ab76c6a1462a6dee538801740acf4858) [@anthonyfok](https://github.com/anthonyfok) 
* Use .Lastmod for og:updated_time [6205d56b](https://github.com/gohugoio/hugo/commit/6205d56b85fea31e008cd0fef26805bab8084786) [@dtip](https://github.com/dtip) 

## Fixes

### Other

* Fix Babel on Windows [723ec555](https://github.com/gohugoio/hugo/commit/723ec555e75fbfa94d90d3ecbcd5775d6c7800e1) [@bep](https://github.com/bep) [#7251](https://github.com/gohugoio/hugo/issues/7251)
* Upgrade chroma to 0.7.3 to fix invalid css [b342e8fb](https://github.com/gohugoio/hugo/commit/b342e8fbdb23157f3979af91cb5d8d3438003707) [@apexskier](https://github.com/apexskier) [#7207](https://github.com/gohugoio/hugo/issues/7207)





