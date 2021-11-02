
---
date: 2021-11-02
title: "Fine Grained File Filters"
description: "Hugo 0.89.0 brings fine grained file filters, archetype rewrite, dependency refresh, and more ..."
categories: ["Releases"]
---

This release is a dependency refresh (the new Goldmark version comes with a lot of bug fixes, as one example), many bug fixes, but also some nice new features:

We have added the [configuration settings](https://gohugo.io/hugo-modules/configuration/#module-config-mounts) **includeFiles** and **excludeFiles** to the mount configuration. This allows fine grained control over what files to include, and it works for all of Hugo's file systems (including `/static`).

We have also [reimplemented archetypes](https://github.com/gohugoio/hugo/pull/9045). The old implementation had some issues, mostly related to the context (e.g. name, file paths) passed to the template. This new implementation is using the exact same code path for evaluating the pages as in a regular build. This also makes it more robust and easier to reason about in a multilingual setup. Now, if you are explicit about the target path, Hugo will now always pick the correct mount and language:

```
hugo new content/en/posts/my-first-post.md
```

This release represents **50 contributions by 13 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@dependabot[bot]](https://github.com/apps/dependabot), [@jmooring](https://github.com/jmooring), and [@anthonyfok](https://github.com/anthonyfok) for their ongoing contributions.
And thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his ongoing work on keeping the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **23 contributions by 9 contributors**. A special thanks to [@jmooring](https://github.com/jmooring), [@bep](https://github.com/bep), [@coliff](https://github.com/coliff), and [@vipkr](https://github.com/vipkr) for their work on the documentation site.


Hugo now has:

* 54999+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 430+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 413+ [themes](http://themes.gohugo.io/)


## Notes

* Hugo now writes an empty file named `.hugo_build.lock` to the root of the project when building (also when doing `hugo new mypost.md` and other commands that requires a build). We recommend you just leave this file alone. Put it in `.gitignore` or similar if you don't want the file in your source repository.
* We have updated to ESBuild `v0.13.12`. The release notes for [v0.13.0](https://github.com/evanw/esbuild/releases/tag/v0.13.0) mentions a potential breaking change.
* We now only build AMD64 release binaries (see [this issue](https://github.com/gohugoio/hugo/issues/9102)) for the Unix OSes (e.g. NetBSD). If you need, say, a binary for ARM64, you need to build it yourself.
* We now build only one release binary/archive for MacOS (see [this issue](https://github.com/gohugoio/hugo/issues/9035)) that works on both Intel and the new Arm M1 systems.
* `.File.ContentBaseName` now returns the owning directory name for all bundles (branch an leaf). This is a bug fix, but worth mentioning. See [this issue](https://github.com/gohugoio/hugo/issues/9112).
* We have updated the Twitter shortcode to use Twitter's new API. See [this issue](https://github.com/gohugoio/hugo/pull/9106) for details.

## Enhancements

### Templates

* Use configured location when date passed to Format is string [e82cbd74](https://github.com/gohugoio/hugo/commit/e82cbd746fd4b07e40fedacc4247b9cd50ef70e7) [@bep](https://github.com/bep) [#9084](https://github.com/gohugoio/hugo/issues/9084)
* Add path.Clean [e55466ce](https://github.com/gohugoio/hugo/commit/e55466ce70363418309d465a0f2aa6c7ada1e51d) [@bradcypert](https://github.com/bradcypert) [#8885](https://github.com/gohugoio/hugo/issues/8885)

### Other

* Regen CLI docs [f503b639](https://github.com/gohugoio/hugo/commit/f503b6395707f8e576af734efab83092d62fae37) [@bep](https://github.com/bep) 
* Make ContentBaseName() return the directory for branch bundles [30aba7fb](https://github.com/gohugoio/hugo/commit/30aba7fb099678363b0a4828936ed28e740e00e2) [@bep](https://github.com/bep) [#9112](https://github.com/gohugoio/hugo/issues/9112)
* Update Twitter shortcode oEmbed endpoint [0cc39af6](https://github.com/gohugoio/hugo/commit/0cc39af68232f1a4981aae2e72cf65da762b5768) [@jmooring](https://github.com/jmooring) [#8130](https://github.com/gohugoio/hugo/issues/8130)
* bump github.com/evanw/esbuild from 0.13.10 to 0.13.12 [7fa66425](https://github.com/gohugoio/hugo/commit/7fa66425aa0a918b4bf5eb9a21f6e567e0a7e876) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/yuin/goldmark from 1.4.1 to 1.4.2 [69210cfd](https://github.com/gohugoio/hugo/commit/69210cfdf341d1faef23f4e9290d51448dd5e0c6) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.40.8 to 1.41.14 [3339c2bb](https://github.com/gohugoio/hugo/commit/3339c2bb618c29bb3ad442c71fe1542ad7195971) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/getkin/kin-openapi from 0.79.0 to 0.80.0 [03bbdba8](https://github.com/gohugoio/hugo/commit/03bbdba8be19929cb6a14243b690372fbfbc6aa6) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.13.8 to 0.13.10 [a772b8fc](https://github.com/gohugoio/hugo/commit/a772b8fc3833e010553c412dd5daa0175e6ccead) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Rename excepted filenames for image golden testdata [dce49d13](https://github.com/gohugoio/hugo/commit/dce49d13336f3dbadaa1359322a277ad4cb55679) [@anthonyfok](https://github.com/anthonyfok) [#6387](https://github.com/gohugoio/hugo/issues/6387)
* bump github.com/frankban/quicktest from 1.13.1 to 1.14.0 [61c5b7a2](https://github.com/gohugoio/hugo/commit/61c5b7a2e623255be99da7adf200f0591c9a1195) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Validate the target path in hugo new [75c9b893](https://github.com/gohugoio/hugo/commit/75c9b893d98961a504cff9ed3c89055d16e315d6) [@bep](https://github.com/bep) [#9072](https://github.com/gohugoio/hugo/issues/9072)
* Set zone of datetime from from `go-toml` [b959ecbc](https://github.com/gohugoio/hugo/commit/b959ecbc8175e2bf260f10b08965531bce9bcb7e) [@satotake](https://github.com/satotake) [#8895](https://github.com/gohugoio/hugo/issues/8895)
* Added nodesource apt repository to snap package [70e45481](https://github.com/gohugoio/hugo/commit/70e454812ef684d02ffa881becf0f8ce6a1b5f8c) [@sergiogarciadev](https://github.com/sergiogarciadev) 
* Set HUGO_ENABLEGITINFO=false override in Set_in_string [355ff83e](https://github.com/gohugoio/hugo/commit/355ff83e74f6e27c79033b8dfb899e3a3b529049) [@anthonyfok](https://github.com/anthonyfok) 
* Add includeFiles and excludeFiles to mount configuration [471ed91c](https://github.com/gohugoio/hugo/commit/471ed91c60cd36645794925cb4892cc820eae626) [@bep](https://github.com/bep) [#9042](https://github.com/gohugoio/hugo/issues/9042)
* bump github.com/mitchellh/mapstructure from 1.4.1 to 1.4.2 [94a5bac5](https://github.com/gohugoio/hugo/commit/94a5bac5b29bbba1ca4809752fe3fd04a58547b6) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Always preserve the original transform error [9830ca9e](https://github.com/gohugoio/hugo/commit/9830ca9e319f6ce313f4e542a202bd0d0469a9ed) [@bep](https://github.com/bep) 
* Add hyperlink to the banner [b64fd057](https://github.com/gohugoio/hugo/commit/b64fd0577b0fb222bea22ae347acb5dd17b2aa04) [@itsAftabAlam](https://github.com/itsAftabAlam) 
* bump github.com/getkin/kin-openapi from 0.78.0 to 0.79.0 [2706437a](https://github.com/gohugoio/hugo/commit/2706437a7d593b66b0fbad0235dbaf917593971b) [@dependabot[bot]](https://github.com/apps/dependabot) 
* github.com/evanw/esbuild v0.13.5 => v0.13.8 [ec7c993c](https://github.com/gohugoio/hugo/commit/ec7c993cfe216b8a3c6fbac85669cefef59778dd) [@bep](https://github.com/bep) 
* Return error on no content dirs [32c6f656](https://github.com/gohugoio/hugo/commit/32c6f656d93ecf4308f7c30848b13b4c6f157436) [@bep](https://github.com/bep) [#9056](https://github.com/gohugoio/hugo/issues/9056)
* Add a cross process build lock and use it in the archetype content builder [ba35e698](https://github.com/gohugoio/hugo/commit/ba35e69856900b6fc92681aa841cdcaefbb4b121) [@bep](https://github.com/bep) [#9048](https://github.com/gohugoio/hugo/issues/9048)
* github.com/alecthomas/chroma v0.9.2 => v0.9.4 [bb053770](https://github.com/gohugoio/hugo/commit/bb053770337e214f41bc1c524d458ba7fbe1fc08) [@bep](https://github.com/bep) [#8532](https://github.com/gohugoio/hugo/issues/8532)
* Reimplement archetypes [9185e11e](https://github.com/gohugoio/hugo/commit/9185e11effa682ea1ef7dc98f2943743671023a6) [@bep](https://github.com/bep) [#9032](https://github.com/gohugoio/hugo/issues/9032)[#7589](https://github.com/gohugoio/hugo/issues/7589)[#9043](https://github.com/gohugoio/hugo/issues/9043)[#9046](https://github.com/gohugoio/hugo/issues/9046)[#9047](https://github.com/gohugoio/hugo/issues/9047)
* bump github.com/tdewolff/minify/v2 from 2.9.21 to 2.9.22 [168a3aab](https://github.com/gohugoio/hugo/commit/168a3aab4622786ccd0943137fce3912707f2a46) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update github.com/evanw/esbuild v0.13.5 [8bcfa3bd](https://github.com/gohugoio/hugo/commit/8bcfa3bdf65492329da8093d841dd04c7a5a10c8) [@bep](https://github.com/bep) 
* bump github.com/mattn/go-isatty from 0.0.13 to 0.0.14 [cd4e67af](https://github.com/gohugoio/hugo/commit/cd4e67af182a1b3aa19db7609c7581c424e9310f) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/getkin/kin-openapi from 0.75.0 to 0.78.0 [e6ad1f0e](https://github.com/gohugoio/hugo/commit/e6ad1f0e763ee891bf4d71df0168b6949369c793) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Allow multiple plugins in the PostCSS options map [64abc83f](https://github.com/gohugoio/hugo/commit/64abc83fc4b70c70458c582ae2cf67fc9c67bb3f) [@jmooring](https://github.com/jmooring) [#9015](https://github.com/gohugoio/hugo/issues/9015)
* Create path.Clean documentation [f8d132d7](https://github.com/gohugoio/hugo/commit/f8d132d731cf8e27c8c17931597fd975e8a7c3cc) [@jmooring](https://github.com/jmooring) 
* Skip a test assertion on CI [26f1919a](https://github.com/gohugoio/hugo/commit/26f1919ae0cf57d754bb029270c20e76cc32cf4d) [@bep](https://github.com/bep) 
* Remove tracking image [ecf025f0](https://github.com/gohugoio/hugo/commit/ecf025f006f22061728e78f2cf50257dde2225ee) [@kambojshalabh35](https://github.com/kambojshalabh35) 
* Revert "Remove credit from release notes" [fab1e43d](https://github.com/gohugoio/hugo/commit/fab1e43de59f3a7596ab23347387d846139bc3a3) [@digitalcraftsman](https://github.com/digitalcraftsman) 
* Pass minification errors to the user [e03f82ee](https://github.com/gohugoio/hugo/commit/e03f82eef2679ec8963894d0b911363eef40941a) [@ptgott](https://github.com/ptgott) [#8954](https://github.com/gohugoio/hugo/issues/8954)
* Clarify "precision" in currency format functions [a864ffe9](https://github.com/gohugoio/hugo/commit/a864ffe9acf295034bb38e789a0efa62906b2ae4) [@ptgott](https://github.com/ptgott) 
* bump github.com/evanw/esbuild from 0.12.24 to 0.12.29 [b49da332](https://github.com/gohugoio/hugo/commit/b49da33280cb01795ce833e70c2b7b78cca1867e) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Use default math/rand.Source for concurrency safety [7c21eca7](https://github.com/gohugoio/hugo/commit/7c21eca74f95b61d6813d0c0b155bf07c9aa8575) [@odeke-em](https://github.com/odeke-em) [#8981](https://github.com/gohugoio/hugo/issues/8981)
* Make the error handling for the mod commands more lenient [13ad8408](https://github.com/gohugoio/hugo/commit/13ad8408fc6b645b12898fb8053388fc4848dfbd) [@bep](https://github.com/bep) 
* Add some help text to the 'unknown revision' error [1cabf61d](https://github.com/gohugoio/hugo/commit/1cabf61ddf96b89c95c3ba77a985168184920feb) [@bep](https://github.com/bep) [#6825](https://github.com/gohugoio/hugo/issues/6825)
* Update github.com/yuin/goldmark v1.4.0 => v1.4.1 [268e3069](https://github.com/gohugoio/hugo/commit/268e3069f37df01a5a58b615844652fb75b8503a) [@jmooring](https://github.com/jmooring) [#8855](https://github.com/gohugoio/hugo/issues/8855)

## Fixes

### Templates

* Fix time.Format with Go layouts [ed6fd26c](https://github.com/gohugoio/hugo/commit/ed6fd26ce884c49b02497728a99e90b92dd65f1f) [@bep](https://github.com/bep) [#9107](https://github.com/gohugoio/hugo/issues/9107)

### Other

* Fix description of lang.FormatNumberCustom [04a3b45d](https://github.com/gohugoio/hugo/commit/04a3b45db4cd28b4821b5c98cd67dfbf1d098957) [@jmooring](https://github.com/jmooring) 
* Fix typo in error message [1d60bd1e](https://github.com/gohugoio/hugo/commit/1d60bd1efa943349636edad3dd8c5427312ab0f1) [@jmooring](https://github.com/jmooring) 
* Fix panic when specifying multiple excludeFiles directives [64e1613f](https://github.com/gohugoio/hugo/commit/64e1613fb390bd893900dc0596e5c3f3c8e1cd8c) [@bep](https://github.com/bep) [#9076](https://github.com/gohugoio/hugo/issues/9076)
* Fix file permissions in new archetype implementation [e02e0727](https://github.com/gohugoio/hugo/commit/e02e0727e57f123f9a8de506e9c098bb374f7a23) [@bep](https://github.com/bep) [#9057](https://github.com/gohugoio/hugo/issues/9057)
* Fix the "page picker" logic in --navigateToChanged [096f5e19](https://github.com/gohugoio/hugo/commit/096f5e19217e985bccbf6c539e1b220541ffa6f6) [@bep](https://github.com/bep) [#9051](https://github.com/gohugoio/hugo/issues/9051)
* Fix a typo on OpenBSD [c7957c90](https://github.com/gohugoio/hugo/commit/c7957c90e83ff2b2cc958bd61486a244f0fd8891) [@nabbisen](https://github.com/nabbisen) 
* Fix value of useResourceCacheWhen in TestResourceChainPostCSS [e6e44b7c](https://github.com/gohugoio/hugo/commit/e6e44b7c41a9b517ffc3775ea0a6aec2b1d4591b) [@jmooring](https://github.com/jmooring) 
