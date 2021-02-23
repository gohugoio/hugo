
---
date: 2020-04-10
title: "Post Build Resource Transformations"
description: "Hugo 0.69.0 allows you to delay resource processing to after the build, the prime use case being removal of unused CSS."
categories: ["Releases"]
---

**It's Eeaster, a time for mysteries and puzzles.** And at first glance, this Hugo release looks a little mysterious. The core of if is a mind-twister:

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $css = $css | resources.PostCSS }}
{{ if hugo.IsProduction }}
{{ $css = $css | minify | fingerprint | resources.PostProcess }}
{{ end }}
<link href="{{ $css.RelPermalink }}" rel="stylesheet" />
```

The above uses the new [resources.PostProcess](https://gohugo.io/hugo-pipes/postprocess/) template function which tells Hugo to postpone the transformation of the Hugo Pipes chain to _after the build_, allowing the build steps to use the build output in `/public` as part of its processing.

The prime current use case for the above is CSS pruning in PostCSS. In simple cases you can use the templates as a base for the content filters, but that has its limitations and can be very hard to setup, especially in themed configurations. So we have added a new [writeStats](https://gohugo.io/getting-started/configuration/#configure-build) configuration that, when enabled, will write a file named `hugo_stats.json` to your project root with some aggregated data about the build, e.g. list of HTML entities published, to be used to do [CSS pruning](https://gohugo.io/hugo-pipes/postprocess/#css-purging-with-postcss). 

This release represents **20 contributions by 10 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@jaywilliams](https://github.com/jaywilliams), and [@satotake](https://github.com/satotake) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **14 contributions by 7 contributors**. A special thanks to [@bep](https://github.com/bep), [@coliff](https://github.com/coliff), [@dmgawel](https://github.com/dmgawel), and [@jasikpark](https://github.com/jasikpark) for their work on the documentation site.


Hugo now has:

* 43052+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 438+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 302+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Extend Jsonify to support options map [8568928a](https://github.com/gohugoio/hugo/commit/8568928aa8e82a6bd7de4555c3703d8835fbd25b) [@moorereason](https://github.com/moorereason) 
* Extend Jsonify to support optional indent parameter [1bc93021](https://github.com/gohugoio/hugo/commit/1bc93021e3dca6405628f6fdd2dc32cff9c9836c) [@moorereason](https://github.com/moorereason) [#5040](https://github.com/gohugoio/hugo/issues/5040)

### Other

* Regen docs helper [b7ff4dc2](https://github.com/gohugoio/hugo/commit/b7ff4dc23e6314fd09ee2c1e24cde96fc833164e) [@bep](https://github.com/bep) 
* Collect HTML elements during the build to use in PurgeCSS etc. [095bf64c](https://github.com/gohugoio/hugo/commit/095bf64c99f57efe083540a50e658808a0a1c32b) [@bep](https://github.com/bep) [#6999](https://github.com/gohugoio/hugo/issues/6999)
* Update to latest emoji package [7791a804](https://github.com/gohugoio/hugo/commit/7791a804e2179667617b3b145b0fe7eba17627a1) [@QuLogic](https://github.com/QuLogic) 
* Update hosting-on-aws-amplify.md [c774b230](https://github.com/gohugoio/hugo/commit/c774b230e941902675af081f118ea206a4f2a04e) [@Helicer](https://github.com/Helicer) 
* Add basic "post resource publish support" [2f721f8e](https://github.com/gohugoio/hugo/commit/2f721f8ec69c52202815cd1b543ca4bf535c0901) [@bep](https://github.com/bep) [#7146](https://github.com/gohugoio/hugo/issues/7146)
* Typo correction [7eba37ae](https://github.com/gohugoio/hugo/commit/7eba37ae9b8653be4fc21a0dbbc6f35ca5b9280e) [@fekete-robert](https://github.com/fekete-robert) 
* Use semver for min_version per recommendations [efc61d6f](https://github.com/gohugoio/hugo/commit/efc61d6f3b9f5fb294411ac1dc872b8fc5bdbacb) [@jaywilliams](https://github.com/jaywilliams) 
* Updateto gitmap v1.1.2 [4de3ecdc](https://github.com/gohugoio/hugo/commit/4de3ecdc2658ffd54d2b5073c5ff303b4bf29383) [@dragtor](https://github.com/dragtor) [#6985](https://github.com/gohugoio/hugo/issues/6985)
* Add data context to the key in ExecuteAsTemplate" [c9dc316a](https://github.com/gohugoio/hugo/commit/c9dc316ad160e78c9dff4e75313db4cac8ea6414) [@bep](https://github.com/bep) [#7064](https://github.com/gohugoio/hugo/issues/7064)

## Fixes

### Other

* Fix hugo mod vendor for regular file mounts [d8d6a25b](https://github.com/gohugoio/hugo/commit/d8d6a25b5755bedaf90261a1539dc37a2f05c3df) [@bep](https://github.com/bep) [#7140](https://github.com/gohugoio/hugo/issues/7140)
* Revert "Revert "common/herrors: Fix typos in comments"" [9f12be54](https://github.com/gohugoio/hugo/commit/9f12be54ee84f24efdf7c58f05867e8d0dea2ccb) [@bep](https://github.com/bep) 
* Fix typos in comments" [4437e918](https://github.com/gohugoio/hugo/commit/4437e918cdab1d84f2f184fe71e5dac14aa48897) [@bep](https://github.com/bep) 
* Fix typos in comments [1123711b](https://github.com/gohugoio/hugo/commit/1123711b0979b1647d7c486f67af7503afb11abb) [@rnazmo](https://github.com/rnazmo) 
* Fix TrimShortHTML [9c998753](https://github.com/gohugoio/hugo/commit/9c9987535f98714c8a4ec98903f54233735ef0e4) [@satotake](https://github.com/satotake) [#7081](https://github.com/gohugoio/hugo/issues/7081)
* Fix IsDescendant/IsAncestor for overlapping section names [4a39564e](https://github.com/gohugoio/hugo/commit/4a39564efe7b02a685598ae9dbae95e2326c0230) [@bep](https://github.com/bep) [#7096](https://github.com/gohugoio/hugo/issues/7096)
* fix typo in getting started [b6e097cf](https://github.com/gohugoio/hugo/commit/b6e097cfe65ecd1d47c805969082e6805563612b) [@matrixise](https://github.com/matrixise) 
* Fix _build.list.local logic [523d5194](https://github.com/gohugoio/hugo/commit/523d51948fc20e2afb4721b43203c5ab696ae220) [@bep](https://github.com/bep) [#7089](https://github.com/gohugoio/hugo/issues/7089)
* Fix cache reset for a page's collections on server live reload [cfa73050](https://github.com/gohugoio/hugo/commit/cfa73050a49b2646fe3557cefa0ed31989b0eeeb) [@bep](https://github.com/bep) [#7085](https://github.com/gohugoio/hugo/issues/7085)





