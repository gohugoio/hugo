
---
date: 2018-01-31
title: "Hugo 0.35: Headless Bundles!"
description: "Headless Bundles, disable languages, improves fast render mode, and much more."
categories: ["Releases"]
---

The most notable new feature in Hugo `0.35` is perhaps **Headless Bundles**.

This means that you in your `index.md` front matter can say:

```yaml
headless: true
```
And

* it will have no `Permalink` and no rendered HTML in `/public`
* it will not be part of `.Site.RegularPages` etc.

But you can get it by:

* `.Site.GetPage ...`

The use cases are many:

* Shared media libraries
* Reusable page content "snippets"
* ...

But this release contains more goodies than possible to sum up in one paragraph, so study the release notes carefully. It represents **42 contributions by 8 contributors** to the main Hugo code base.

[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@vassudanagunta](https://github.com/vassudanagunta), [@yanzay](https://github.com/yanzay), and [@robertbasic](https://github.com/robertbasic) for their ongoing contributions.

And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **28 contributions by 5 contributors**. A special thanks to [@bep](https://github.com/bep), [@kaushalmodi](https://github.com/kaushalmodi), [@regisphilibert](https://github.com/regisphilibert), and [@salim-b](https://github.com/salim-b) for their work on the documentation site.


Hugo now has:

* 22967+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 448+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 197+ [themes](http://themes.gohugo.io/)


## Notes

* Deprecate `useModTimeAsFallback` [adfd4370](https://github.com/gohugoio/hugo/commit/adfd4370b67fd7181178bd6b3b1d07356beaac71) [@bep](https://github.com/bep) [#4351](https://github.com/gohugoio/hugo/issues/4351)
* Deprecate CLI flags `canonifyURLs`, `pluralizeListTitles`, `preserveTaxonomyNames`, `uglyURLs` [f08ea02d](https://github.com/gohugoio/hugo/commit/f08ea02d24d42929676756950f3affaca7fd8c01) [@bep](https://github.com/bep) [#4347](https://github.com/gohugoio/hugo/issues/4347)
* Remove undraft command [2fa70c93](https://github.com/gohugoio/hugo/commit/2fa70c9344b231c9d999bbafdfa4acbf27ed9f6e) [@robertbasic](https://github.com/robertbasic) [#4353](https://github.com/gohugoio/hugo/issues/4353)

## Enhancements

### Templates

* Update Twitter card to also consider images in `.Resources` [25d691da](https://github.com/gohugoio/hugo/commit/25d691daff57d7c6d7d0f63af3991d22e3f788fe) [@bep](https://github.com/bep) [#4349](https://github.com/gohugoio/hugo/issues/4349)
* Seed random on init only [83c761b7](https://github.com/gohugoio/hugo/commit/83c761b71a980aee6331179b271c7e24e999e8eb) [@liguoqinjim](https://github.com/liguoqinjim) 
* Remove duplicate layout lookup layouts [b2fcbb1f](https://github.com/gohugoio/hugo/commit/b2fcbb1f9774aa1e929b8575c0e1ac366ab2fb73) [@bep](https://github.com/bep) [#4319](https://github.com/gohugoio/hugo/issues/4319)
* Add "removable-media" interface to `snapcraft.yaml` [f0c0ece4](https://github.com/gohugoio/hugo/commit/f0c0ece44d55b6c2997cbd106d1bc099ea1a2fa7) [@anthonyfok](https://github.com/anthonyfok) [#3837](https://github.com/gohugoio/hugo/issues/3837)

### Other

* Add a way to disable one or more languages [6413559f](https://github.com/gohugoio/hugo/commit/6413559f7575e2653d76227a8037a7edbaae82aa) [@bep](https://github.com/bep) [#4297](https://github.com/gohugoio/hugo/issues/4297)[#4329](https://github.com/gohugoio/hugo/issues/4329)
* Handle newly created files in Fast Render Mode [1707dae8](https://github.com/gohugoio/hugo/commit/1707dae8d3634006017eb6d040df4dbafc53d92f) [@yanzay](https://github.com/yanzay) [#4339](https://github.com/gohugoio/hugo/issues/4339)
* Extract the Fast Render Mode logic into a method [94e736c5](https://github.com/gohugoio/hugo/commit/94e736c5e167a0ee70a528e1c19d64a47e7929c2) [@bep](https://github.com/bep) [#4339](https://github.com/gohugoio/hugo/issues/4339)
* Remove unused code [ae5a45be](https://github.com/gohugoio/hugo/commit/ae5a45be6f0ee4d5c52b38fd28b22b55d9cd7b2d) [@bep](https://github.com/bep) 
* Add the last lookup variant for the `GetPage` index [3446fe9b](https://github.com/gohugoio/hugo/commit/3446fe9b8937610b8b628b2c212eb25888a7c1bb) [@bep](https://github.com/bep) [#4312](https://github.com/gohugoio/hugo/issues/4312)
* Simplify bundle lookup via `.Site.GetPage`, `ref`, `relref` [517b6b62](https://github.com/gohugoio/hugo/commit/517b6b62389d23bfe41fe3ae551a691b11bdcaa7) [@bep](https://github.com/bep) [#4312](https://github.com/gohugoio/hugo/issues/4312)
* Remove some now superflous Fast Render Mode code [feeed073](https://github.com/gohugoio/hugo/commit/feeed073c3320b09fb38168ce272ac88b987f1d2) [@bep](https://github.com/bep) [#4339](https://github.com/gohugoio/hugo/issues/4339)
* Make resource counters for `name` and `title` independent [df20b054](https://github.com/gohugoio/hugo/commit/df20b05463fef42aba93d5208e410a7ecc56da5d) [@bep](https://github.com/bep) [#4335](https://github.com/gohugoio/hugo/issues/4335)
* Provide .Name to the archetype templates [863a812e](https://github.com/gohugoio/hugo/commit/863a812e07193541b42732b0e227f3d320433f01) [@bep](https://github.com/bep) [#4348](https://github.com/gohugoio/hugo/issues/4348)
* Only set `url` if permalink in metadata and remove duplicate confirm msg [3752348e](https://github.com/gohugoio/hugo/commit/3752348ef13ced8f6f528b42ee7d76a12a97ae5c) [@lildude](https://github.com/lildude) [#1887](https://github.com/gohugoio/hugo/issues/1887)
* Start Resources :counter first time they're used [7b472e46](https://github.com/gohugoio/hugo/commit/7b472e46084b603045b87cea870ffc73ac1cf7e7) [@bep](https://github.com/bep) [#4335](https://github.com/gohugoio/hugo/issues/4335)
* Update to Go 1.9.3 [a91aba1c](https://github.com/gohugoio/hugo/commit/a91aba1c1562259dffd321a608f38c38dd4d5aeb) [@bep](https://github.com/bep) [#4328](https://github.com/gohugoio/hugo/issues/4328)
* Support pages without front matter [91bb774a](https://github.com/gohugoio/hugo/commit/91bb774ae4e129f7ed0624754b31479c960ef774) [@vassudanagunta](https://github.com/vassudanagunta) [#4320](https://github.com/gohugoio/hugo/issues/4320)
* Add page metadata dates tests [3f0379ad](https://github.com/gohugoio/hugo/commit/3f0379adb72389954ca2be6a9f2ebfcd65c6c440) [@vassudanagunta](https://github.com/vassudanagunta) 
* Re-generate CLI docs [1e27d058](https://github.com/gohugoio/hugo/commit/1e27d0589118a114e49c032e4bd68b4798e44a5b) [@bep](https://github.com/bep) 
* Remove and update deprecation status [d418c2c2](https://github.com/gohugoio/hugo/commit/d418c2c2eacdc1dc6fffe839e0a90600867878ca) [@bep](https://github.com/bep) 
* Shorten the stale setup [4a7c2b36](https://github.com/gohugoio/hugo/commit/4a7c2b3695fe7b88861f2155ea7ef635fe425cd4) [@bep](https://github.com/bep) 
* Add a `GetPage` to the site benchmarks [a1956391](https://github.com/gohugoio/hugo/commit/a19563910eec5fed08f3b02563b9a7b38026183d) [@bep](https://github.com/bep) 
* Add headless bundle support [0432c64d](https://github.com/gohugoio/hugo/commit/0432c64dd22e4610302162678bb93661ba68d758) [@bep](https://github.com/bep) [#4311](https://github.com/gohugoio/hugo/issues/4311)
* Merge matching resources params maps [5a0819b9](https://github.com/gohugoio/hugo/commit/5a0819b9b5eb9e79826cfa0a65f235d9821b1ac4) [@bep](https://github.com/bep) [#4315](https://github.com/gohugoio/hugo/issues/4315)
* Add some general code contribution criteria [78c86330](https://github.com/gohugoio/hugo/commit/78c863305f337ed4faf3cf0a23675f28b0ae5641) [@bep](https://github.com/bep) 
* Tighten page kind logic, introduce tests [8125b4b0](https://github.com/gohugoio/hugo/commit/8125b4b03d10eb73f8aea3f9ea41172aba8df082) [@vassudanagunta](https://github.com/vassudanagunta) 

## Fixes
* Fix `robots.txt` in multihost mode [4d912e2a](https://github.com/gohugoio/hugo/commit/4d912e2aad39bfe8d76672cf53b01317792e02c5) [@bep](https://github.com/bep) [#4193](https://github.com/gohugoio/hugo/issues/4193)
* Fix `--uglyURLs` from comand line regression [016398ff](https://github.com/gohugoio/hugo/commit/016398ffe2e0a073453cf46a9d6bf72d693c11e5) [@bep](https://github.com/bep) [#4343](https://github.com/gohugoio/hugo/issues/4343)
* Avoid unescape in `highlight` [ebdd8cba](https://github.com/gohugoio/hugo/commit/ebdd8cba3f5965a8ac897833f313d772271de649) [@bep](https://github.com/bep) [#4219](https://github.com/gohugoio/hugo/issues/4219)
* Fix Docker build [a34213f0](https://github.com/gohugoio/hugo/commit/a34213f0b5624de101272aab469ca9b6fe0c273f) [@skoblenick](https://github.com/skoblenick) [#4076](https://github.com/gohugoio/hugo/issues/4076)[#4077](https://github.com/gohugoio/hugo/issues/4077)
* Fix language params handling [ae742cb1](https://github.com/gohugoio/hugo/commit/ae742cb1bdf35b81aa0ede5453da6b0c4a4fccf2) [@bep](https://github.com/bep) [#4356](https://github.com/gohugoio/hugo/issues/4356)[#4352](https://github.com/gohugoio/hugo/issues/4352)
* Fix handling of top-level page bundles [4eb2fec6](https://github.com/gohugoio/hugo/commit/4eb2fec67c3a72a3ac98aa834dc56fd4504626d8) [@bep](https://github.com/bep) [#4332](https://github.com/gohugoio/hugo/issues/4332)
* Fix `baseURL` server regression for multilingual sites [ed4a00e4](https://github.com/gohugoio/hugo/commit/ed4a00e46f2344320a22f07febe5aec4075cb3fb) [@bep](https://github.com/bep) [#4333](https://github.com/gohugoio/hugo/issues/4333)
* Fix "date" page param [322c5672](https://github.com/gohugoio/hugo/commit/322c567220aa4123a5d707629c1bebd375599912) [@vassudanagunta](https://github.com/vassudanagunta) [#4323](https://github.com/gohugoio/hugo/issues/4323)
* Fix typo in comment [912147ab](https://github.com/gohugoio/hugo/commit/912147ab896e69a450b7100c3d6bf81a7bf78b5a) [@yanzay](https://github.com/yanzay) 





