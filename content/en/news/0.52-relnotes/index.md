
---
date: 2018-11-28
title: "And Now: Hugo 0.52"
description: "Configurable file caches, inline shortcodes, and more ..."
categories: ["Releases"]
---

The two big new items in this release is [Inline Shortcodes](https://gohugo.io//templates/shortcode-templates/#inline-shortcodes) and [Consolidated File Caches](https://gohugo.io/getting-started/configuration/#configure-file-caches). In Hugo we really care about build speed, and caching is important. With this release, you get much better control over your cache configuration, which is especially useful when building on a Continuous Integration server (Netlify, CircleCI or similar). Inline Shortcodes was implemented to help the Bootstrap project [move their documentation site](https://github.com/twbs/bootstrap/issues/24475#issuecomment-441238128) to Hugo. Note that this feature is disabled by default. To enable, set `enableInlineShortcodes = true` in your site config. Worth mentioning is also the new `param` shortcode, which looks up the param in page front matter with the site's parameter as a fall back.

This release represents **33 contributions by 7 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@emirb](https://github.com/emirb), and [@allizad](https://github.com/allizad) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **10 contributions by 4 contributors**. A special thanks to [@budparr](https://github.com/budparr), [@bep](https://github.com/bep), [@allizad](https://github.com/allizad), and [@funkydan2](https://github.com/funkydan2) for their work on the documentation site.

Hugo now has:

* 30595+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 270+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Add tests [ed698e94](https://github.com/gohugoio/hugo/commit/ed698e94c12c05bfc392eaca4f0c8442eac64906) [@moorereason](https://github.com/moorereason)
* Regenerate templates [89e2716d](https://github.com/gohugoio/hugo/commit/89e2716d290708ccde0a6f65504c1650c2f41b3d) [@bep](https://github.com/bep)
* Add "param" shortcode [f37c5a25](https://github.com/gohugoio/hugo/commit/f37c5a25676db89c0e804ccaac69bb392758192b) [@bep](https://github.com/bep) [#4010](https://github.com/gohugoio/hugo/issues/4010)
* Add float64 support to where [112461fd](https://github.com/gohugoio/hugo/commit/112461fded0d7970817ce7bf476c4763922ad314) [@moorereason](https://github.com/moorereason) [#5466](https://github.com/gohugoio/hugo/issues/5466)

### Core

* Fall back to title in ByLinkTitle sort [a9a93d08](https://github.com/gohugoio/hugo/commit/a9a93d082d8640684b7fd0076c64ea808ea7f762) [@bep](https://github.com/bep) [#4953](https://github.com/gohugoio/hugo/issues/4953)
* Improve nil handling in IsDescendant and IsAncestor [b09a4033](https://github.com/gohugoio/hugo/commit/b09a40333f382cc1034d2eda856230258ab6b8cc) [@bep](https://github.com/bep) [#5461](https://github.com/gohugoio/hugo/issues/5461)

### Other

* Remove duplicate mapstructure depdendency [7e75aeca](https://github.com/gohugoio/hugo/commit/7e75aeca80aead50d64902d2ff47e4ad4d013352) [@bep](https://github.com/bep)
* Add dependency list to README [e14e0b19](https://github.com/gohugoio/hugo/commit/e14e0b192f39812e3c3d5202d34ee907021412bb) [@bep](https://github.com/bep)
* Document inline shortcodes [aded0f25](https://github.com/gohugoio/hugo/commit/aded0f25fd23a78804b10e127aebe0e4b6fed2ac) [@bep](https://github.com/bep) [#4011](https://github.com/gohugoio/hugo/issues/4011)
* Add inline shortcode support [bc337e6a](https://github.com/gohugoio/hugo/commit/bc337e6ab5a75f1f1bfe3a83f3786d0afdb6346c) [@bep](https://github.com/bep) [#4011](https://github.com/gohugoio/hugo/issues/4011)
* Include drafts in convert command [dcfeed35](https://github.com/gohugoio/hugo/commit/dcfeed35c6e14c1ce593d23be9d2b89c66ce9bee) [@bep](https://github.com/bep) [#5457](https://github.com/gohugoio/hugo/issues/5457)
* Handle themes in the new file cache (for images, assets) [f9b4eb4f](https://github.com/gohugoio/hugo/commit/f9b4eb4f3968d32f45e0168c854e6b0c7f3a90b0) [@bep](https://github.com/bep) [#5460](https://github.com/gohugoio/hugo/issues/5460)
* Add tests for permalink on Resource with baseURL with path [12742bac](https://github.com/gohugoio/hugo/commit/12742bac71c65d65dc56548b643debda94757aee) [@bep](https://github.com/bep) [#5226](https://github.com/gohugoio/hugo/issues/5226)
* Add a comment about file mode for new files [fabf026f](https://github.com/gohugoio/hugo/commit/fabf026f4937bf6fbbb944aa7d6e721839ae4c92) [@bep](https://github.com/bep) [#5434](https://github.com/gohugoio/hugo/issues/5434)
* Add a :project placeholder [94f0f7e5](https://github.com/gohugoio/hugo/commit/94f0f7e59788e802e706a55cac0d52a9e70ff745) [@bep](https://github.com/bep) [#5439](https://github.com/gohugoio/hugo/issues/5439)
* Add a cache prune func [3c29c5af](https://github.com/gohugoio/hugo/commit/3c29c5af8ee865ef20741f576088e031e940c3d2) [@bep](https://github.com/bep) [#5439](https://github.com/gohugoio/hugo/issues/5439)
* Add a filecache root dir [33502667](https://github.com/gohugoio/hugo/commit/33502667fbacf57167ede66df8f13e308a4a9aec) [@bep](https://github.com/bep)
* Use time.Duration for maxAge [d3489eba](https://github.com/gohugoio/hugo/commit/d3489eba5dfc0ecdc032016d9db0746213dd5f0e) [@bep](https://github.com/bep) [#5438](https://github.com/gohugoio/hugo/issues/5438)
* Split implementation and config into separate files [17d7ecde](https://github.com/gohugoio/hugo/commit/17d7ecde2b261d2ab29049d12361b66504e3f995) [@bep](https://github.com/bep)
* Update to LibSASS 3.5.5 [e4b25728](https://github.com/gohugoio/hugo/commit/e4b2572880550a997d51dab3b198dac1fd642690) [@bep](https://github.com/bep) [#5432](https://github.com/gohugoio/hugo/issues/5432)[#5435](https://github.com/gohugoio/hugo/issues/5435)
* More spelling corrections [782dd158](https://github.com/gohugoio/hugo/commit/782dd15858128d8dfe78970c86e543b6590a004c) [@bep](https://github.com/bep)
* Spelling corrections [aff9c091](https://github.com/gohugoio/hugo/commit/aff9c091669a022b59f493c9dccf72be29511299) [@bep](https://github.com/bep)
* Remove appveyor [fdd4a768](https://github.com/gohugoio/hugo/commit/fdd4a768f053b21271d4520bf0d43baf62d516da) [@bep](https://github.com/bep)
* Document the new file cache [abeeff13](https://github.com/gohugoio/hugo/commit/abeeff1325267f8d8f1f66f0ec4ed175ffc140ad) [@bep](https://github.com/bep) [#5404](https://github.com/gohugoio/hugo/issues/5404)
* Add a consolidated file cache [f7aeaa61](https://github.com/gohugoio/hugo/commit/f7aeaa61291dd75f92901bcbeecc7fce07a28dec) [@bep](https://github.com/bep) [#5404](https://github.com/gohugoio/hugo/issues/5404)
* Add Windows build config to Travis [7d78a2af](https://github.com/gohugoio/hugo/commit/7d78a2afd3c4a6c4af77a4ddcbd2a82f15986048) [@emirb](https://github.com/emirb)
* Add Elasticsearch/bonsai.io to services doc. [c0b3a1af](https://github.com/gohugoio/hugo/commit/c0b3a1af0354e3aa9979cc00ae8630d7f0be63dc) [@allizad](https://github.com/allizad)

## Fixes

### Templates

* Fix whitespace issue [aba2647c](https://github.com/gohugoio/hugo/commit/aba2647c152ffff927f42523b77ee6651630cd67) [@max-arnold](https://github.com/max-arnold)
* Fix test to pass with gccgo [a8cb1b07](https://github.com/gohugoio/hugo/commit/a8cb1b07b4cf7fcf0e949657cb03c1a4838f975e) [@ianlancetaylor](https://github.com/ianlancetaylor)

### Other

* Fix handling of commented out front matter [7540a628](https://github.com/gohugoio/hugo/commit/7540a62834d4465af8936967e430a9e05a1e1359) [@bep](https://github.com/bep) [#5478](https://github.com/gohugoio/hugo/issues/5478)
* Fix when only shortcode and then summary [94ab125b](https://github.com/gohugoio/hugo/commit/94ab125b27a29a65e5ea45efd99dd247084b4c37) [@bep](https://github.com/bep) [#5464](https://github.com/gohugoio/hugo/issues/5464)
* Fix ignored --config flag with 'new' command [e82b2dc8](https://github.com/gohugoio/hugo/commit/e82b2dc8c1628f2da33e5fb0bae1b03e0594ad2c) [@krisbudhram](https://github.com/krisbudhram)
* Fix Permalink for resource, baseURL with path and canonifyURLs set [5df2b79d](https://github.com/gohugoio/hugo/commit/5df2b79dd2734e9a00ed1692328f58c385676468) [@bep](https://github.com/bep) [#5226](https://github.com/gohugoio/hugo/issues/5226)
