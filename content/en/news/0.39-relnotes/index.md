
---
date: 2018-04-16
title: "Hugo 0.39: The Nat King Cole Stabilizer Edition"
description: "Hugo 0.39: Rewrite of the /commands package, `Resource.Content`, several new template funcs, and more …"
categories: ["Releases"]
---

	
**Nat King Cole** was a fantastic American jazz pianist. When his bass player had visited the bar a little too often, he started with his percussive	piano playing to keep the tempo flowing. Oscar Peterson called it _"Nat's stabilizers"_. This release is the software equivalent of that. We have been doing frequent main releases this year, but looking back, the patch releases that followed them seemed unneeded. And looking at the regressions, most of them stem from the `commands` package, a package that before this release was filled with globals and high coupling. This package is now rewritten and accompanied with decent test coverage.

But this release isn't all boring and technical: It includes several important bug fixes, several useful new template functions, and `Resource.Content` allows you to get any resource's content without having to fiddle with file paths and `readFile`.

This release represents **61 contributions by 4 contributors** to the main Hugo code base. A shout-out to [@bep](https://github.com/bep) for the implementation and [@it-gro](https://github.com/it-gro) and [@RickCogley](https://github.com/RickCogley) for the help testing it.

[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions and his witty Norwegian humour, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@thedodobird2](https://github.com/thedodobird2), and [@neurocline](https://github.com/neurocline) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **6 contributions by 5 contributors**. A special thanks to [@kaushalmodi](https://github.com/kaushalmodi), [@regisphilibert](https://github.com/regisphilibert), [@bep](https://github.com/bep), and [@tomanistor](https://github.com/tomanistor) for their work on the documentation site.

Hugo now has:

* 24911+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 446+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 218+ [themes](http://themes.gohugo.io/)

## Notes

* The `main.Execute` function now returns a `Response` object and the global `Hugo` variable is removed. This is only relevant for people building some kind of API around Hugo.
* Remove deprecated `File.Bytes` [94c8b29c](https://github.com/gohugoio/hugo/commit/94c8b29c39d0c485ee91d98c08fd615c28802496) [@bep](https://github.com/bep) 

## Enhancements

### Templates

* Add `anchorize` template func [4dba6ce1](https://github.com/gohugoio/hugo/commit/4dba6ce15ae9b5208b1e2d68c96d7b1dce0a07ab) [@bep](https://github.com/bep) 
* Add `path.Join` [880ca19f](https://github.com/gohugoio/hugo/commit/880ca19f209e68e6a8daa6686b361515ecacc91e) [@bep](https://github.com/bep) 
* Add `path.Split` template func [01b72eb5](https://github.com/gohugoio/hugo/commit/01b72eb592d0e0aefc5f7ae42f9f6ff112883bb6) [@bep](https://github.com/bep) 

### Other

* Implement `Resource.Content` [0e7716a4](https://github.com/gohugoio/hugo/commit/0e7716a42450401c7998aa81ad2ed98c8ab109e8) [@bep](https://github.com/bep) [#4622](https://github.com/gohugoio/hugo/issues/4622)
* Make `Page.Content` a method that returns interface{} [417c5e2b](https://github.com/gohugoio/hugo/commit/417c5e2b67b97fa80a0b6f77d259966f03b95344) [@bep](https://github.com/bep) [#4622](https://github.com/gohugoio/hugo/issues/4622)
* Remove accidental and breaking space in baseURL flag [1b4e0c41](https://github.com/gohugoio/hugo/commit/1b4e0c4161fb631add62e77f494a7e62c3619020) [@bep](https://github.com/bep) [#4607](https://github.com/gohugoio/hugo/issues/4607)
* Properly handle CLI slice arguments [27a524b0](https://github.com/gohugoio/hugo/commit/27a524b0905ec73c1eef233f94700feb9f465011) [@bep](https://github.com/bep) [#4607](https://github.com/gohugoio/hugo/issues/4607)
* Correctly handle destination and i18n-warnings [bede93de](https://github.com/gohugoio/hugo/commit/bede93de005dcf934f3ec9be6388310ac6c57acd) [@bep](https://github.com/bep) [#4607](https://github.com/gohugoio/hugo/issues/4607)
* Allow "*/" inside commented out shortcodes [14c35c8a](https://github.com/gohugoio/hugo/commit/14c35c8a56c4dc9a1ee0053e9ff976be7715ba99) [@bep](https://github.com/bep) [#4608](https://github.com/gohugoio/hugo/issues/4608)
* Make commands.Execute return a Response object [96689a5c](https://github.com/gohugoio/hugo/commit/96689a5c319f720368491226f034d0ff9585217c) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Remove some TODOs [e7010c1b](https://github.com/gohugoio/hugo/commit/e7010c1b621d68ee53411a5ba8143d07b976d9fe) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Add basic server test [a7d00fc3](https://github.com/gohugoio/hugo/commit/a7d00fc39e87a5cac99b3a2380f5cc8c135d2b4b) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Remove the Hugo global [b110d0ae](https://github.com/gohugoio/hugo/commit/b110d0ae04e13fb45c739bcebb580709745082e6) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the limit command work again [73825cfc](https://github.com/gohugoio/hugo/commit/73825cfc1c0b007830b24bb1947a565175b52d36) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Move the commands related logic to its own file [a8f7fbbb](https://github.com/gohugoio/hugo/commit/a8f7fbbb10aa78f3ebac008d29d9969bb197393c) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Add CLI tests [e8d6ca95](https://github.com/gohugoio/hugo/commit/e8d6ca9531d19e4e898c57d77d2fd627ea38ade0) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the hugo command non-global [4d32f2fa](https://github.com/gohugoio/hugo/commit/4d32f2fa8969f368b088dc9bcedb45f2c986cb27) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Extract some common types into its own file [018602c4](https://github.com/gohugoio/hugo/commit/018602c46db8d729af2871bd5f4c1e7480420f09) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the server command non-global [2f0d98a1](https://github.com/gohugoio/hugo/commit/2f0d98a19b021d03930003217b0519afaef3a391) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the gen commands non-global [e0621d20](https://github.com/gohugoio/hugo/commit/e0621d207ce3278a82f8a60607e9cdd304149029) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the list commands non-global [e26a8b24](https://github.com/gohugoio/hugo/commit/e26a8b242a6434117d089a0799238add7025dbf4) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the import commands non-global [2a2c9838](https://github.com/gohugoio/hugo/commit/2a2c9838671b5401331d20f8c72e2b934fe34e8d) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the config command non-global [15b1e269](https://github.com/gohugoio/hugo/commit/15b1e269ade91ddc6a74c552bc61b0c5e527d268) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make the new commands non-global [56a13080](https://github.com/gohugoio/hugo/commit/56a13080446283ed1cde6b69fc6f4fac85076c84) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make convert command non-global [4b780ca7](https://github.com/gohugoio/hugo/commit/4b780ca778ee7f25af808da38ede964a01698c70) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make more commands non-global [7bc5e89f](https://github.com/gohugoio/hugo/commit/7bc5e89fbaa5c613b8853ff7b69fae570bd0b56d) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Make benchmark non-global [fdf1d94e](https://github.com/gohugoio/hugo/commit/fdf1d94ebc7d1aa4855c62237f2edbd4bdade1a7) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Start of flag cleaning [1157fef8](https://github.com/gohugoio/hugo/commit/1157fef85908ea54883fe0dba6adc4861ba02162) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Use short date format in CLI docs [e614d8a5](https://github.com/gohugoio/hugo/commit/e614d8a57c2ff5eef9270d51fcc6518398d7ff88) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Update README.md [fca49d6c](https://github.com/gohugoio/hugo/commit/fca49d6c608d227049cb2f26895cfecc685f1c89) [@thedodobird2](https://github.com/thedodobird2) 
* Sync dependencies [0e8b3cbc](https://github.com/gohugoio/hugo/commit/0e8b3cbcd274e1f2e14be694c794a544f49efb56) [@bep](https://github.com/bep) 
* Bump Go versions [230f2b8c](https://github.com/gohugoio/hugo/commit/230f2b8c4fce03f14847de2b22402e64d4d69783) [@bep](https://github.com/bep) [#4545](https://github.com/gohugoio/hugo/issues/4545)
* Add bash completion [874159b5](https://github.com/gohugoio/hugo/commit/874159b5436bc9080aec71a9c26d35f8f62c9fd0) [@anthonyfok](https://github.com/anthonyfok) 
* Handle mass content etc. edits in server mode [730b66b6](https://github.com/gohugoio/hugo/commit/730b66b6520f263af16f555d1d7be51205a8e51d) [@bep](https://github.com/bep) [#4563](https://github.com/gohugoio/hugo/issues/4563)

## Fixes

* Fix `livereload` of bundled pages [f3775877](https://github.com/gohugoio/hugo/commit/f3775877c61c11ab7c8fd1fc3e15470bf5da4820) [@bep](https://github.com/bep) [#4607](https://github.com/gohugoio/hugo/issues/4607)
* Do not reset `.Page.Scratch` on rebuilds [61d52f14](https://github.com/gohugoio/hugo/commit/61d52f146297950e283ae086d8b1af61099d22a0) [@bep](https://github.com/bep) [#4627](https://github.com/gohugoio/hugo/issues/4627)
* Fix failing Travis server test [9c782d51](https://github.com/gohugoio/hugo/commit/9c782d5147bfea0dd85cf3374f598f0176f204eb) [@bep](https://github.com/bep) 
* Fix the config command [f396cffa](https://github.com/gohugoio/hugo/commit/f396cffa239e948075af2224208671956d8b4a84) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Fix some flag diff [24d5c219](https://github.com/gohugoio/hugo/commit/24d5c219424a9777bb1dd366b43e68e6f47e1adb) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Fix TestFixURL [1e233b1c](https://github.com/gohugoio/hugo/commit/1e233b1c4598fd8cbce7da8a67bf2c4918c6047e) [@bep](https://github.com/bep) [#4598](https://github.com/gohugoio/hugo/issues/4598)
* Disable shallow clone to fix TestPageWithLastmodFromGitInfo [094ec171](https://github.com/gohugoio/hugo/commit/094ec171420e659cdf962a19dd90105912ce9901) [@anthonyfok](https://github.com/anthonyfok) [#4584](https://github.com/gohugoio/hugo/issues/4584)
* Fix livereload for the home page bundle [f87239e4](https://github.com/gohugoio/hugo/commit/f87239e4cab958bf59ecfb1beb8cac439441a553) [@bep](https://github.com/bep) [#4576](https://github.com/gohugoio/hugo/issues/4576)
* Fix empty `BuildDate` in "hugo version" [294c0f80](https://github.com/gohugoio/hugo/commit/294c0f8001fe598278c1eb8015deb6b98e8de686) [@anthonyfok](https://github.com/anthonyfok) 
* Fix some livereload content regressions [a4deaeff](https://github.com/gohugoio/hugo/commit/a4deaeff0cfd70abfbefa6d40c0b86839a216f6d) [@bep](https://github.com/bep) [#4566](https://github.com/gohugoio/hugo/issues/4566)
* Fix two tests that are broken on Windows [26f34fd5](https://github.com/gohugoio/hugo/commit/26f34fd59da1ce1885d4f2909c5d9ef9c1726944) [@neurocline](https://github.com/neurocline) 
