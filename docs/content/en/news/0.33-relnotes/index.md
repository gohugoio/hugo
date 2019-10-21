
---
date: 2018-01-18
title: "Hugo 0.33: The New Kinder Surprise!"
description: "Hugo 0.33 comes with resource (images etc.) metadata, `type` and `layout` for all page types, `url` in front matter for list pages …"
categories: ["Releases"]
---

	Hugo `0.33` is the first main Hugo release of the new year, and it is safe to say that [@bep](https://github.com/bep)  has turned off his lazy Christmas mode :smiley:

This is a full makeover of the layout selection logic with full custom `layout` and `type` support (many have asked for this). Also, Hugo now respects the `url` value in front matter for all page types, including sections. Also, you can now configure `uglyURLs` per section.

But this release is also a follow-up to the `0.32` release which was all about bundles with resources and powerful image processing. With this release it is now simple to add metadata to your images and other bundle resources. 

[@bep](https://github.com/bep)  has added a section with examples of both `resources` configuration in both `YAML` and `TOML` front matter in his [test site](http://hugotest.bep.is/resourcemeta/). The example below shows a sample of how it would look like in `YAML`:

```yaml
date: 2017-01-17
title: My Bundle With YAML Resource Metadata
resources:
- src: "image-4.png"
  title: "The Fourth Image"
- src: "*.png"
  name: "my-cool-image-:counter"
  title: "The Image #:counter"
  params:
    byline: "bep"
```

This release represents **41 contributions by 3 contributors** to the main Hugo code base.

Hugo now has:

* 22553+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 448+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 197+ [themes](http://themes.gohugo.io/)

## Notes
* We have re-implemented and unified the template layout lookup logic. This has made it more powerful and much simpler to understand. We don't expect any sites to break because of this. We have tested lots of Hugo sites, including the 200 [themes](http://themes.gohugo.io/).
*  The `indexes` type is removed from template lookup. It's not in the documentation, and is a legacy term inherited from very old Hugo versions.
* If you have sub-dirs in your shiny new bundles (e.g. `my-bundle/images`) and use the `*Prefix*` methods to find them, we have made an unintended change that affects you. See [this issue](https://github.com/gohugoio/hugo/issues/4295).

## Enhancements

### Templates

* Respect `Type` and `Layout` for list template selection [51dd462c](https://github.com/gohugoio/hugo/commit/51dd462c3958f7cf032b06503f1f200a6aceebb9) [@bep](https://github.com/bep) [#3005](https://github.com/gohugoio/hugo/issues/3005)[#3245](https://github.com/gohugoio/hugo/issues/3245)

### Core

* Allow `url` in front matter for list type pages [8a409894](https://github.com/gohugoio/hugo/commit/8a409894bdb0972e152a2eccc47a2738568e1cfc) [@bep](https://github.com/bep) [#4263](https://github.com/gohugoio/hugo/issues/4263)
* Improve `.Site.GetPage` for regular translated pages. Before this change it was not possible to say "get me the current language edition of the given content page if possible." Now you can do that by doing a lookup without any extensions:  `.Site.GetPage "page" "post/mypost"` [9409bc0f](https://github.com/gohugoio/hugo/commit/9409bc0f799a8057836a14ccdf2833a55902175e) [@bep](https://github.com/bep) [#4285](https://github.com/gohugoio/hugo/issues/4285)
* Add front matter metadata to `Resource` [20c9b6ec](https://github.com/gohugoio/hugo/commit/20c9b6ec81171d1c586ea31d5d08b40b0edaffc6) [@bep](https://github.com/bep) [#4244](https://github.com/gohugoio/hugo/issues/4244)
* Implement `Resources.ByPrefix` [46db900d](https://github.com/gohugoio/hugo/commit/46db900dab9c0e6fcd9d227f10a32fb24f5c8bd9) [@bep](https://github.com/bep) [#4266](https://github.com/gohugoio/hugo/issues/4266)
* Make `GetByPrefix` work for Page resources [60c9f3b1](https://github.com/gohugoio/hugo/commit/60c9f3b1c34b69771e25a66906f150f460d73223) [@bep](https://github.com/bep) [#4264](https://github.com/gohugoio/hugo/issues/4264)
* Make `Resources.GetByPrefix` case insensitive [db85e834](https://github.com/gohugoio/hugo/commit/db85e83403913cff4b8737b138932b28e5bf6160) [@bep](https://github.com/bep) [#4258](https://github.com/gohugoio/hugo/issues/4258)
* Update `Chroma` and other third-party deps [64f0e9d1](https://github.com/gohugoio/hugo/commit/64f0e9d1c1d4ff2249fd9cf9749e70485002b36d) [@bep](https://github.com/bep) [#4267](https://github.com/gohugoio/hugo/issues/4267)
* Remove superflous `BuildDate` logic [13d53b31](https://github.com/gohugoio/hugo/commit/13d53b31f19240879122d6b7e4aaeb60b5130a3c) [@bep](https://github.com/bep) [#4272](https://github.com/gohugoio/hugo/issues/4272)
* Run benchmarks 3 times [b6ea6d07](https://github.com/gohugoio/hugo/commit/b6ea6d07d0b072d850fb066c78976acd6c2f5e81) [@bep](https://github.com/bep) 
* Support `uglyURLs` per section [57e10f17](https://github.com/gohugoio/hugo/commit/57e10f174e51cc5e1cf5f37eed30a0f3b153dd64) [@bep](https://github.com/bep) [#4256](https://github.com/gohugoio/hugo/issues/4256)
* Update CONTRIBUTING.md [1046e936](https://github.com/gohugoio/hugo/commit/1046e9363f2e382fd0b4aac838735ae4cbbebe5a) [@vassudanagunta](https://github.com/vassudanagunta) 
* Support offline builds [d5803da1](https://github.com/gohugoio/hugo/commit/d5803da1befba5446d1b2c1ad16f6467dc7b3991) [@vassudanagunta](https://github.com/vassudanagunta) 

## Fixes

* Fix handling of mixed-case taxonomy folders with content file [2d3189b2](https://github.com/gohugoio/hugo/commit/2d3189b22760e0a8995dae082a6bc5480f770bfe) [@bep](https://github.com/bep) [#4238](https://github.com/gohugoio/hugo/issues/4238)
* Fix handling of very long image file names [ecaf1451](https://github.com/gohugoio/hugo/commit/ecaf14514e06321823bdd10235cf23e7d654ba77) [@bep](https://github.com/bep) [#4261](https://github.com/gohugoio/hugo/issues/4261)
* Update `Afero` to avoid panic on "file name is too long" [f8a119b6](https://github.com/gohugoio/hugo/commit/f8a119b606d55aa4f31f16e5a3cadc929c99e4f8) [@bep](https://github.com/bep) [#4240](https://github.com/gohugoio/hugo/issues/4240)
* And now really fix the server watch logic [d4f8f88e](https://github.com/gohugoio/hugo/commit/d4f8f88e67f958b8010f90cb9b9854114e52dac2) [@bep](https://github.com/bep) [#4275](https://github.com/gohugoio/hugo/issues/4275)
* Fix server without watch [4e524ffc](https://github.com/gohugoio/hugo/commit/4e524ffcfff48c017717e261c6067416aa56410f) [@bep](https://github.com/bep) [#4275](https://github.com/gohugoio/hugo/issues/4275)






