
---
date: 2019-08-14
title: "Hugo 0.57: The Cascading Edition"
description: "Hugo 0.57 brings cascading front matter, alphabetical sorting, resource loading from Assets with wildcards. And it's faster."
categories: ["Releases"]
---

Hugo 0.57 brings **Cascading Front Matter**, **Alphabetical Sorting**, **Resources Loading from Assets with Wildcards**. And more.

**Cascading Front Matter**: We have added a new and powerful `cascade` keyword to Hugo's front matter. This can be added to any index node in `_index.md`. Any values in `cascade` will be merged into itself and all the descendants.

```yaml
title: "My Blog"
icon: "world.png"
cascade:
   icon: "flag.png"
   outputs: ["HTML"]
   type: "blog"
 ```
 
It's worth noting that the `cascade` element itself will also be merged. Also, to grasp the full value of this feature, remember that front matter in Hugo is both **data** and **behaviour**: You can tell Hugo how to process a subset of the pages (some example keywords are `layout`, `type`, `outputs`, `weight`) using the `cascade` keyword, e.g. _"I want this subsection to be rendered in both the `HTML` and `Calendar` Output Formats"_.
 
This feature is created by[@regisphilibert](https://github.com/regisphilibert) and [@bep](https://github.com/bep) See [#6041](https://github.com/gohugoio/hugo/issues/6041) for details.

**Resources Loading from Assets with Wildcards**: We have added two new sought after template functions to the `resources` namespace: `resources.Match` and `resources.GetMatch`. These behaves like their namesake methods on `Page` (with [super-asterisk wildcard support](https://github.com/gobwas/glob)), but searches in all the resources in Assets. E.g. `{{ $prettyImages := resources.Match "images/**pretty.jpg" }}` will give a slice of all "pretty pictures". Another relevant example: `{{ $js := resources.Match "libs/*.js" | resources.Concat "js/bundle.js" }}`.

**Performance:** In general, this version is slightly faster and more memory effective. In particular, we have fixed a performance issue with the replacer step that greatly improves the build speed of certain large and content-rich sites (thanks to [@vazrupe](https://github.com/vazrupe) for the fix).  

This release represents **46 contributions by 8 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@muesli](https://github.com/muesli), [@XhmikosR](https://github.com/XhmikosR), and [@vazrupe](https://github.com/vazrupe) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **13 contributions by 7 contributors**. A special thanks to [@regisphilibert](https://github.com/regisphilibert), [@bep](https://github.com/bep), [@kenberkeley](https://github.com/kenberkeley), and [@davidsneighbour](https://github.com/davidsneighbour) for their work on the documentation site.

Hugo now has:

* 37336+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 440+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 334+ [themes](http://themes.gohugo.io/)

## Notes

* All string sorting in Hugo is now alphabetical/lexicographical.
* `home.Pages` now only returns pages in the top level section. Before this release, it included _all regular pages in the site_. This made it easy to list all the pages on home page, but it also meant that you needed to take special care if you wanted to navigate the page tree from top to bottom. If you need _all regular pages_, use `.Site.RegularPages`.  Also see [#6153](https://github.com/gohugoio/hugo/issues/6153).
* `.Pages` now include sections. We have added `.RegularPages` as a convenience method if you want the old behaviour. See [#6154](https://github.com/gohugoio/hugo/issues/6154) for details.
* Hugo now only "auto create" sections for the home page and the top level folders. The other sections need a `_index.md` file. See [#6171](https://github.com/gohugoio/hugo/issues/6171) for details.


## Enhancements

### Templates

* Regenerate templates [2d1d3367](https://github.com/gohugoio/hugo/commit/2d1d33673d82c5073335e18944744606a71a5029) [@bep](https://github.com/bep) 
* Always load GitHub Gists over HTTPS [be0d4efc](https://github.com/gohugoio/hugo/commit/be0d4efc3db18035a04b188e089c09cdd8e04365) [@coliff](https://github.com/coliff) 

### Core

* Remove temporary warning [4644b95b](https://github.com/gohugoio/hugo/commit/4644b95bd568946429482aa36eeaff1eec6a7075) [@bep](https://github.com/bep) 
* Add some more site benchmarks [df374851](https://github.com/gohugoio/hugo/commit/df374851a0683f1446f33a4afef74c42f7d3eaaf) [@bep](https://github.com/bep) 

### Other

* Add FileInfo to resources created with resources.Match etc. [1089cfe4](https://github.com/gohugoio/hugo/commit/1089cfe4e1c35bec1f269b8280da43b367b5d070) [@bep](https://github.com/bep) [#6190](https://github.com/gohugoio/hugo/issues/6190)
* Improve the server assets cache invalidation logic [cd575023](https://github.com/gohugoio/hugo/commit/cd575023af846aa18ffa709f37bc70277e98cad3) [@bep](https://github.com/bep) [#6199](https://github.com/gohugoio/hugo/issues/6199)
* Do not fail build on errors in theme.toml [63150981](https://github.com/gohugoio/hugo/commit/6315098104ff80f8be6d5ae812835b4b4079582e) [@bep](https://github.com/bep) [#6162](https://github.com/gohugoio/hugo/issues/6162)
* Add resources.Match and resources.GetMatch [b64617fe](https://github.com/gohugoio/hugo/commit/b64617fe4f90da030bcf4a9c5a4913393ce96b14) [@bep](https://github.com/bep) [#6190](https://github.com/gohugoio/hugo/issues/6190)
* Convert from testify to quicktest [9e571827](https://github.com/gohugoio/hugo/commit/9e571827055dedb46b78c5db3d17d6913f14870b) [@bep](https://github.com/bep) 
* Avoid unnecessary conversions [6027ee11](https://github.com/gohugoio/hugo/commit/6027ee11082d0b9d72de1d4d1980a702be294ad2) [@muesli](https://github.com/muesli) 
* Simplify code [a93cbb0d](https://github.com/gohugoio/hugo/commit/a93cbb0d6cc6e3a78ba34aa372abc5b41ca24b2c) [@muesli](https://github.com/muesli) 
* Implement cascading front matter [bd98182d](https://github.com/gohugoio/hugo/commit/bd98182dbde893a8a809661c70633741bbf63911) [@bep](https://github.com/bep) [#6041](https://github.com/gohugoio/hugo/issues/6041)
* Use the SVG logo in README.md [c0eef3b4](https://github.com/gohugoio/hugo/commit/c0eef3b401615e85bb74baee6a515abcf531fc2c) [@XhmikosR](https://github.com/XhmikosR) 
* Add a branch bundle test case [82439520](https://github.com/gohugoio/hugo/commit/824395204680496d528684587a1f2977394aff3d) [@bep](https://github.com/bep) [#6173](https://github.com/gohugoio/hugo/issues/6173)
* Simplify page tree logic [7ff0a8ee](https://github.com/gohugoio/hugo/commit/7ff0a8ee9fe8d710d407e57faf1fda43bd635f28) [@bep](https://github.com/bep) [#6154](https://github.com/gohugoio/hugo/issues/6154)[#6153](https://github.com/gohugoio/hugo/issues/6153)[#6152](https://github.com/gohugoio/hugo/issues/6152)
* Merge pull request #6149 from bep/sort-caseinsensitive [53077b0d](https://github.com/gohugoio/hugo/commit/53077b0da54906feee64a03612e5186043e17341) [@bep](https://github.com/bep) 
* Regenerate CLI docs [02b947ea](https://github.com/gohugoio/hugo/commit/02b947eaa3cc68404180d796a2f7119dce074539) [@bep](https://github.com/bep) 
* Add "hugo config mounts" command [d7c233af](https://github.com/gohugoio/hugo/commit/d7c233afee6a16b1947f60b7e5450e40612997bb) [@bep](https://github.com/bep) [#6144](https://github.com/gohugoio/hugo/issues/6144)
* Cleanup the hugo config command [45ee8a7a](https://github.com/gohugoio/hugo/commit/45ee8a7a52213bf394c7f41a72be78084ddc789a) [@bep](https://github.com/bep) [#6144](https://github.com/gohugoio/hugo/issues/6144)
* Move the mount duplicate filter to the modules package [4b6c5eba](https://github.com/gohugoio/hugo/commit/4b6c5eba306e6e69f3dd07a6c102bfc8040b38c9) [@bep](https://github.com/bep) 
* Allow overlap in module mounts [edf9f0a3](https://github.com/gohugoio/hugo/commit/edf9f0a354e5eaa556f8faed70b5243b7273b35c) [@bep](https://github.com/bep) [#6146](https://github.com/gohugoio/hugo/issues/6146)
* Add some more content language test assertions [84bc8d84](https://github.com/gohugoio/hugo/commit/84bc8d84e4d2ec1fc94aee3113ebc570a28d1d16) [@bep](https://github.com/bep) [#6136](https://github.com/gohugoio/hugo/issues/6136)
* Add proper error message when receiving nil in Resource transformation [e5f96024](https://github.com/gohugoio/hugo/commit/e5f960245938d8d8b4e99f312e9907f8d3aebf7a) [@bep](https://github.com/bep) [#6128](https://github.com/gohugoio/hugo/issues/6128)
* Merge branch 'release-0.56.1' [9f497e7b](https://github.com/gohugoio/hugo/commit/9f497e7b5f77d0eb45d932a2301e648a3cd2d88f) [@bep](https://github.com/bep) 
* Update go-org to v0.1.2 [56908509](https://github.com/gohugoio/hugo/commit/56908509eb3a5779743a2314c05693a732b7feb3) [@niklasfasching](https://github.com/niklasfasching) 
* Do not return error on params dot access on incompatible types [e393c629](https://github.com/gohugoio/hugo/commit/e393c6290e827111a8a2e486791dc21f63a92b55) [@bep](https://github.com/bep) [#6121](https://github.com/gohugoio/hugo/issues/6121)
* Set GO111MODULE=on [e5fe3789](https://github.com/gohugoio/hugo/commit/e5fe378925c16c75902bbb46499c376c530ebdb5) [@bep](https://github.com/bep) [#6114](https://github.com/gohugoio/hugo/issues/6114)
* Skip resource cache init if the fs is missing [da4c4a77](https://github.com/gohugoio/hugo/commit/da4c4a7789d403af3f4f4fdd5dfd3327535e4050) [@bep](https://github.com/bep) [#6113](https://github.com/gohugoio/hugo/issues/6113)

## Fixes

### Core

* Fix output format handling of mix cased page kinds [de876242](https://github.com/gohugoio/hugo/commit/de87624241daa86660f205cc72a745409b9c9238) [@bep](https://github.com/bep) [#4528](https://github.com/gohugoio/hugo/issues/4528)
* Fix broken test [9ef4dca3](https://github.com/gohugoio/hugo/commit/9ef4dca361727a78e0f66f8f4e54c64e4c4781cb) [@bep](https://github.com/bep) 
* Fix bundle header clone logic [0e086785](https://github.com/gohugoio/hugo/commit/0e086785fa4be8086256e9d7de6cda78e18d00ee) [@bep](https://github.com/bep) [#6136](https://github.com/gohugoio/hugo/issues/6136)

### Other

* Fix faulty -h logic in hugo mod get [17ca8f0c](https://github.com/gohugoio/hugo/commit/17ca8f0c4c636752fb9da2ad551679275dc03dd3) [@bep](https://github.com/bep) [#6197](https://github.com/gohugoio/hugo/issues/6197)
* Fixed ineffectual assignments [c577a9ed](https://github.com/gohugoio/hugo/commit/c577a9ed2347559783c44232e1f08414008c5203) [@muesli](https://github.com/muesli) 
* Fixed tautological error conditions [e88d7989](https://github.com/gohugoio/hugo/commit/e88d7989907108b656eccd92bccc076be72a5c03) [@muesli](https://github.com/muesli) 
* Fix static sync issue with virtual mounts [166a394a](https://github.com/gohugoio/hugo/commit/166a394a2fef6f2990e264cc8dfb722af2cc6a67) [@bep](https://github.com/bep) [#6165](https://github.com/gohugoio/hugo/issues/6165)
* Cache the next position of `urlreplacer.prefix` [a843ca53](https://github.com/gohugoio/hugo/commit/a843ca53b5e0f29df9535fa0e88408a63cdc2cd7) [@vazrupe](https://github.com/vazrupe) [#5942](https://github.com/gohugoio/hugo/issues/5942)
* Fix no-map vs noMap discrepancy [02397e76](https://github.com/gohugoio/hugo/commit/02397e76cece28b467de30ff0cb0f471d9b212ee) [@bep](https://github.com/bep) [#6166](https://github.com/gohugoio/hugo/issues/6166)
* Fix assorted typos [f7f549e3](https://github.com/gohugoio/hugo/commit/f7f549e3a7492c787c6abb4900cc0f57c8ab1826) [@XhmikosR](https://github.com/XhmikosR) 
* common/collections: Fix typo [6512d128](https://github.com/gohugoio/hugo/commit/6512d128c6d33b86f376764ab1d622a89ea18d20) [@shawnps](https://github.com/shawnps) 
* Fix multilingual example compatibility with latest version [b8758de1](https://github.com/gohugoio/hugo/commit/b8758de19ec75b4565075314f9578270a092bc6f) [@robinwassen](https://github.com/robinwassen) 
* Fix self-mounts on the main project [36220851](https://github.com/gohugoio/hugo/commit/36220851e4ed7fc3fa78aa250d001d5f922210e7) [@bep](https://github.com/bep) [#6143](https://github.com/gohugoio/hugo/issues/6143)
* Fix config reloading in Vim and similar [6eca0a3d](https://github.com/gohugoio/hugo/commit/6eca0a3dee77f0e764b1de2e10c10ec2b7cf8ef1) [@bep](https://github.com/bep) [#6139](https://github.com/gohugoio/hugo/issues/6139)
* Fix Jekyll import [e28bd4c0](https://github.com/gohugoio/hugo/commit/e28bd4c0f843f39cfcb715b6c9c7d249bad5b500) [@bep](https://github.com/bep) [#6131](https://github.com/gohugoio/hugo/issues/6131)
* Fix image format detection for upper case extensions, e.g. JPG [c62bbf7b](https://github.com/gohugoio/hugo/commit/c62bbf7b11d68d52ef11a4c6c70660914c473d08) [@bep](https://github.com/bep) [#6137](https://github.com/gohugoio/hugo/issues/6137)
* Fix i18n project vs theme order [00a238e3](https://github.com/gohugoio/hugo/commit/00a238e32c82b0651e4145e306840cffa46e535d) [@bep](https://github.com/bep) [#6134](https://github.com/gohugoio/hugo/issues/6134)
* Fix image Width/Height regression [93d02aab](https://github.com/gohugoio/hugo/commit/93d02aabe6e611d65c428a9c5669b422e1bcf5e8) [@bep](https://github.com/bep) [#6120](https://github.com/gohugoio/hugo/issues/6120)





