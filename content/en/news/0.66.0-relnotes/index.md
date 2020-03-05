
---
date: 2020-03-03
title: "Hugo 0.66.0: PostCSS Edition"
description: "Native inline, recursive import support in PostCSS/Tailwind, \"dependency-less\" builds, and more …"
categories: ["Releases"]
---

This release adds [inline `@import`](/hugo-pipes/postcss/#options) support to `resources.PostCSS`, with imports relative to Hugo's virtual, composable file system. Another useful addition is the new `build` [configuration section](/getting-started/configuration/#configure-build). As an example in `config.toml`:

```toml
[build]
  useResourceCacheWhen = "always"
```

The above will tell Hugo to _always_ use the cached build resources inside `resources/_gen` for the build steps requiring a non-standard dependency (PostCSS and SCSS/SASS). Valid values are `never`, `always` and `fallback` (default).


This release represents **27 contributions by 8 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@carlmjohnson](https://github.com/carlmjohnson), and [@sams96](https://github.com/sams96) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **8 contributions by 5 contributors**. A special thanks to [@bep](https://github.com/bep), [@nantipov](https://github.com/nantipov), [@regisphilibert](https://github.com/regisphilibert), and [@inwardmovement](https://github.com/inwardmovement) for their work on the documentation site.


Hugo now has:

* 41984+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 299+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Change error message on missing resource [d7798906](https://github.com/gohugoio/hugo/commit/d7798906d8e152a5d33f76ed0362628da8dd2c35) [@sams96](https://github.com/sams96) [#6942](https://github.com/gohugoio/hugo/issues/6942)
* Add math.Sqrt [d184e505](https://github.com/gohugoio/hugo/commit/d184e5059c72c15df055192b01da0fd8c5b0fc5c) [@StarsoftAnalysis](https://github.com/StarsoftAnalysis) [#6941](https://github.com/gohugoio/hugo/issues/6941)

### Other

* Skip some tests on CircleCI [6a34f88d](https://github.com/gohugoio/hugo/commit/6a34f88dcc1ac229247decc008471d7449d6d316) [@bep](https://github.com/bep) 
* {{ in }} should work with html.Template type [ae383f04](https://github.com/gohugoio/hugo/commit/ae383f04c806687cdae184d6138bcf51edbffcb2) [@carlmjohnson](https://github.com/carlmjohnson) [#7002](https://github.com/gohugoio/hugo/issues/7002)
* Regen CLI docs [ee31e61f](https://github.com/gohugoio/hugo/commit/ee31e61fb06bb6e26c9d66d78d8763aabd19e11d) [@bep](https://github.com/bep) 
* Add --all flag to hugo mod clean [760a87a4](https://github.com/gohugoio/hugo/commit/760a87a45a0a3e6a581851e5cf4fe440e9a8c655) [@bep](https://github.com/bep) 
* Add build.UseResourceCacheWhen [3d3fa5c3](https://github.com/gohugoio/hugo/commit/3d3fa5c3fe5ee0c9df59d682ee0acaba71a06ae1) [@bep](https://github.com/bep) [#6993](https://github.com/gohugoio/hugo/issues/6993)
* Update dependency list in README.md [ee3d0213](https://github.com/gohugoio/hugo/commit/ee3d02134d9b46b10e5a0403c9986ee1833ae6c1) [@anthonyfok](https://github.com/anthonyfok) 
* Add full filename to image when processing fails [305ce1c9](https://github.com/gohugoio/hugo/commit/305ce1c9ec746d3b8f6c9306b7014bfd621478a5) [@bep](https://github.com/bep) [#7000](https://github.com/gohugoio/hugo/issues/7000)
* Update dependency list in README [449deb7f](https://github.com/gohugoio/hugo/commit/449deb7f9ce089236f8328dd4fa585bea6e9bfde) [@anthonyfok](https://github.com/anthonyfok) 
* Add basic @import support to resources.PostCSS [b66d38c4](https://github.com/gohugoio/hugo/commit/b66d38c41939252649365822d9edb10cf5990617) [@bep](https://github.com/bep) [#6957](https://github.com/gohugoio/hugo/issues/6957)[#6961](https://github.com/gohugoio/hugo/issues/6961)
* Implement include/exclude filters for deploy [05a74eae](https://github.com/gohugoio/hugo/commit/05a74eaec0d944a4b29445c878a431cd6ae12277) [@vangent](https://github.com/vangent) [#6922](https://github.com/gohugoio/hugo/issues/6922)
* Update to Go 1.14 and 1.13.8 [33ae6210](https://github.com/gohugoio/hugo/commit/33ae62108325f703f1eaeabef1e8a80950229415) [@bep](https://github.com/bep) [#6958](https://github.com/gohugoio/hugo/issues/6958)
* Add hugo.IsProduction function [1352bc88](https://github.com/gohugoio/hugo/commit/1352bc880df4cd25eff65843973fcc0dd21b6304) [@hcwong](https://github.com/hcwong) [#6873](https://github.com/gohugoio/hugo/issues/6873)
* Apply missing go fmt [76b2afe6](https://github.com/gohugoio/hugo/commit/76b2afe642c37aedc7269b41d6fca5b78f467ce4) [@bep](https://github.com/bep) 

## Fixes

### Output

* Fix panic on no output formats [f4605303](https://github.com/gohugoio/hugo/commit/f46053034759c4f9790a79e0a146dbc1b426b1ff) [@bep](https://github.com/bep) [#6924](https://github.com/gohugoio/hugo/issues/6924)

### Core

* Fix error handling in page collector [3e9db2ad](https://github.com/gohugoio/hugo/commit/3e9db2ad951dbb1000cd0f8f25e4a95445046679) [@bep](https://github.com/bep) [#6988](https://github.com/gohugoio/hugo/issues/6988)
* Fix 2 Paginator.Pages taxonomy regressions [7ef5a4c8](https://github.com/gohugoio/hugo/commit/7ef5a4c83e4560bced3eee0ccf0e0db176146f44) [@bep](https://github.com/bep) [#6921](https://github.com/gohugoio/hugo/issues/6921)[#6918](https://github.com/gohugoio/hugo/issues/6918)
* Fix deletion of orphaned sections [a70bbd06](https://github.com/gohugoio/hugo/commit/a70bbd0696df3b0a6889650e48a07f8223151da4) [@bep](https://github.com/bep) [#6920](https://github.com/gohugoio/hugo/issues/6920)

### Other

* Fix ref/relref short lookup for pages in sub-folder [8947c3fa](https://github.com/gohugoio/hugo/commit/8947c3fa0beec021e14b3f8040857335e1ecd473) [@bep](https://github.com/bep) [#6952](https://github.com/gohugoio/hugo/issues/6952)
* Fix ref/relRef regression for relative refs from bundles [1746e8a9](https://github.com/gohugoio/hugo/commit/1746e8a9b2be46dcd6cecbb4bc90983a9c69b333) [@bep](https://github.com/bep) [#6952](https://github.com/gohugoio/hugo/issues/6952)
* Fix potential infinite recursion in server change detection [6f48146e](https://github.com/gohugoio/hugo/commit/6f48146e75e9877c4271ec239b763e6f3bc3babb) [@bep](https://github.com/bep) [#6986](https://github.com/gohugoio/hugo/issues/6986)
* Fix rebuild logic when editing template using a base template [b0d85032](https://github.com/gohugoio/hugo/commit/b0d850321e58a052ead25f7014b7851f63497601) [@bep](https://github.com/bep) [#6968](https://github.com/gohugoio/hugo/issues/6968)
* Fix panic when home page is drafted [0bd6356c](https://github.com/gohugoio/hugo/commit/0bd6356c6d2a2bac06d0c3705bf13a90cb7a2688) [@bep](https://github.com/bep) [#6927](https://github.com/gohugoio/hugo/issues/6927)
* Fix goldmark toc rendering [ca68abf0](https://github.com/gohugoio/hugo/commit/ca68abf0bc2fa003c2052143218f7b2ab195a46e) [@satotake](https://github.com/satotake) [#6736](https://github.com/gohugoio/hugo/issues/6736)[#6809](https://github.com/gohugoio/hugo/issues/6809)
* Fix crashes for 404 in IsAncestor etc. [a524124b](https://github.com/gohugoio/hugo/commit/a524124beb0e7ca226c207ea48a90cea2cbef76e) [@bep](https://github.com/bep) [#6931](https://github.com/gohugoio/hugo/issues/6931)
* Fix panic in 404.Parent [4c2a0de4](https://github.com/gohugoio/hugo/commit/4c2a0de412a850745ad32e580fcd65575192ca53) [@bep](https://github.com/bep) [#6924](https://github.com/gohugoio/hugo/issues/6924)





