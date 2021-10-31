
---
date: 2021-09-02
title: "Go 1.17 Update"
description: "Hugo 0.88.0 brings Go 1.17 support, a dependency refresh, and more."
categories: ["Releases"]
---

The most important piece in this release is the Go 1.17 update. This release is built with that new Go version, but also and more importantly, the Hugo Modules logic has been updated to support Go's new way of (lazy) loading transitive dependencies. If you already have Go 1.17 installed, building existing Hugo Modules backed projects have not been an issue, but `hugo mod init` for a new project could give you _too new_ versions of transitive dependencies. Hugo 0.88 fixes this.

This release represents **26 contributions by 6 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@dependabot[bot]](https://github.com/apps/dependabot), [@helfper](https://github.com/helfper), and [@wzshiming](https://github.com/wzshiming) for their ongoing contributions.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **9 contributions by 6 contributors**. A special thanks to [@bep](https://github.com/bep), [@jmooring](https://github.com/jmooring), [@StevenMaude](https://github.com/StevenMaude), and [@coliff](https://github.com/coliff) for their work on the documentation site.

Hugo now has:

* 53915+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 430+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 404+ [themes](http://themes.gohugo.io/)

## Notes
* We have fixed a bug with the import order in `js.Build` for the case `./foo` when both `./foo.js` and `./foo/index.js` exists. This is in line with how both Node and ESBuild's native import resolver does it. We discovered this trying to build AlpineJS v3 from source mounted in `/assets`. See [cf73cc2e](https://github.com/gohugoio/hugo/commit/cf73cc2ececd4e794df09ea382a38ab18960d84e) [@bep](https://github.com/bep) [#8945](https://github.com/gohugoio/hugo/issues/8945).

## Enhancements

### Templates

* Handle nil values in time.AsTime [abd969a6](https://github.com/gohugoio/hugo/commit/abd969a670852f9ed57c1a26434445aa985706fe) [@bep](https://github.com/bep) 
* Handle nil values in time.AsTime [3e110728](https://github.com/gohugoio/hugo/commit/3e11072892ca31bb76980ee38890a4bd92d83dfd) [@bep](https://github.com/bep) [#8865](https://github.com/gohugoio/hugo/issues/8865)

### Other

* Run go mod tidy [6631c9c7](https://github.com/gohugoio/hugo/commit/6631c9c7e00fb9dc237b4ec2fbb261d05df268d1) [@bep](https://github.com/bep) 
* Don't fail on template errors on go mod graph etc. [7d1f806e](https://github.com/gohugoio/hugo/commit/7d1f806ecb3621ae7b545a686d04de4568814055) [@bep](https://github.com/bep) [#8942](https://github.com/gohugoio/hugo/issues/8942)
* bump github.com/getkin/kin-openapi from 0.74.0 to 0.75.0 [04b59599](https://github.com/gohugoio/hugo/commit/04b59599613a62d378bf3710ac0eb06c9543b96d) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/frankban/quicktest from 1.13.0 to 1.13.1 [c278b6e4](https://github.com/gohugoio/hugo/commit/c278b6e45d56b101db9691347f9e5a99a9319572) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.12.22 to 0.12.24 [107c86fe](https://github.com/gohugoio/hugo/commit/107c86febbb7057c4ae90c6a35b3e8eda24297c7) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Avoid failing with "module not found" for hugo mod init and similar [a0489c2d](https://github.com/gohugoio/hugo/commit/a0489c2dfd3ceb4d0702de0da7a4af3eabce05e5) [@bep](https://github.com/bep) [#8940](https://github.com/gohugoio/hugo/issues/8940)
* Update to Go 1.17 [0fc2ce9e](https://github.com/gohugoio/hugo/commit/0fc2ce9e4bf0524994a861b7300e4332f6f8d390) [@bep](https://github.com/bep) [#8930](https://github.com/gohugoio/hugo/issues/8930)
* Remove Pygments from snapcraft.yml [32569285](https://github.com/gohugoio/hugo/commit/32569285c181c8798ef594c12d3cfd7f9a252a04) [@anthonyfok](https://github.com/anthonyfok) 
* bump github.com/fsnotify/fsnotify from 1.4.9 to 1.5.0 [5a46eefb](https://github.com/gohugoio/hugo/commit/5a46eefbc6da3463b796ada8d15902be197455a3) [@bep](https://github.com/bep) [#8920](https://github.com/gohugoio/hugo/issues/8920)
* Add tabindex when code is not highlighted [7a15edaf](https://github.com/gohugoio/hugo/commit/7a15edafe240471c072d3548b72ccda0271ffd8f) [@helfper](https://github.com/helfper) 
* bump github.com/evanw/esbuild from 0.12.17 to 0.12.22 [2f0945ba](https://github.com/gohugoio/hugo/commit/2f0945bafe501103abe97b2f2b5566b28ec48e52) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump golang.org/x/text from 0.3.6 to 0.3.7 [7ba3f3d2](https://github.com/gohugoio/hugo/commit/7ba3f3d201e386cb9c7c15df5a6cc1c4b46473bd) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/fsnotify/fsnotify from 1.4.9 to 1.5.0 [f7016524](https://github.com/gohugoio/hugo/commit/f70165242b98e3ee182fbac08bf2893a7f09e961) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Prevent minifier from removing quotes around post-processed attributes [bc0743ed](https://github.com/gohugoio/hugo/commit/bc0743ed8eafc3c2d9b21a1e8f1b05d64b85e8ba) [@bep](https://github.com/bep) [#8884](https://github.com/gohugoio/hugo/issues/8884)
* Avoid too many watch file handles causing the server to fail to start [ffa2fe61](https://github.com/gohugoio/hugo/commit/ffa2fe61172aa0d892234b23d1497c77a6a7f5c4) [@bep](https://github.com/bep) 
* Remove some pygments references [d966f5d0](https://github.com/gohugoio/hugo/commit/d966f5d08d7f75f1ae9acd94e292bf61de2adf0d) [@helfper](https://github.com/helfper) 
* Avoid too many watch file handles causing the server to fail to start [3f38c785](https://github.com/gohugoio/hugo/commit/3f38c785b7208440e2a9dd9a80cb39d4ae23e676) [@wzshiming](https://github.com/wzshiming) [#8904](https://github.com/gohugoio/hugo/issues/8904)
* bump github.com/getkin/kin-openapi from 0.68.0 to 0.74.0 [24589c08](https://github.com/gohugoio/hugo/commit/24589c0814bc5d21565470bec6215ee792f1655e) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update github.com/spf13/cast v1.4.0 => v1.4.1 [efebd756](https://github.com/gohugoio/hugo/commit/efebd756eb1f35c515ac82ccc85ec520bac91240) [@bep](https://github.com/bep) [#8891](https://github.com/gohugoio/hugo/issues/8891)
* Import time/tzdata on Windows [58b6742c](https://github.com/gohugoio/hugo/commit/58b6742cfeb6d4cd04450cbe9592209510c2b977) [@bep](https://github.com/bep) [#8892](https://github.com/gohugoio/hugo/issues/8892)
* Indent TOML tables [9bba9a3a](https://github.com/gohugoio/hugo/commit/9bba9a3a98fa268391597d8d7a52112fb401d952) [@bep](https://github.com/bep) [#8850](https://github.com/gohugoio/hugo/issues/8850)

## Fixes

### Other

* Fix import order for ./foo when both ./foo.js and ./foo/index.js exists [cf73cc2e](https://github.com/gohugoio/hugo/commit/cf73cc2ececd4e794df09ea382a38ab18960d84e) [@bep](https://github.com/bep) [#8945](https://github.com/gohugoio/hugo/issues/8945)
* Fix it so disableKinds etc. does not get merged in from theme [f4ffeea7](https://github.com/gohugoio/hugo/commit/f4ffeea71dd3d044a2628bbb5d6634680667398f) [@bep](https://github.com/bep) [#8866](https://github.com/gohugoio/hugo/issues/8866)
* Fix `lang.FormatPercent` description [d6c8cd77](https://github.com/gohugoio/hugo/commit/d6c8cd771834ae2913658c652e30a9feadc2a7b7) [@salim-b](https://github.com/salim-b) 
