
---
date: 2020-12-31
title: "Hugo 0.80: Last Release of 2020!"
description: "This release brings Dart Sass support, a new image overlay function, and more."
categories: ["Releases"]
---

The last Hugo release of the year brings a new [images.Overlay](https://gohugo.io/functions/images/#overlay) filter to overlay an image on top of another, e.g. for watermarking, and [Dart Sass](https://gohugo.io/hugo-pipes/scss-sass/#options) support.

This release represents **29 contributions by 12 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), and [@davidsneighbour](https://github.com/davidsneighbour) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **22 contributions by 6 contributors**. A special thanks to [@bep](https://github.com/bep), [@D4D3VD4V3](https://github.com/D4D3VD4V3), [@chrischute](https://github.com/chrischute), and [@azenk](https://github.com/azenk) for their work on the documentation site.


Hugo now has:

* 49096+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 436+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 369+ [themes](http://themes.gohugo.io/)

## Notes

* Resource.ResourceType now always returns MIME's main type [81975f84](https://github.com/gohugoio/hugo/commit/81975f847dc19c21c2321207645807771db97fab) [@bep](https://github.com/bep) [#8052](https://github.com/gohugoio/hugo/issues/8052)

## Enhancements

### Templates

* Regenerate templates [a2d146ec](https://github.com/gohugoio/hugo/commit/a2d146ec32a26ccca9ffa68d3c840ec5b08cca96) [@bep](https://github.com/bep) 
* tpl/internal/go_templates: Revert formatting [718e09ed](https://github.com/gohugoio/hugo/commit/718e09ed4bc538f4fccc4337f99e9eb86aea31f3) [@bep](https://github.com/bep) 
* Add title parameter to YouTube shortcode [4fc918e0](https://github.com/gohugoio/hugo/commit/4fc918e02cfc7f260d6312248ff9d33e95b27943) [@azenk](https://github.com/azenk) 

### Output

* Add missing OutputStyle option [428b0b32](https://github.com/gohugoio/hugo/commit/428b0b32947ec16f8585b8c33548d72fd4fb025d) [@bep](https://github.com/bep) 

### Other

* Allow Dart Sass transformations to be cached on disk [ffbf5e45](https://github.com/gohugoio/hugo/commit/ffbf5e45fa0617a37950b34deab63736b1c6b1d3) [@bep](https://github.com/bep) 
* Dart Sass only supports `expanded` and `compressed` [48994ea7](https://github.com/gohugoio/hugo/commit/48994ea766f08332f57c0f8e74843b6c8617c3d1) [@bep](https://github.com/bep) 
* Update emoji import paths and version [1f7e9f73](https://github.com/gohugoio/hugo/commit/1f7e9f733397b891cefc725ffc94ba901e70425a) [@moorereason](https://github.com/moorereason) 
* Add Dart Sass support [cea15740](https://github.com/gohugoio/hugo/commit/cea157402365f34a69882110a4208999728007a6) [@bep](https://github.com/bep) [#7380](https://github.com/gohugoio/hugo/issues/7380)[#8102](https://github.com/gohugoio/hugo/issues/8102)
* GroupByParamDate now supports datetimes [f9f77978](https://github.com/gohugoio/hugo/commit/f9f779786edcefc4449a14cfc04dd93379f71373) [@zerok](https://github.com/zerok) 
* Skip para test when not on CI [a9718f44](https://github.com/gohugoio/hugo/commit/a9718f44cd6c938448fc697f0ec720ebed7d863a) [@bep](https://github.com/bep) [#6963](https://github.com/gohugoio/hugo/issues/6963)
* Update SECURITY.md [f802bb23](https://github.com/gohugoio/hugo/commit/f802bb236a60dcc6c64d53edac634891272e0c07) [@bep](https://github.com/bep) 
* Improve LookPath [10ae7c32](https://github.com/gohugoio/hugo/commit/10ae7c3210cd1add14d3750aa9512a87df0e1146) [@bep](https://github.com/bep) 
* create a SECURITY.md [ae2d1bd5](https://github.com/gohugoio/hugo/commit/ae2d1bd52df0099190ef9195666d0788708b0385) [@davidsneighbour](https://github.com/davidsneighbour) [#8074](https://github.com/gohugoio/hugo/issues/8074)
* Show more detail on failed time test [8103188b](https://github.com/gohugoio/hugo/commit/8103188b9b9e8eeb3bcb53c8b64e2b83397e82ae) [@moorereason](https://github.com/moorereason) [#6963](https://github.com/gohugoio/hugo/issues/6963)
* Add images.Overlay filter [3ba147e7](https://github.com/gohugoio/hugo/commit/3ba147e702a5ae0af6e8b3b0296d256c3246a546) [@bep](https://github.com/bep) [#8057](https://github.com/gohugoio/hugo/issues/8057)[#4595](https://github.com/gohugoio/hugo/issues/4595)[#6731](https://github.com/gohugoio/hugo/issues/6731)
* Bump github.com/spf13/cobra from 0.15.0 to 0.20.0 [c84ad8db](https://github.com/gohugoio/hugo/commit/c84ad8db821c10225c0e603c6ec920c67b6ce36f) [@anthonyfok](https://github.com/anthonyfok) 
* configure proper link to discourse.gohugo.io (#8020) [4e0acb89](https://github.com/gohugoio/hugo/commit/4e0acb89b793d8895dc53eb8887be27430c3ab31) [@davidsneighbour](https://github.com/davidsneighbour) 
* Format code with gofumpt [d90e37e0](https://github.com/gohugoio/hugo/commit/d90e37e0c6e812f9913bf256c9c81aa05b7a08aa) [@bep](https://github.com/bep) 
* bump github.com/evanw/esbuild from 0.8.15 to 0.8.17 [32471b57](https://github.com/gohugoio/hugo/commit/32471b57bde51c55a15dbf1db75d6e5f7232c347) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Use --baseURL path for live-reload URL [0ad378b0](https://github.com/gohugoio/hugo/commit/0ad378b09cea90a2a70d7ff06af668abe22475a1) [@sth](https://github.com/sth) [#6595](https://github.com/gohugoio/hugo/issues/6595)
* bump github.com/getkin/kin-openapi from 0.31.0 to 0.32.0 [907d9e92](https://github.com/gohugoio/hugo/commit/907d9e92682ed56a57a2206ae9bd9a985b3e1870) [@dependabot[bot]](https://github.com/apps/dependabot) 

## Fixes

### Templates

* Fix series detection in opengraph [d2d493ab](https://github.com/gohugoio/hugo/commit/d2d493ab5d6a054001a8448ea0de2949dac4b30e) [@Humberd](https://github.com/Humberd) 
* Fix substr when length parameter is zero [5862fd2a](https://github.com/gohugoio/hugo/commit/5862fd2a60b5d16f2437bd8c8b7bac700de5f047) [@moorereason](https://github.com/moorereason) [#7993](https://github.com/gohugoio/hugo/issues/7993)
* Refactor and fix substr logic [64789fb5](https://github.com/gohugoio/hugo/commit/64789fb5dcf8326f14f13d69a2576ae3aa2bbbaa) [@moorereason](https://github.com/moorereason) [#7993](https://github.com/gohugoio/hugo/issues/7993)

### Other

* Fix Resource.ResourceType so it always returns MIME's main type [81975f84](https://github.com/gohugoio/hugo/commit/81975f847dc19c21c2321207645807771db97fab) [@bep](https://github.com/bep) [#8052](https://github.com/gohugoio/hugo/issues/8052)
* hugolib/paths: Fix typo [ce96895d](https://github.com/gohugoio/hugo/commit/ce96895debb67df20ae24fb5f0f04b98a30cc6cc) [@mayocream](https://github.com/mayocream) 
* Fix minor typos [04b89857](https://github.com/gohugoio/hugo/commit/04b89857e104ac7dcbf9fc65d8d4f1a1178123e6) [@phil-davis](https://github.com/phil-davis) 
* Fix BenchmarkMergeByLanguage [21fa1e86](https://github.com/gohugoio/hugo/commit/21fa1e86f2aa929fb0983a0cc3dc4e271ea1cc54) [@bep](https://github.com/bep) [#7914](https://github.com/gohugoio/hugo/issues/7914)
* Fix RelURL and AbsURL when path starts with language [aebfe156](https://github.com/gohugoio/hugo/commit/aebfe156fb2f27057e61b2e50c7576e6b06dab58) [@ivan-meridianbanc-com](https://github.com/ivan-meridianbanc-com) 





