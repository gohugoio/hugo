
---
date: 2018-09-24
title: "Hugo 0.49: Directory Based Archetypes"
description: "Hugo 0.49 brings archetype bundle support and collection goodness."
categories: ["Releases"]
---

	
Hugo `0.49` brings [directory based archetypes](https://gohugo.io/content-management/archetypes/#directory-based-archetypes) and also improves the language handling in `hugo new`. This should simplify working with [page bundles](https://gohugo.io/content-management/page-bundles/). One example of this would be how you now can create a [new showcase](https://gohugo.io/showcase/template/) for the Hugo web site.

But this release is also about _collections_. Go 1.11 gave us variable overwrite support in Go templates. That made it possible to simplify a lot of template constructs, but it also showed some limitations in Hugo's template functions. So with this release we have:

* added [append](https://gohugo.io/functions/append/) function to append to collections.
* added [group](https://gohugo.io/functions/group/) to create custom page groups.
* improved the type support in [slice](https://gohugo.io/functions/slice/).

This release represents **66 contributions by 9 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@felicianotech](https://github.com/felicianotech), and [@vdanjean](https://github.com/vdanjean) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **37 contributions by 20 contributors**. A special thanks to [@bep](https://github.com/bep), [@kaushalmodi](https://github.com/kaushalmodi), [@AlexChambers](https://github.com/AlexChambers), and [@shaform](https://github.com/shaform) for their work on the documentation site.


Hugo now has:

* 28985+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 265+ [themes](http://themes.gohugo.io/)


## Notes

* Remove deprecated `rssURI` [f1a00b20](https://github.com/gohugoio/hugo/commit/f1a00b2069ede85feb487d29b9f690396e2402c6) [@bep](https://github.com/bep) 
* Remove deprecated flags [df4cbbd3](https://github.com/gohugoio/hugo/commit/df4cbbd3bdc05aa14a67b3a0a29a0db75b82e640) [@bep](https://github.com/bep) 
* Deprecate `Pages.Sort`. Use `.ByWeight` [2e2e34a9](https://github.com/gohugoio/hugo/commit/2e2e34a9350edec0220462aa3d47ecc9d428a0fb) [@bep](https://github.com/bep)
* When setting `preserveTaxonomyNames` Hugo now _really_ preserves them. Before this release, we would make the first character upper case. If this is the behaviour you want you can use the new `strings.FirstUpper` function.

## Enhancements

### Templates

* Allow `first` function to return an empty slice [cae07ce8](https://github.com/gohugoio/hugo/commit/cae07ce84b3bd4a33fd18b5109a1a3c3dce2191c) [@felicianotech](https://github.com/felicianotech) [#5235](https://github.com/gohugoio/hugo/issues/5235)
* Use `safeHTMLAttr` instead of `safeHTML` for HTML attributes [4f9c109d](https://github.com/gohugoio/hugo/commit/4f9c109dc5431553e5dbf98e0ed37487c12e8d16) [@felicianotech](https://github.com/felicianotech) [#5236](https://github.com/gohugoio/hugo/issues/5236)
* Add `collections.Append` (with alias `append`) [e27fd4c1](https://github.com/gohugoio/hugo/commit/e27fd4c1b80b7acb43290ac50e9f140d690cf042) [@bep](https://github.com/bep) [#5190](https://github.com/gohugoio/hugo/issues/5190)
* Improve type handling in `collections.Slice` [fe6676c7](https://github.com/gohugoio/hugo/commit/fe6676c775b8d917a661238f24fd4a9088f25d50) [@bep](https://github.com/bep) [#5188](https://github.com/gohugoio/hugo/issues/5188)
* Add `group` template func [6667c6d7](https://github.com/gohugoio/hugo/commit/6667c6d7430acc16b3683fbbacd263f1d00c8672) [@bep](https://github.com/bep) [#4865](https://github.com/gohugoio/hugo/issues/4865)
* Add `strings.FirstUpper` [e5d66074](https://github.com/gohugoio/hugo/commit/e5d66074ce1ed4e0fe329e3fdef66f8b6fd5dc55) [@bep](https://github.com/bep) [#5174](https://github.com/gohugoio/hugo/issues/5174)

### Core

* Minor cleaning in the sorting code [2eed35c8](https://github.com/gohugoio/hugo/commit/2eed35c826e5de6aae432b36969a28c2ae3e0f02) [@bep](https://github.com/bep) 
* Make sure ambiguous lookups in GetPage gets an error [75e54345](https://github.com/gohugoio/hugo/commit/75e54345f9a7d786bb28af64ad80eb9502fee7c7) [@bep](https://github.com/bep) [#5138](https://github.com/gohugoio/hugo/issues/5138)
* Allow creating page groups from any page collection [cfda13b3](https://github.com/gohugoio/hugo/commit/cfda13b36367465016f4458ab9924c948ed02b6f) [@vdanjean](https://github.com/vdanjean) [#4865](https://github.com/gohugoio/hugo/issues/4865)
* Do not FirstUpper taxonomy titles [be3ae3ec](https://github.com/gohugoio/hugo/commit/be3ae3ec92da972a55112af39ce2e1c45121b9a5) [@Schnouki](https://github.com/Schnouki) [#5172](https://github.com/gohugoio/hugo/issues/5172)
* Simplify some code [e38e8812](https://github.com/gohugoio/hugo/commit/e38e881248b7d20927eab0e56c85732e1acbc45e) [@moorereason](https://github.com/moorereason) 
* Add missing error checks [0665a395](https://github.com/gohugoio/hugo/commit/0665a3951be6ffc2551ef6664856da4cfccd45fa) [@moorereason](https://github.com/moorereason) 
* Remove extraneous createStaticFs call [1ed8c363](https://github.com/gohugoio/hugo/commit/1ed8c363367c2532014154e91eeade9b3a923f91) [@moorereason](https://github.com/moorereason) 

### Other

* Add "go mod download" to CI scripts [cf47f43f](https://github.com/gohugoio/hugo/commit/cf47f43ff12ca5f5ea851a8b6761b19b5e4d4eba) [@bep](https://github.com/bep) 
* Regenerate CLI docs [3b6bd121](https://github.com/gohugoio/hugo/commit/3b6bd1210a2792c51c34b9c655cb8b7e9a0f15d7) [@bep](https://github.com/bep) 
* Document group [8388cd90](https://github.com/gohugoio/hugo/commit/8388cd90e89358f73ddcb7f496a1a8dc5c30c36c) [@bep](https://github.com/bep) 
* Make Data.Integrity be of type template.HTMLAttr [fe6a6f27](https://github.com/gohugoio/hugo/commit/fe6a6f2737769070fd64a5192ff685c9c89020bd) [@bep](https://github.com/bep) 
* Add directory based archetypes [2650fa77](https://github.com/gohugoio/hugo/commit/2650fa772b40846d9965f8c5f169286411f3beb2) [@bep](https://github.com/bep) [#4535](https://github.com/gohugoio/hugo/issues/4535)
* Build on CircleCI outside of GOPATH [ef525b15](https://github.com/gohugoio/hugo/commit/ef525b15d4584886b52428bd7a35de835ab07a48) [@felicianotech](https://github.com/felicianotech) [#5135](https://github.com/gohugoio/hugo/issues/5135)
* Prevent symbolic links from themes [f9168146](https://github.com/gohugoio/hugo/commit/f9168146978bd970d1f4fb061eff75264af88bb1) [@bep](https://github.com/bep) 
* Update releasenotes_writer.go [4b82f748](https://github.com/gohugoio/hugo/commit/4b82f74848836efbcf453c0122bd35555ee7517d) [@bep](https://github.com/bep) 
* Add docs for append [df50c108](https://github.com/gohugoio/hugo/commit/df50c108ba2f24936eff20b51d23f9328adb2d87) [@bep](https://github.com/bep) [#5190](https://github.com/gohugoio/hugo/issues/5190)
* Set minifier to KeepEndTags [9b26b548](https://github.com/gohugoio/hugo/commit/9b26b5487b5c5142fe9fb58681fe7d1dac95a291) [@onedrawingperday](https://github.com/onedrawingperday) 
* Make JSON minification more generic [3dafe206](https://github.com/gohugoio/hugo/commit/3dafe206e31bb92f27802a04bf9159cbc20af234) [@zinefer](https://github.com/zinefer) 
* Update Mage [37d64634](https://github.com/gohugoio/hugo/commit/37d6463479952f7dfba59d899eed38b41e223283) [@bep](https://github.com/bep) 
* Update dependencies [bb2fe814](https://github.com/gohugoio/hugo/commit/bb2fe814c2db0c494b3b678a5da20a6cc0538857) [@bep](https://github.com/bep) 
* Improve some godoc comments [30bc4ed0](https://github.com/gohugoio/hugo/commit/30bc4ed0a01f965cc2f9187ccb6ab5d28a3149f6) [@moorereason](https://github.com/moorereason) 
* Update to latest Mage [3b103cb7](https://github.com/gohugoio/hugo/commit/3b103cb7b74228f26af5beb4cefc47edee794ce9) [@bep](https://github.com/bep) 
* Remove some duplicate code [c15c7da4](https://github.com/gohugoio/hugo/commit/c15c7da42a1c7bc535cc16cca2b341526f8cf169) [@bep](https://github.com/bep) 
* Update Dockerfile to Go 1.11 [bcbe57c6](https://github.com/gohugoio/hugo/commit/bcbe57c6e9243cbf3823f11b755f57c091cc1866) [@zyfdegh](https://github.com/zyfdegh) [#5145](https://github.com/gohugoio/hugo/issues/5145)
* Init packages once [ea8ef573](https://github.com/gohugoio/hugo/commit/ea8ef573c6f869de95fdf4b19765d34026de6471) [@bep](https://github.com/bep) 
* Update script to Go 1.11 [293e1235](https://github.com/gohugoio/hugo/commit/293e12355dd9d9361774f5ab340cd8a03b4828a1) [@bep](https://github.com/bep) [#5127](https://github.com/gohugoio/hugo/issues/5127)
* Remove the remains of Go Dep [fdf3c3b8](https://github.com/gohugoio/hugo/commit/fdf3c3b8234ed340f40a85fb76d96ae3a9ccf195) [@bep](https://github.com/bep) [#5115](https://github.com/gohugoio/hugo/issues/5115)
* Update CONTRIBUTING.md [312d2252](https://github.com/gohugoio/hugo/commit/312d2252be6b7bf250fa4f8b1b541fdc13641940) [@bep](https://github.com/bep) 
* Update README.md [f627903e](https://github.com/gohugoio/hugo/commit/f627903efaa1a5f7e137c2d409efd1e1e2db47f6) [@bep](https://github.com/bep) 
* Fix golint issues [400fe96a](https://github.com/gohugoio/hugo/commit/400fe96aee8e38112e347e762661b8389701c938) [@moorereason](https://github.com/moorereason) 
* Fix golint godoc issues [3f45e729](https://github.com/gohugoio/hugo/commit/3f45e729f4e0296bb1a3558d60087bec8321444b) [@moorereason](https://github.com/moorereason) 
* Fix godoc comment [e03eb90a](https://github.com/gohugoio/hugo/commit/e03eb90a366159ed9ef9888246de87f283508866) [@moorereason](https://github.com/moorereason) 
* Fix typo in private func name [c915d0d3](https://github.com/gohugoio/hugo/commit/c915d0d3252007d61b680a388dcbe6b035d0adc8) [@moorereason](https://github.com/moorereason)
* Fix golint godoc issues [f6f22ad9](https://github.com/gohugoio/hugo/commit/f6f22ad944a1c308fd823792b2fbff1504f42cef) [@moorereason](https://github.com/moorereason) 
* Fix filepath issue in test [d970327d](https://github.com/gohugoio/hugo/commit/d970327d7b994b495ef3bb468c3e0599b0deef5a) [@moorereason](https://github.com/moorereason) 
* Fix golint issues [0013bea9](https://github.com/gohugoio/hugo/commit/0013bea901ee2124f4c18f9728abf47c3880f97d) [@moorereason](https://github.com/moorereason) 
* Fix golint godoc issue [ffaa73dc](https://github.com/gohugoio/hugo/commit/ffaa73dc8aa860edb3476b2a460774071b8470a4) [@moorereason](https://github.com/moorereason) 
* Fix golint godoc issue [f8d8c854](https://github.com/gohugoio/hugo/commit/f8d8c85428f527139c20369910230741dcaf2969) [@moorereason](https://github.com/moorereason) 
* Fix golint issue [10dc87bf](https://github.com/gohugoio/hugo/commit/10dc87bf866f7a4f99c248436c38edf0ecdd157f) [@moorereason](https://github.com/moorereason) 
* Fix golint issues [7231869b](https://github.com/gohugoio/hugo/commit/7231869ba87f4e8d08e94dce18f20b7ed4fa2e15) [@moorereason](https://github.com/moorereason) 
* Fix golint godoc issues [600047ff](https://github.com/gohugoio/hugo/commit/600047ff1cb95d061af1983b9a755157eb4941f8) [@moorereason](https://github.com/moorereason) 
* Fix golint issues [5f2e1cb8](https://github.com/gohugoio/hugo/commit/5f2e1cb8969c2adac6c866b57cc331e1bc16d4e9) [@moorereason](https://github.com/moorereason) 
* Fix golint issues [c8ce6504](https://github.com/gohugoio/hugo/commit/c8ce65046dc7539f3bf5f6dd35fa7ece2bec866d) [@moorereason](https://github.com/moorereason) 
* Fix most golint issues [a53f9623](https://github.com/gohugoio/hugo/commit/a53f962312e273cea9fe460b40655350a82210f2) [@moorereason](https://github.com/moorereason) 
* Fix golint issues [daffeec3](https://github.com/gohugoio/hugo/commit/daffeec30d9d67017ec84064e15fd946b0b0cb0d) [@moorereason](https://github.com/moorereason) 
* Fix golint errors [b8b91f55](https://github.com/gohugoio/hugo/commit/b8b91f550646b2620649c3504e14a441975bea9f) [@moorereason](https://github.com/moorereason) 
* Fix golint issues [f0effac8](https://github.com/gohugoio/hugo/commit/f0effac80426325040c4bc703cd610f434d0b5a8) [@moorereason](https://github.com/moorereason) 
* Fix broken Travis config [2cf8fe2e](https://github.com/gohugoio/hugo/commit/2cf8fe2ea218d37776af72893691e772737750e3) [@bep](https://github.com/bep) 
* Fix error message for go vet [47d4edce](https://github.com/gohugoio/hugo/commit/47d4edce6083bab1c190dad99fefb7c73afc6af8) [@mdhender](https://github.com/mdhender) 


## Fixes

* Compare every element in pages cache [ed4f1edb](https://github.com/gohugoio/hugo/commit/ed4f1edbd729bf75af89879b76fbad931693cd67) [@bep](https://github.com/bep) [#5239](https://github.com/gohugoio/hugo/issues/5239)
* Revise error handling in `getJSON` and `getCSV` [43d44652](https://github.com/gohugoio/hugo/commit/43d446522a7c09af4bf6879f93341d8ff62654d1) [@bep](https://github.com/bep) [#5076](https://github.com/gohugoio/hugo/issues/5076)
* Show error on `union` or `intersect` of uncomparable types [4f72e791](https://github.com/gohugoio/hugo/commit/4f72e79120a4f964330d10c8ebe9aceb2b5761a7) [@moorereason](https://github.com/moorereason) [#3820](https://github.com/gohugoio/hugo/issues/3820)
* Do not set RSS as Kind in RSS output [555a5612](https://github.com/gohugoio/hugo/commit/555a5612b2641075b3e1b3b7af8ce9b5aba9f200) [@bep](https://github.com/bep) [#5138](https://github.com/gohugoio/hugo/issues/5138)








