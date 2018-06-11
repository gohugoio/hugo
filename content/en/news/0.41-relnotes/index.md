
---
date: 2018-05-25
title: "Hugo 0.41: Privacy Configuration for GDPR"
description: "Hugo 0.41 adds a new site configuration section to meet the new General Data Protection Regulation (GDPR)."
categories: ["Releases"]
---

	
In Hugo `0.41` we add a new **Privacy Configuration** section to meet the new regulations in the new **General Data Protection Regulation ([GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation))**. Many have contributed work on this project, but a special thanks to [@onedrawingperday](https://github.com/onedrawingperday), [@jhabdas](https://github.com/jhabdas), and [@it-gro](https://github.com/it-gro).

> You can read more about the new Privacy Config [here](https://gohugo.io/about/hugo-and-gdpr/).


Hugo now has:

* 25917+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 446+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 231+ [themes](http://themes.gohugo.io/)

## Notes

* We have fixed an issue where we sent the wrong `.Site` object into the archetype templates used in `hugo new`. This meant that if you wanted to access site params and some other fields you needed to use `.Site.Info.Params` etc. This release fixes that so now you can use the same construct in the archetype templates as in the regular templates.

## Enhancements

* Alias tweet shortode to twitter [3bfe8f4b](https://github.com/gohugoio/hugo/commit/3bfe8f4be653f44674293685cb5750d90668b2f6) [@bep](https://github.com/bep) [#4765](https://github.com/gohugoio/hugo/issues/4765)
* Adjust GA templates [f45b522e](https://github.com/gohugoio/hugo/commit/f45b522ebffafc61a3cb9b694bc3542747c73e07) [@bep](https://github.com/bep) [#4616](https://github.com/gohugoio/hugo/issues/4616)
* Wrap the relevant templates with the privacy policy disable check [67892073](https://github.com/gohugoio/hugo/commit/6789207347fc2df186741644a6fe968d41ea9077) [@bep](https://github.com/bep) [#4616](https://github.com/gohugoio/hugo/issues/4616)
* Extract internal templates [34ad9a4f](https://github.com/gohugoio/hugo/commit/34ad9a4f178fcf50abe7246ad9d30b294327da16) [@bep](https://github.com/bep) [#4457](https://github.com/gohugoio/hugo/issues/4457)
* Use double quotes instead of back quotes [b2b500f5](https://github.com/gohugoio/hugo/commit/b2b500f563c3bb36751a4c1610df113c4daad604) [@anthonyfok](https://github.com/anthonyfok) 
* Provide the correct .Site object to archetype templates [ab02594e](https://github.com/gohugoio/hugo/commit/ab02594e09c0414124186e42d67d52d474dd341a) [@bep](https://github.com/bep) [#4732](https://github.com/gohugoio/hugo/issues/4732)
* Document the GDPR Privacy Config [c71f201f](https://github.com/gohugoio/hugo/commit/c71f201fd93287afa7cb7b875bd523c25e48400c) [@bep](https://github.com/bep) [#4751](https://github.com/gohugoio/hugo/issues/4751)
* Add no-cookie variants of the Google Analytics templates [a51945ea] (https://github.com/gohugoio/hugo/commit/a51945ea4b99d17501d73cf3367926683e4a4dfd) [@bep](https://github.com/bep) [#4775](https://github.com/gohugoio/hugo/issues/4775)
* Remove youtube_simple for now [448081b8](https://github.com/gohugoio/hugo/commit/448081b840db4a23c0c49c2d869ac207dcb6ac40) [@bep](https://github.com/bep) [#4751](https://github.com/gohugoio/hugo/issues/4751)
* Add anonymizeIP to GA privacy config [1f1d955b](https://github.com/gohugoio/hugo/commit/1f1d955b56471e41d5288c57f1ef8333dc297120) [@bep](https://github.com/bep) [#4751](https://github.com/gohugoio/hugo/issues/4751)
* Support DNT in Twitter shortcode for GDPR [9753cb59](https://github.com/gohugoio/hugo/commit/9753cb59f1f1d866943a485dd7c917d1b68f6eda) [@bep](https://github.com/bep) [#4765](https://github.com/gohugoio/hugo/issues/4765)
* Regenerate embedded templates [6aa2c385](https://github.com/gohugoio/hugo/commit/6aa2c38507aa1c2246222684717b4d69d26b03d7) [@bep](https://github.com/bep) [#4761](https://github.com/gohugoio/hugo/issues/4761)
* Add instagram_simple shortcode [9ad46a20](https://github.com/gohugoio/hugo/commit/9ad46a20357a7e28b405feef5c8f7d4501186da6) [@bep](https://github.com/bep) [#4748](https://github.com/gohugoio/hugo/issues/4748)
* Go fmt [4256de33](https://github.com/gohugoio/hugo/commit/4256de3392d320a5a47fcab49882f2a3249c2163) [@bep](https://github.com/bep) 
* Remove the id from youtube_simple [bed7a0fa](https://github.com/gohugoio/hugo/commit/bed7a0faff90bbe389629347026853b7bc4c8c3f) [@bep](https://github.com/bep) [#4751](https://github.com/gohugoio/hugo/issues/4751)
* Add an unified .Site.Config with a services section [4ddcf52c](https://github.com/gohugoio/hugo/commit/4ddcf52ccc7af3e23109ebaac1f0486087a212ba) [@bep](https://github.com/bep) [#4751](https://github.com/gohugoio/hugo/issues/4751)
* Move the privacy config into a parent [353148c2](https://github.com/gohugoio/hugo/commit/353148c2bc2cdb9f2eb8ee967ba756ce09323801) [@bep](https://github.com/bep) [#4751](https://github.com/gohugoio/hugo/issues/4751)
* Make the simple mode YouTube links schemaless [69ee6b41](https://github.com/gohugoio/hugo/commit/69ee6b41e36625595e2bcabcde0bc58663e5b93c) [@bep](https://github.com/bep) [#4616](https://github.com/gohugoio/hugo/issues/4616)
* Add YouTube shortcode simple mode [88e35686](https://github.com/gohugoio/hugo/commit/88e356868062cc618385cd22b6730df2459518cd) [@bep](https://github.com/bep) [#4616](https://github.com/gohugoio/hugo/issues/4616)
* Do not return error on .Get "class" and vice versa in shortcodes [2f17f937](https://github.com/gohugoio/hugo/commit/2f17f9378ad96c4a9f6d7d24b0776ed3a25a08a3) [@bep](https://github.com/bep) [#4745](https://github.com/gohugoio/hugo/issues/4745)
* Create SUPPORT.md [0a7027e2](https://github.com/gohugoio/hugo/commit/0a7027e2a87283743d5310b74e18666e4a64d3e1) [@coliff](https://github.com/coliff) 
* Add PrivacyEnhanced mode for YouTube to the GDPR Policy [5f24a2c0](https://github.com/gohugoio/hugo/commit/5f24a2c047db0bff8c9e267bfa8ef8e43e6bd24e) [@bep](https://github.com/bep) [#4616](https://github.com/gohugoio/hugo/issues/4616)
* Add RespectDoNotTrack to GDPR privacy policy for Google Analytics [71014201](https://github.com/gohugoio/hugo/commit/710142016b140538bfc11e48bb32d26fa685b2ad) [@bep](https://github.com/bep) [#4616](https://github.com/gohugoio/hugo/issues/4616)
* Add the foundation for GDPR privacy configuration [0bbdef98](https://github.com/gohugoio/hugo/commit/0bbdef986d8eecf4fabe9a372e33626dbdfeb36b) [@bep](https://github.com/bep) [#4616](https://github.com/gohugoio/hugo/issues/4616)
* Show site build warning in TestPageBundlerSiteRegular [9bd4236e](https://github.com/gohugoio/hugo/commit/9bd4236e1b3bee332439eef50e12d4481340c3eb) [@anthonyfok](https://github.com/anthonyfok) [#4672](https://github.com/gohugoio/hugo/issues/4672)
* Do not show empty BuildDate in version [4eedb377](https://github.com/gohugoio/hugo/commit/4eedb377b60fb6742c97398942a0045ff2a824c4) [@anthonyfok](https://github.com/anthonyfok) 
* Improve markup determination logic [2fb9af59](https://github.com/gohugoio/hugo/commit/2fb9af59c14b1732ba1a2f21794e2cf8dfca0604) [@vassudanagunta](https://github.com/vassudanagunta) 
* Update CONTRIBUTING.md [b6ededf0](https://github.com/gohugoio/hugo/commit/b6ededf0591a81667754f1dccef2c6fe6d342811) [@domdomegg](https://github.com/domdomegg) 

## Fixes

* Avoid ANSI character output on Windows [568b4335](https://github.com/gohugoio/hugo/commit/568b4335c20effb46168bd639317a3420f563463) [@LorenzCK](https://github.com/LorenzCK) [#4462](https://github.com/gohugoio/hugo/issues/4462)






