
---
date: 2019-12-23
title: "Hugo Christmas Edition!"
description: "Hugo 0.62 brings Markdown Render Hooks. And it's faster!"
categories: ["Releases"]
---

From all of us to all of you, a **very Merry Christmas** -- and Hugo `0.62.0`! This version brings [Markdown Render Hooks](https://gohugo.io/getting-started/configuration-markup/#markdown-render-hooks). This gives you full control over how links and images in Markdown are rendered without using any shortcodes. With this, you can get Markdown links that work on both GitHub and Hugo, resize images etc. It is a very long sought after feature, that has been hard to tackle until we got [Goldmark](https://github.com/yuin/goldmark/), the new Markdown engine, by [@yuin](https://github.com/yuin). When you read up on this new feature in the documentation, also note the new [.RenderString](https://gohugo.io/functions/renderstring/) method on `Page`.

Adding these render hooks also had the nice side effect of making Hugo **faster and more memory effective**. We could have just added this feature on top of what we got, getting it to work. But you like Hugo's fast builds, you love instant browser-refreshes on change. So we had to take a step back and redesign how we detect "what changed?" for templates referenced from content files, either directly or indirectly. And by doing that we greatly simplified how we handle all the templates. Which accidentally makes this version  **the fastest to date**. It's not an "every site will be much faster" statement. This depends. Sites with many languages and/or many templates will benefit more from this. We have benchmarks with site-building showing about 15% improvement in build speed and memory efficiency.

This release represents **25 contributions by 5 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@gavinhoward](https://github.com/gavinhoward), [@niklasfasching](https://github.com/niklasfasching), and [@zaitseff](https://github.com/zaitseff) for their ongoing contributions. And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), which has received **8 contributions by 5 contributors**. A special thanks to [@bep](https://github.com/bep), [@DirtyF](https://github.com/DirtyF), [@pfhawkins](https://github.com/pfhawkins), and [@bubelov](https://github.com/bubelov) for their work on the documentation site.

Also a big shoutout and thanks to the very active and helpful moderators on the [Hugo Discourse](https://discourse.gohugo.io/), making it a first class forum for Hugo questions and discussions.

Hugo now has:

* 40362+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 440+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 284+ [themes](http://themes.gohugo.io/)

## Notes

* Ace and Amber support is now removed from Hugo. See [#6609](https://github.com/gohugoio/hugo/issues/6609) for more information.
* The `markdownify` template function does not, yet, support render hooks. We recommend you look at the new and more powerful [.RenderString](https://gohugo.io/functions/renderstring/) method on `Page`.
* If you have output format specific behaviour in a template used from a content file, you must create a output format specific template, e.g. `myshortcode.amp.html`. This also applies to the new rendering hooks introduced in this release. This has been the intended behaviour all the time, but a failing test (now fixed) shows that the implementation of this has not been as strict as specified, hence this note.
* The `errorf` does not return any value anymore. This means that the ERROR will just be printed to the console. We have also added a `warnf` template func.


## Enhancements

### Templates

* Do not return any value in errorf [50cc7fe5](https://github.com/gohugoio/hugo/commit/50cc7fe54580018239ea95aafe67f6a158cdcc9f) [@bep](https://github.com/bep) [#6653](https://github.com/gohugoio/hugo/issues/6653)
* Add a warnf template func [1773d71d](https://github.com/gohugoio/hugo/commit/1773d71d5b40f5a6a14edca417d2818607a499f1) [@bep](https://github.com/bep) [#6628](https://github.com/gohugoio/hugo/issues/6628)
* Some more params merge adjustments [ccb1bf1a](https://github.com/gohugoio/hugo/commit/ccb1bf1abb7341fa1be23a90b66c14ae89790f49) [@bep](https://github.com/bep) [#6633](https://github.com/gohugoio/hugo/issues/6633)
* Get rid of the custom template truth logic [d20ca370](https://github.com/gohugoio/hugo/commit/d20ca3700512d661247b44d953515b9455e57ed6) [@bep](https://github.com/bep) [#6615](https://github.com/gohugoio/hugo/issues/6615)
* Add some comments [92c7f7ab](https://github.com/gohugoio/hugo/commit/92c7f7ab85a40cae8f36f2348d86f3e47d811eb5) [@bep](https://github.com/bep) 

### Core

* Improve error and reload handling  of hook templates in server mode [8a58ebb3](https://github.com/gohugoio/hugo/commit/8a58ebb311fd079f65068e7e37725e4d43f17ab5) [@bep](https://github.com/bep) [#6635](https://github.com/gohugoio/hugo/issues/6635)

### Other

* Update Goldmark to v1.1.18 [1fb17be9](https://github.com/gohugoio/hugo/commit/1fb17be9a008b549d11b622849adbaad01d4023d) [@bep](https://github.com/bep) [#6649](https://github.com/gohugoio/hugo/issues/6649)
* Update go-org [51d89dab](https://github.com/gohugoio/hugo/commit/51d89dab1827ae80f9d865f5c38cb5f6a3a11f68) [@niklasfasching](https://github.com/niklasfasching) 
* More on hooks [c8bfe47c](https://github.com/gohugoio/hugo/commit/c8bfe47c6a740c5fedfdb5b7465d7ae1db44cb65) [@bep](https://github.com/bep) 
* Update to Goldmark v1.1.17 [04536838](https://github.com/gohugoio/hugo/commit/0453683816cfbc94e1e19c644f5f84213bb8cf35) [@bep](https://github.com/bep) [#6641](https://github.com/gohugoio/hugo/issues/6641)
* Regen docshelper [55c29d4d](https://github.com/gohugoio/hugo/commit/55c29d4de38df67dd116f1845f7cc69ca7e35843) [@bep](https://github.com/bep) 
* Preserve HTML Text for image render hooks [a67d95fe](https://github.com/gohugoio/hugo/commit/a67d95fe1a033ca4934957b5a98b12ecc8a9edbd) [@bep](https://github.com/bep) [#6639](https://github.com/gohugoio/hugo/issues/6639)
* Update Goldmark [eef934ae](https://github.com/gohugoio/hugo/commit/eef934ae7eabc38eeba386831de6013eec0285f2) [@bep](https://github.com/bep) [#6626](https://github.com/gohugoio/hugo/issues/6626)
* Preserve HTML Text for link render hooks [00954c5d](https://github.com/gohugoio/hugo/commit/00954c5d1fda0b18cd1b847ee580d5f4caa76c70) [@bep](https://github.com/bep) [#6629](https://github.com/gohugoio/hugo/issues/6629)
* Footnote [3e316155](https://github.com/gohugoio/hugo/commit/3e316155c5d4fbf166d38e997a41101b6aa501d5) [@bep](https://github.com/bep) 
* Add render template hooks for links and images [e625088e](https://github.com/gohugoio/hugo/commit/e625088ef5a970388ad50e464e87db56b358dac4) [@bep](https://github.com/bep) [#6545](https://github.com/gohugoio/hugo/issues/6545)[#4663](https://github.com/gohugoio/hugo/issues/4663)[#6043](https://github.com/gohugoio/hugo/issues/6043)
* Enhance accessibility to issues [0947cf95](https://github.com/gohugoio/hugo/commit/0947cf958358e5a45b4f605e2a5b2504896fa360) [@peaceiris](https://github.com/peaceiris) [#6233](https://github.com/gohugoio/hugo/issues/6233)
* Re-introduce the correct version of Goldmark [03d6960a](https://github.com/gohugoio/hugo/commit/03d6960a15dcc8efc164e5ed310b12bd1ffdd930) [@bep](https://github.com/bep) 
* Rework template handling for function and map lookups [a03c631c](https://github.com/gohugoio/hugo/commit/a03c631c420a03f9d90699abdf9be7e4fca0ff61) [@bep](https://github.com/bep) [#6594](https://github.com/gohugoio/hugo/issues/6594)
* Create lightweight forks of text/template and html/template [167c0153](https://github.com/gohugoio/hugo/commit/167c01530bb295c8b8d35921eb27ffa5bee76dfe) [@bep](https://github.com/bep) [#6594](https://github.com/gohugoio/hugo/issues/6594)
* Add config option for ordered list [4c804319](https://github.com/gohugoio/hugo/commit/4c804319f6db0b8459cc9b5df4a904fd2c55dedd) [@gavinhoward](https://github.com/gavinhoward) 

## Fixes

### Templates

* Fix merge vs Params [1b785a7a](https://github.com/gohugoio/hugo/commit/1b785a7a6d3c264e39e4976c59b618c0ac1ba5f9) [@bep](https://github.com/bep) [#6633](https://github.com/gohugoio/hugo/issues/6633)

### Core

* Fix test [3c24ae03](https://github.com/gohugoio/hugo/commit/3c24ae030fe08ba259dd3de7ffea6c927c01e070) [@bep](https://github.com/bep) 

### Other

* Fix abs path handling in module mounts [ad6504e6](https://github.com/gohugoio/hugo/commit/ad6504e6b504277bbc7b60d093cdccd4f3baaa4f) [@bep](https://github.com/bep) [#6622](https://github.com/gohugoio/hugo/issues/6622)
* Fix incorrect MIME type from image/jpg to image/jpeg [158e7ec2](https://github.com/gohugoio/hugo/commit/158e7ec204e5149d77893d353cac9f55946d3e9a) [@zaitseff](https://github.com/zaitseff) 





