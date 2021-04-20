
---
date: 2020-10-06
title: "Multiple Cascades With Page Filters"
description: "Hugo 0.76.0 brings multiple cascade blocks per page with filters for path, kind and language."
categories: ["Releases"]
---

In **Hugo 0.76.0** you can now have a list of [cascade](https://gohugo.io/content-management/front-matter#front-matter-cascade) blocks per page and a new `_target` keyword where you can select which pages to _cascade_ upon using [Glob](https://github.com/gobwas/glob) patterns for a `Page`'s `Kind`, `Lang` and/or `Path`:

```toml
title ="Blog"
[[cascade]]
background = "yosemite.jpg"
[cascade._target]
path="/blog/**"
lang="en"
kind="page"
[[cascade]]
background = "goldenbridge.jpg"
[cascade._target]
kind="section"
```

Tasks that were earlier hard/borderline impossible to do are now simple. One common example would to apply a different template set to nested sections; you can now apply a custom `Type` to these sections using  `path="/blog/*/**"` and similar.

A related improvement is that the [build option](https://gohugo.io/content-management/build-options/#readout) `render` is now an enum. In addition to turning on/off rendering of a given page you can tell Hugo to not render, but you want to preserve the `.Permalink`, useful for SPA applications.

This release represents **35 contributions by 8 contributors** to the main Hugo code base. A big shoutout to [@bep](https://github.com/bep), [@ai](https://github.com/ai), and [@jmooring](https://github.com/jmooring) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **11 contributions by 6 contributors**. A special thanks to [@amdw](https://github.com/amdw), [@davidsneighbour](https://github.com/davidsneighbour), [@samrobbins85](https://github.com/samrobbins85), and [@yaythomas](https://github.com/yaythomas) for their work on the documentation site.


Hugo now has:

* 47025+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 438+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 354+ [themes](http://themes.gohugo.io/)

## Notes


We have added a `force` flag to the [server redirects](https://gohugo.io/getting-started/configuration/#configure-server) configuration, configuring whether to override any existing content in the path or not. This is inline with how [Netlify](https://docs.netlify.com/routing/redirects/#syntax-for-the-netlify-configuration-file) does it.

This is set to default `false`. If you want the old behaviour you need to add this flag to your configuration:

{{< code-toggle file="config" >}}
[[redirects]]
from = "/myspa/**"
to = "/myspa/"
status = 200
force = true
{{< /code-toggle >}}

## Enhancements

### Templates

* Add Do Not Track (dnt) option to Vimeo shortcode [edc5c474](https://github.com/gohugoio/hugo/commit/edc5c4741caaee36ba4d42b5947c195a3e02e6aa) [@joshgerdes](https://github.com/joshgerdes) [#7700](https://github.com/gohugoio/hugo/issues/7700)

### Other

* Regen docshelper [b9318e43](https://github.com/gohugoio/hugo/commit/b9318e4315d9112f727140c0950d8836bf26eb87) [@bep](https://github.com/bep) 
* Make BuildConfig.Render an enum [63493890](https://github.com/gohugoio/hugo/commit/634938908ec8f393b9a05d26b4cfe19ca7abb0d0) [@bep](https://github.com/bep) [#7783](https://github.com/gohugoio/hugo/issues/7783)
* Allow cascade to be a slice with a _target discriminator [c63db7f1](https://github.com/gohugoio/hugo/commit/c63db7f1f6774a2d661af1d8197c6fe377e3ad25) [@bep](https://github.com/bep) [#7782](https://github.com/gohugoio/hugo/issues/7782)
* Add force flag to server redirects config [5e2a547c](https://github.com/gohugoio/hugo/commit/5e2a547cb594b31ecb0f089b08db2e15c6dc381a) [@bep](https://github.com/bep) [#7778](https://github.com/gohugoio/hugo/issues/7778)
* bump github.com/evanw/esbuild from 0.7.8 to 0.7.9 [ee090c09](https://github.com/gohugoio/hugo/commit/ee090c0940cdbf636e3a55a40b41612d92b9c62d) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/tdewolff/minify/v2 from 2.9.5 to 2.9.7 [05e358fd](https://github.com/gohugoio/hugo/commit/05e358fd335bcb5c7bdc2783ab0c17ec42667df6) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.34.34 to 1.35.0 [a2e85d9a](https://github.com/gohugoio/hugo/commit/a2e85d9a75aca59fd720cce6561ff64997858cd2) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/getkin/kin-openapi from 0.22.0 to 0.22.1 [4fba78dd](https://github.com/gohugoio/hugo/commit/4fba78dd0e950742132954a5d24629e4adfa1bb1) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.34.33 to 1.34.34 [c011b466](https://github.com/gohugoio/hugo/commit/c011b4667f3e1e3c6ecea2fe8f251578884c53b6) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.7.7 to 0.7.8 [35348b4b](https://github.com/gohugoio/hugo/commit/35348b4b343600ec24b1eb1a06f4d3c59199df25) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.34.27 to 1.34.33 [34915777](https://github.com/gohugoio/hugo/commit/34915777c2e8bc1457ff90d09cf814d494d9eece) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.7.4 to 0.7.7 [0f4a837e](https://github.com/gohugoio/hugo/commit/0f4a837ed1fd903bb6740b512683528ddb917918) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/tdewolff/minify/v2 from 2.9.4 to 2.9.5 [b395d686](https://github.com/gohugoio/hugo/commit/b395d686e9a77bf4e0d587ee9a3af4ae6e1aee02) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Upgrade to go-i18n v2 [97987e5c](https://github.com/gohugoio/hugo/commit/97987e5c0254e35668dca7f89e67b79553e617c8) [@bep](https://github.com/bep) [#5242](https://github.com/gohugoio/hugo/issues/5242)
* bump github.com/evanw/esbuild from 0.7.2 to 0.7.4 [4855c186](https://github.com/gohugoio/hugo/commit/4855c186d8f05e5e1b0f681b4aa6482a033df241) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.34.26 to 1.34.27 [6f07ec7e](https://github.com/gohugoio/hugo/commit/6f07ec7e9ec5c43f78100aa36b82786ba0260d75) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/alecthomas/chroma from 0.8.0 to 0.8.1 [4318dc72](https://github.com/gohugoio/hugo/commit/4318dc72f8c562b3bc106cd953d9fce58a93455d) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.7.1 to 0.7.2 [acdc27a3](https://github.com/gohugoio/hugo/commit/acdc27a32de83f32557e7a108797ddbebe4eb464) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Make sure CSS is rebuilt when postcss.config.js or tailwind.config.js changes [3acde9ae](https://github.com/gohugoio/hugo/commit/3acde9ae04fbf4a8c635d404608cb87218a8b803) [@bep](https://github.com/bep) [#7715](https://github.com/gohugoio/hugo/issues/7715)
* bump github.com/aws/aws-sdk-go from 1.34.22 to 1.34.26 [0bce9770](https://github.com/gohugoio/hugo/commit/0bce97703c17318b13b95d78ba41f40efb06aea7) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update to  github.com/tdewolff/minify v2.9.4 [b254532b](https://github.com/gohugoio/hugo/commit/b254532b52785954c98a473a635b9cea016d8565) [@bep](https://github.com/bep) 
* Bump bundled Node.js from v12.18.3 to v12.18.4 [05a22892](https://github.com/gohugoio/hugo/commit/05a22892921bd4618efe6135ce0d6fe2be545607) [@anthonyfok](https://github.com/anthonyfok) 
* Add preserveTOC option [8e553dcd](https://github.com/gohugoio/hugo/commit/8e553dcdefe50ab534f1199c006ae7754e14bee5) [@helfper](https://github.com/helfper) 
* bump github.com/frankban/quicktest from 1.10.2 to 1.11.0 [d4fc70a3](https://github.com/gohugoio/hugo/commit/d4fc70a3b320a55c4f571eed806d5ad5fdf1ef14) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.6.32 to 0.7.1 [d905abc0](https://github.com/gohugoio/hugo/commit/d905abc002aa6fd260e82063ef1edb8876aa76fd) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/rogpeppe/go-internal from 1.5.1 to 1.6.2 [8f394674](https://github.com/gohugoio/hugo/commit/8f3946746dda444f183ba235288c2b39d0d6a943) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/jdkato/prose from 1.1.1 to 1.2.0 [b01b2564](https://github.com/gohugoio/hugo/commit/b01b2564eefe342c9bf9767ffc256ebd04b94c71) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/spf13/afero from 1.2.2 to 1.4.0 [9fa5ebe2](https://github.com/gohugoio/hugo/commit/9fa5ebe2c42fbb37d066ffcd36bad4d08efe879a) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Preserve the original package.json if it exists [214afe4c](https://github.com/gohugoio/hugo/commit/214afe4c1bb9c37bc6159e659d66ba9a268a2849) [@bep](https://github.com/bep) [#7690](https://github.com/gohugoio/hugo/issues/7690)

## Fixes

### Templates

* Fix grammar in the new 'requires non-zero' error message [cd830bb0](https://github.com/gohugoio/hugo/commit/cd830bb0275fc39240861627ef26e146985b5c86) [@nekr0z](https://github.com/nekr0z) 

### Other

* Fix writeStats with quote inside quotes [11134411](https://github.com/gohugoio/hugo/commit/111344113bf8c16ae45528d67ff408da15961727) [@bep](https://github.com/bep) [#7746](https://github.com/gohugoio/hugo/issues/7746)
* Fix CLI example for PostCSS 8 [0c3d2b67](https://github.com/gohugoio/hugo/commit/0c3d2b67e0af38a4c3935fb04f722a73ec1d3f8b) [@ai](https://github.com/ai) 
* Fix typo in redirect error message [473b6610](https://github.com/gohugoio/hugo/commit/473b6610d51d4a33ba35917f95b0d97ea78dad2b) [@jmooring](https://github.com/jmooring) 
* Fix nilpointer for images with no Exif [cd00f7f9](https://github.com/gohugoio/hugo/commit/cd00f7f9661d67951ef16c5198541f09f1c058b4) [@bep](https://github.com/bep) [#7688](https://github.com/gohugoio/hugo/issues/7688)





