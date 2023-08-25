
---
date: 2020-10-30
title: "Hugo 0.77.0: Hugo Modules Improvements and More "
description: "New Replacements config option for simpler development workflows, ignore errors from getJSON, localized dates, and more."
categories: ["Releases"]
---

Hugo `0.77.0` is a small, but useful release. Some notable updates are:

* **time.AsTime** accepts an optional location as second parameter, allowing timezone aware printing of dates.
* You can now build with `go install -tags nodeploy` if you don't need the **`hugo deploy`** feature.
* Remote **`getJSON`** errors can now be ignored by adding `ignoreErrors = ["error-remote-getjson"]` to your site config.

There are also several useful **[Hugo Modules](https://gohugo.io/hugo-modules/)** enhancements:

* We have added `Replacements` to the [Module Configuration](https://gohugo.io/hugo-modules/configuration/#module-configuration-top-level). This should enable a much simpler developer workflow, simpler to set up preview sites for your remote theme etc, as you now can do `env HUGO_MODULE_REPLACEMENTS="github.com/bep/myprettytheme -> ../.." hugo` and similar.
* The module `Path` for local modules can now be absolute for imports defined in the project.

This release represents **38 contributions by 11 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), and [@anthonyfok](https://github.com/anthonyfok) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **3 contributions by 3 contributors**.

Hugo now has:

* 47530+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 438+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 361+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Refactor time.AsTime location implementation [807db97a](https://github.com/gohugoio/hugo/commit/807db97af83ff61b022cbc8af80b9dc9cdb8dd43) [@moorereason](https://github.com/moorereason) 
* Update Hugo time to support optional [LOCATION] parameter [26eeb291](https://github.com/gohugoio/hugo/commit/26eeb2914720929d2d778f14d6a4bf737014e9e3) [@virgofx](https://github.com/virgofx) 
* Improve layout path construction [acfa1538](https://github.com/gohugoio/hugo/commit/acfa153863d6ff2acf17ffb4395e05d102229905) [@moorereason](https://github.com/moorereason) 
* Test all lookup permutations in TestLayout [78b26d53](https://github.com/gohugoio/hugo/commit/78b26d538c716d463b30c23de7df5eaa4d5504fd) [@moorereason](https://github.com/moorereason) 
* Reformat TestLayout table [28179bd5](https://github.com/gohugoio/hugo/commit/28179bd55619847f46ca0ffd316ef52fc9c96f1e) [@moorereason](https://github.com/moorereason) 

### Other

* Allow absolute paths for project imports [beabc8d9](https://github.com/gohugoio/hugo/commit/beabc8d998249ecc5dd522d696dc6233a29131c2) [@bep](https://github.com/bep) [#7910](https://github.com/gohugoio/hugo/issues/7910)
* Regen docs helper [332b65e4](https://github.com/gohugoio/hugo/commit/332b65e4ccb6ac0d606de2a1b23f5189c72542be) [@bep](https://github.com/bep) 
* Add module.replacements [173187e2](https://github.com/gohugoio/hugo/commit/173187e2633f3fc037c83e1e3de2902ae3c93b92) [@bep](https://github.com/bep) [#7904](https://github.com/gohugoio/hugo/issues/7904)[#7908](https://github.com/gohugoio/hugo/issues/7908)
* Do not call CDN service invalidation when executing a dry run deployment [56a34350](https://github.com/gohugoio/hugo/commit/56a343507ca28254edb891bc1c21b6c8ca017982) [@zemanel](https://github.com/zemanel) [#7884](https://github.com/gohugoio/hugo/issues/7884)
* Pass editor arguments from newContentEditor correctly [d48a98c4](https://github.com/gohugoio/hugo/commit/d48a98c477a818d28008d9771050d2681e63e880) [@bhavin192](https://github.com/bhavin192) 
* Bump github.com/spf13/cobra from 0.0.7 to 1.1.1 [3261678f](https://github.com/gohugoio/hugo/commit/3261678f63fd66810db77ccaf9a0c0e426be5380) [@anthonyfok](https://github.com/anthonyfok) 
* Allow optional "nodeploy" tag to exclude deploy command from bin [f465c5c3](https://github.com/gohugoio/hugo/commit/f465c5c3079261eb7fa513e2d2793851b9c52b83) [@emhagman](https://github.com/emhagman) [#7826](https://github.com/gohugoio/hugo/issues/7826)
* Allow cascade _target to work with non toml fm [3400aff2](https://github.com/gohugoio/hugo/commit/3400aff2588cbf9dd4629c05537d16b019d0fdf5) [@gwatts](https://github.com/gwatts) [#7874](https://github.com/gohugoio/hugo/issues/7874)
* Allow getJSON errors to be ignored [fdfa4a5f](https://github.com/gohugoio/hugo/commit/fdfa4a5fe62232f65f1dd8d6fe0c500374228788) [@bep](https://github.com/bep) [#7866](https://github.com/gohugoio/hugo/issues/7866)
* bump github.com/evanw/esbuild from 0.7.15 to 0.7.18 [8cbe2bbf](https://github.com/gohugoio/hugo/commit/8cbe2bbfad6aa4de267921e24e166d4addf47040) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Revert "Add benchmark for building docs site" [b886fa46](https://github.com/gohugoio/hugo/commit/b886fa46bb92916152476cfac45c7a5ee5e5820a) [@bep](https://github.com/bep) 
* Avoid making unnecessary allocation [14bce18a](https://github.com/gohugoio/hugo/commit/14bce18a6c5aca8cb3e70a74d5045ca8b2358fee) [@moorereason](https://github.com/moorereason) 
* Add benchmark for building docs site [837e084b](https://github.com/gohugoio/hugo/commit/837e084bbe53e9e2e6cd471d2a3daf273a874d92) [@moorereason](https://github.com/moorereason) 
* Always show page number when 5 pages or less [08e4f9ff](https://github.com/gohugoio/hugo/commit/08e4f9ff9cc448d5fea9b8a62a23aed8aad0d047) [@moorereason](https://github.com/moorereason) [#7523](https://github.com/gohugoio/hugo/issues/7523)
* bump github.com/frankban/quicktest from 1.11.0 to 1.11.1 [f033d9f0](https://github.com/gohugoio/hugo/commit/f033d9f01d13d8cd08205ccfaa09919ed15dca77) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.7.14 to 0.7.15 [59fe2794](https://github.com/gohugoio/hugo/commit/59fe279424c66ac6a89cafee01a5b2e34dbcc1fb) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Merge branch 'release-0.76.5' [62119022](https://github.com/gohugoio/hugo/commit/62119022d1be41e423ef3bcf467a671ce6c4f7dd) [@bep](https://github.com/bep) 
* Render aliases even if render=link [79a022a1](https://github.com/gohugoio/hugo/commit/79a022a15c5f39b8ae87a94665f14bf1797b605c) [@bep](https://github.com/bep) [#7832](https://github.com/gohugoio/hugo/issues/7832)
* Render aliases even if render=link [ead5799f](https://github.com/gohugoio/hugo/commit/ead5799f7ea837fb2ca1879a6d37ba364e53827f) [@bep](https://github.com/bep) [#7832](https://github.com/gohugoio/hugo/issues/7832)
* bump github.com/spf13/afero from 1.4.0 to 1.4.1 [d57be113](https://github.com/gohugoio/hugo/commit/d57be113243be4b76310d4476fbb7525d1452658) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.7.9 to 0.7.14 [d0705966](https://github.com/gohugoio/hugo/commit/d070596694a3edbf42fc315bb326505aa39fce90) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update to Go 1.15 and Alpine 3.12 [f5ea359d](https://github.com/gohugoio/hugo/commit/f5ea359dd34bf59a2944f1d9667838202af13c93) [@ducksecops](https://github.com/ducksecops) 
* Install postcss v8 explicitly as it is now a peer dependency [e9a7ebaf](https://github.com/gohugoio/hugo/commit/e9a7ebaf67a63ffe5e64c3b3aaefe66feb7f1868) [@anthonyfok](https://github.com/anthonyfok) 
* Merge branch 'release-0.76.3' [49972d07](https://github.com/gohugoio/hugo/commit/49972d07925604fea45afe1ace7b5dcc6efc30bf) [@bep](https://github.com/bep) 
* Add merge helper [c98132e3](https://github.com/gohugoio/hugo/commit/c98132e30e01a9638e61bd888c769d30e4e43ad5) [@bep](https://github.com/bep) 
* Add workaround for known language, but missing plural rule error [33e9d79b](https://github.com/gohugoio/hugo/commit/33e9d79b78b32d0cc19693ab3c29ba9941d80f8f) [@bep](https://github.com/bep) [#7798](https://github.com/gohugoio/hugo/issues/7798)
* Update to  github.com/tdewolff/minify v2.9.4" [6dd60fca](https://github.com/gohugoio/hugo/commit/6dd60fca73ff96b48064bb8c6586631a2370ffc6) [@bep](https://github.com/bep) [#7792](https://github.com/gohugoio/hugo/issues/7792)

## Fixes

### Templates

* Fix reflection bug in merge [6d95dc9d](https://github.com/gohugoio/hugo/commit/6d95dc9d74681cba53b46e79c6e1d58d27fcdfb0) [@moorereason](https://github.com/moorereason) [#7899](https://github.com/gohugoio/hugo/issues/7899)

### Other

* Fix setting HUGO_MODULE_PROXY etc. via env vars [8a1c637c](https://github.com/gohugoio/hugo/commit/8a1c637c4494751046142e0ef345fce38fc1431b) [@bep](https://github.com/bep) [#7903](https://github.com/gohugoio/hugo/issues/7903)
* Fix for language code case issue with pt-br etc. [50682043](https://github.com/gohugoio/hugo/commit/506820435cacb39ce7bb1835f46a15e913b95828) [@bep](https://github.com/bep) [#7804](https://github.com/gohugoio/hugo/issues/7804)
* Fix for bare TOML keys [fc6abc39](https://github.com/gohugoio/hugo/commit/fc6abc39c75c152780151c35bc95b12bee01b09c) [@bep](https://github.com/bep) 
* Fix i18n .Count regression [f9e798e8](https://github.com/gohugoio/hugo/commit/f9e798e8c4234bd60277e3cb10663ba254d4ecb7) [@bep](https://github.com/bep) [#7787](https://github.com/gohugoio/hugo/issues/7787)
* Fix typo in 0.76.0 release note [ee56efff](https://github.com/gohugoio/hugo/commit/ee56efffcb3f81120b0d3e0297b4fb5966124354) [@digitalcraftsman](https://github.com/digitalcraftsman) 
