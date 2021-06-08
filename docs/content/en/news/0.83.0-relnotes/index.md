
---
date: 2021-05-01
title: "Hugo 0.83: WebP Support!"
description: "WebP image encoding support, some important i18n fixes, and more."
categories: ["Releases"]
---

**Note:** If you use i18n, there is an unfortunate regression bug in this release (see [issue](https://github.com/gohugoio/hugo/issues/8492)). A patch release coming Sunday.


Hugo `0.83` finally brings [WebP](https://gohugo.io/content-management/image-processing/) image processing support. Note that you need the [extended version](https://gohugo.io/troubleshooting/faq/#i-get-tocss--this-feature-is-not-available-in-your-current-hugo-version) of Hugo to encode to WebP. If you want to target all Hugo versions, you may use a construct such as this:

```go-html-template
{{ $images := slice }}
{{ $images = $images | append ($img.Resize "300x") }}
{{ if hugo.IsExtended }}
  {{ $images = $images | append ($img.Resize "300x webp") }}
{{ end }}
```

Also worth highlighting:

* Some important language/i18n fixes (thanks to [@jmooring](https://github.com/jmooring) for helping out with these):
    * Fix multiple unknown language codes [7eb80a9e](https://github.com/gohugoio/hugo/commit/7eb80a9e6fcb6d31711effa20310cfefb7b23c1b) [@bep](https://github.com/bep) [#7838](https://github.com/gohugoio/hugo/issues/7838)
    * Improve plural handling of floats [eebde0c2](https://github.com/gohugoio/hugo/commit/eebde0c2ac4964e91d26d8b0cf0ac43afcfd207f) [@bep](https://github.com/bep) [#8464](https://github.com/gohugoio/hugo/issues/8464)
    * Revise the plural implementation [537c905e](https://github.com/gohugoio/hugo/commit/537c905ec103dc5adaf8a1b2ccdef5da7cc660fd) [@bep](https://github.com/bep) [#8454](https://github.com/gohugoio/hugo/issues/8454)[#7822](https://github.com/gohugoio/hugo/issues/7822)
* You can now use slice syntax in the sections permalinks config[2dc222ce](https://github.com/gohugoio/hugo/commit/2dc222cec4460595af8569165d1c498bb45aac84) [@bep](https://github.com/bep) [#8363](https://github.com/gohugoio/hugo/issues/8363).

This release represents **61 contributions by 9 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@dependabot[bot]](https://github.com/apps/dependabot), [@jmooring](https://github.com/jmooring), and [@anthonyfok](https://github.com/anthonyfok) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **10 contributions by 5 contributors**. A special thanks to [@lupsa](https://github.com/lupsa), [@jmooring](https://github.com/jmooring), [@bep](https://github.com/bep), and [@arhuman](https://github.com/arhuman) for their work on the documentation site.


Hugo now has:

* 51594+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 432+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 370+ [themes](http://themes.gohugo.io/)

## Notes

* We have updated ESBUild to v0.11.16. There are no breaking changes on the API side, but you may want to read the release upstream release notes: https://github.com/evanw/esbuild/releases/tag/v0.10.0 https://github.com/evanw/esbuild/releases/tag/v0.11.0

## Enhancements

### Templates

* Remove the FuzzMarkdownify func for now [5656a908](https://github.com/gohugoio/hugo/commit/5656a908d837f2aa21837d39712b8ab4aa6db842) [@bep](https://github.com/bep) 

### Output

* Make the shortcode template lookup for output formats stable [0d86a32d](https://github.com/gohugoio/hugo/commit/0d86a32d8f3031e2124c8005b680b597f3c0e558) [@bep](https://github.com/bep) [#7774](https://github.com/gohugoio/hugo/issues/7774)
* Only output mediaType once in docshelper JSON [7b4ade56](https://github.com/gohugoio/hugo/commit/7b4ade56dd50d89a91760fc5ef8e2f151874de96) [@bep](https://github.com/bep) [#8379](https://github.com/gohugoio/hugo/issues/8379)

### Other

* Regenerate docs helper [a9b52b41](https://github.com/gohugoio/hugo/commit/a9b52b41758d20ae4c10b71721b22175395c69e9) [@bep](https://github.com/bep) 
* Regenerate CLI docs [b073a1c9](https://github.com/gohugoio/hugo/commit/b073a1c9723980eeb58717884006148dfc0e0c8e) [@bep](https://github.com/bep) 
* Remove all dates from gendoc [4227cc1b](https://github.com/gohugoio/hugo/commit/4227cc1bd308d1ef1ea151c86f72f537b5e77b1d) [@bep](https://github.com/bep) 
* Update getkin/kin-openapi v0.60.0 => v0.61. [3cc4fdd6](https://github.com/gohugoio/hugo/commit/3cc4fdd6f358263ffde33ccbf61546f073979e32) [@bep](https://github.com/bep) 
* Update github.com/evanw/esbuild v0.11.14 => v0.11.16 [78c1a6a7](https://github.com/gohugoio/hugo/commit/78c1a6a7c6e14f006854ee97ec561abdcf6203fc) [@bep](https://github.com/bep) 
* Remove .Site.Authors from embedded templates [f6745ad3](https://github.com/gohugoio/hugo/commit/f6745ad3588a7b3aaae228fec18fe0027affd566) [@jmooring](https://github.com/jmooring) [#4458](https://github.com/gohugoio/hugo/issues/4458)
* Don't treat a NotFound response for Delete as a fatal error. [f523e9f0](https://github.com/gohugoio/hugo/commit/f523e9f0fd0e0b0ce75879532caa834742297d16) [@vangent](https://github.com/vangent) 
* Switch to deb packages of nodejs and python3-pygments [63cd05ce](https://github.com/gohugoio/hugo/commit/63cd05ce5ae308c496b848f6b11bcb3fdbdf5cb2) [@anthonyfok](https://github.com/anthonyfok) 
* Install bin/node from node/14/stable [902535ef](https://github.com/gohugoio/hugo/commit/902535ef11fce449b377896ab7498c4799beb9ce) [@anthonyfok](https://github.com/anthonyfok) 
* bump github.com/getkin/kin-openapi from 0.55.0 to 0.60.0 [70aebba0](https://github.com/gohugoio/hugo/commit/70aebba04d801fe6a3784394d25c433ffeb6d123) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.11.13 to 0.11.14 [3e3b7d44](https://github.com/gohugoio/hugo/commit/3e3b7d4474ea97a1990f303482a12f0c3031bd07) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update to Chroma v0.9.1 [048418ba](https://github.com/gohugoio/hugo/commit/048418ba749d02eb3dde9d6895cedef2adaefefd) [@caarlos0](https://github.com/caarlos0) 
* Improve plural handling of floats [eebde0c2](https://github.com/gohugoio/hugo/commit/eebde0c2ac4964e91d26d8b0cf0ac43afcfd207f) [@bep](https://github.com/bep) [#8464](https://github.com/gohugoio/hugo/issues/8464)
* bump github.com/evanw/esbuild from 0.11.12 to 0.11.13 [65c502cc](https://github.com/gohugoio/hugo/commit/65c502cc8110e49540cbe2b49ecd5a8ede9e67a1) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Revise the plural implementation [537c905e](https://github.com/gohugoio/hugo/commit/537c905ec103dc5adaf8a1b2ccdef5da7cc660fd) [@bep](https://github.com/bep) [#8454](https://github.com/gohugoio/hugo/issues/8454)[#7822](https://github.com/gohugoio/hugo/issues/7822)
* Update to "base: core20" [243951eb](https://github.com/gohugoio/hugo/commit/243951ebe9715d3da3968e96e6f60dcd53e25d92) [@anthonyfok](https://github.com/anthonyfok) 
* bump github.com/frankban/quicktest from 1.11.3 to 1.12.0 [fe2ee028](https://github.com/gohugoio/hugo/commit/fe2ee028024836695c99e28595393588e3930136) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump google.golang.org/api from 0.44.0 to 0.45.0 [316d65cd](https://github.com/gohugoio/hugo/commit/316d65cd7049d60b0d5ac0080a87236198e74fc9) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.37.11 to 1.38.23 [b95229ab](https://github.com/gohugoio/hugo/commit/b95229ab49ac2126aefe7802392ef34fdd021c3b) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Correct function name in comment [0551df09](https://github.com/gohugoio/hugo/commit/0551df090e6b2a391941bf7383b79c2dbc11d416) [@xhit](https://github.com/xhit) 
* Upgraded github.com/evanw/esbuild v0.11.0 => v0.11.12 [057e5a22](https://github.com/gohugoio/hugo/commit/057e5a22af937459082c3096ba3095b343d1a8bf) [@bep](https://github.com/bep) 
* Regen docs helper [fd96f65a](https://github.com/gohugoio/hugo/commit/fd96f65a3d7755e49b4a70fb276dfffcba4e541a) [@bep](https://github.com/bep) 
* bump github.com/tdewolff/minify/v2 from 2.9.15 to 2.9.16 [d3a64708](https://github.com/gohugoio/hugo/commit/d3a64708f49139552ca79a199a4cbf6544375443) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump golang.org/x/text from 0.3.5 to 0.3.6 [3b56244f](https://github.com/gohugoio/hugo/commit/3b56244f425a72c783bb58c30542aeb4b045acca) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Remove some unreachable code [f5d3d635](https://github.com/gohugoio/hugo/commit/f5d3d635e6b88d7c5d304b80f04e7b4361349fd6) [@bep](https://github.com/bep) 
* bump github.com/getkin/kin-openapi from 0.39.0 to 0.55.0 [0d3c42da](https://github.com/gohugoio/hugo/commit/0d3c42da56151325f16802b3b1a4105a21ce250e) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Some performance tweaks for the HTML elements collector [ef34dd8f](https://github.com/gohugoio/hugo/commit/ef34dd8f0e94e52ba6f1d5d607e4ac3ae98a7abb) [@bep](https://github.com/bep) 
* Exclude comment and doctype elements from writeStats [bc80022e](https://github.com/gohugoio/hugo/commit/bc80022e033a5462d1a9ce541f40a050994011cc) [@dirkolbrich](https://github.com/dirkolbrich) [#8396](https://github.com/gohugoio/hugo/issues/8396)[#8417](https://github.com/gohugoio/hugo/issues/8417)
* Merge branch 'release-0.82.1' [2bb9496c](https://github.com/gohugoio/hugo/commit/2bb9496ce29dfe90e8b3664ed8cf7f895011b2d4) [@bep](https://github.com/bep) 
* bump github.com/yuin/goldmark from 1.3.2 to 1.3.5 [3ddffd06](https://github.com/gohugoio/hugo/commit/3ddffd064dbacf62aa854b26ea8ddc5d15ba1ef8) [@jmooring](https://github.com/jmooring) [#8377](https://github.com/gohugoio/hugo/issues/8377)
* Remove duplicate references from release notes [6fc52d18](https://github.com/gohugoio/hugo/commit/6fc52d185a98b86c70b6ba862549cc6aae782691) [@jmooring](https://github.com/jmooring) [#8360](https://github.com/gohugoio/hugo/issues/8360)
* bump github.com/spf13/afero from 1.5.1 to 1.6.0 [73c3ae81](https://github.com/gohugoio/hugo/commit/73c3ae818a7fc78febff092ac74772a114a2cbd2) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/pelletier/go-toml from 1.8.1 to 1.9.0 [7ca118fd](https://github.com/gohugoio/hugo/commit/7ca118fdfd9f0d1c636ef5e266c9000a20099e03) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Add webp image encoding support [33d5f805](https://github.com/gohugoio/hugo/commit/33d5f805923eb50dfb309d024f6555c59a339846) [@bep](https://github.com/bep) [#5924](https://github.com/gohugoio/hugo/issues/5924)
* bump google.golang.org/api from 0.40.0 to 0.44.0 [509d39fa](https://github.com/gohugoio/hugo/commit/509d39fa6ddbba106c127b7923a41b0dcaea9381) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/nicksnyder/go-i18n/v2 from 2.1.1 to 2.1.2 [7725c41d](https://github.com/gohugoio/hugo/commit/7725c41d40b7009c2701a5ad3fa6bc9de57b88ee) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/rogpeppe/go-internal from 1.6.2 to 1.8.0 [5d36d801](https://github.com/gohugoio/hugo/commit/5d36d801534c0823697610fdb32e1eeb61f70e33) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Remove extraneous space from figure shortcode [9b34d42b](https://github.com/gohugoio/hugo/commit/9b34d42bb2ff05deaeeef63ff4b5b993f35f0451) [@jmooring](https://github.com/jmooring) [#8401](https://github.com/gohugoio/hugo/issues/8401)
* bump github.com/magefile/mage from 1.10.0 to 1.11.0 [c2d8f87c](https://github.com/gohugoio/hugo/commit/c2d8f87cfc1c4ae666fbb1fb5b8983d43492333f) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/google/go-cmp from 0.5.4 to 0.5.5 [cbc24661](https://github.com/gohugoio/hugo/commit/cbc246616e88729322dad70971eae18ef59dd5d4) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Disable broken pretty relative links feature [fa432b17](https://github.com/gohugoio/hugo/commit/fa432b17b349ed7e914af3625187e2c1dc2e243b) [@niklasfasching](https://github.com/niklasfasching) 
* Update go-org to v1.5.0 [0cd55c66](https://github.com/gohugoio/hugo/commit/0cd55c66d370559b66eea220626c4842efaf7039) [@niklasfasching](https://github.com/niklasfasching) 
* bump github.com/jdkato/prose from 1.2.0 to 1.2.1 [0d5cf256](https://github.com/gohugoio/hugo/commit/0d5cf256e4f2a5babcbcf7b49a6818869c3c0691) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/spf13/cobra from 1.1.1 to 1.1.3 [36527576](https://github.com/gohugoio/hugo/commit/36527576b30224dff2eae7f6c9f27eff807d5402) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Add complete dependency list in "hugo env -v" [9b83f45b](https://github.com/gohugoio/hugo/commit/9b83f45b6dcafa6e50df80a4786d6a36400a47fe) [@bep](https://github.com/bep) [#8400](https://github.com/gohugoio/hugo/issues/8400)
* Add hugo.IsExtended [7fdd2b95](https://github.com/gohugoio/hugo/commit/7fdd2b95e20f322b0a47f63ff1010a04f47ce67b) [@bep](https://github.com/bep) [#8399](https://github.com/gohugoio/hugo/issues/8399)
* Also test minified HTML in the element collector [3d5dbdcb](https://github.com/gohugoio/hugo/commit/3d5dbdcb1a11b059fc2f93ed6fadb9009bf72673) [@bep](https://github.com/bep) [#7567](https://github.com/gohugoio/hugo/issues/7567)
* Skip script, pre and textarea content when looking for HTML elements [8a308944](https://github.com/gohugoio/hugo/commit/8a308944e46f8c2aa054005d5aed89f2711f9c1d) [@bep](https://github.com/bep) [#7567](https://github.com/gohugoio/hugo/issues/7567)
* Add slice syntax to sections permalinks config [2dc222ce](https://github.com/gohugoio/hugo/commit/2dc222cec4460595af8569165d1c498bb45aac84) [@bep](https://github.com/bep) [#8363](https://github.com/gohugoio/hugo/issues/8363)
* Upgrade github.com/evanw/esbuild v0.9.6 => v0.11.0 [4d22ad58](https://github.com/gohugoio/hugo/commit/4d22ad580ec8c8e5e27cf4f5cce69b6828aa8501) [@bep](https://github.com/bep) 

## Fixes

### Templates

* Fix where on type mismatches [e4dc9a82](https://github.com/gohugoio/hugo/commit/e4dc9a82b557a417b1552c533b0df605c6ff1cc0) [@bep](https://github.com/bep) [#8353](https://github.com/gohugoio/hugo/issues/8353)

### Output

* Regression in media type suffix lookup [6e9d2bf0](https://github.com/gohugoio/hugo/commit/6e9d2bf0c936900f8f676d485098755b3f463373) [@bep](https://github.com/bep) [#8406](https://github.com/gohugoio/hugo/issues/8406)
* Regression in media type suffix lookup [e73f7a77](https://github.com/gohugoio/hugo/commit/e73f7a770dfb06f23d842d589bdd3d0fb53c7eed) [@bep](https://github.com/bep) [#8406](https://github.com/gohugoio/hugo/issues/8406)

### Other

* Fix multiple unknown language codes [7eb80a9e](https://github.com/gohugoio/hugo/commit/7eb80a9e6fcb6d31711effa20310cfefb7b23c1b) [@bep](https://github.com/bep) [#7838](https://github.com/gohugoio/hugo/issues/7838)
* Fix permalinks pattern detection for some of the sections variants [c13d3687](https://github.com/gohugoio/hugo/commit/c13d368746992eb39a33f065ca808e129baec4ef) [@bep](https://github.com/bep) [#8363](https://github.com/gohugoio/hugo/issues/8363)
* Fix Params case handling in where with slices of structs (e.g. Pages) [bca40cf0](https://github.com/gohugoio/hugo/commit/bca40cf0c9c7b75e6d5b4a9ac8b927eb17590c7e) [@bep](https://github.com/bep) [#7009](https://github.com/gohugoio/hugo/issues/7009)
* Fix typo in docshelper.go [7c7974b7](https://github.com/gohugoio/hugo/commit/7c7974b711879938eafc08a2ce242d0f00c8e9e6) [@jmooring](https://github.com/jmooring) [#8380](https://github.com/gohugoio/hugo/issues/8380)
* Try to fix the fuzz build [5e2f1289](https://github.com/gohugoio/hugo/commit/5e2f1289118dc1489fb782bf289298a05104eeaf) [@bep](https://github.com/bep) 





