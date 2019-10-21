
---
date: 2019-10-21
title: "Hugo 0.59.0"
description: "Set image target format and background color, and more ..."
categories: ["Releases"]
---

The timing of this release is motivated by getting the copies of the docs repositories in synch, now fully "Hugo Modularized". But it also comes with some very nice additions:

It is now possible to set the target format and the background fill color when processing images, e.g.:

```
{{ $image.Resize "600x jpg #b31280" }}
```

See [Image Processing Options](https://gohugo.io/content-management/image-processing/#image-processing-options).

Another useful addon is the `$pages.Next` and `$pages.Prev` methods on the core page collections in Hugo. These works the same way as the built-in static variants one `Page`, e.g. `.Next` and `.NextInSection`:

```
{{with .Site.RegularPages.Next . }}{{.RelPermalink}}{{end}}
```

The above is a functionally equivalent (but slightly slower) variant of:

```
{{with .Next }}{{.RelPermalink}}{{end}}
```

See [Pages Methods](https://gohugo.io/variables/pages/) for more information.


This release represents **45 contributions by 13 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@BaibhaVatsa](https://github.com/BaibhaVatsa), and [@XhmikosR](https://github.com/XhmikosR) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **34 contributions by 20 contributors**. A special thanks to [@bep](https://github.com/bep), [@celtic-coder](https://github.com/celtic-coder), [@napcs](https://github.com/napcs), and [@bmackinney](https://github.com/bmackinney) for their work on the documentation site.


Hugo now has:

* 38843+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 255+ [themes](http://themes.gohugo.io/)

## Notes


* Shortcode params can now be typed (supported types are `string`, `bool` `int` and `float64`, see [#6376](https://github.com/gohugoio/hugo/pull/6376).
* Pages.Next/.Prev as described above has existed for a long time, but they have been undocumented. They have been reimplemented for this release and now works like their namesakes on `Page`. This may be considered a breaking change, but it should be a welcome one, as the old behaviour wasn't very useful. See [#4500](https://github.com/gohugoio/hugo/issues/4500)

## Enhancements

### Templates

* Add optional "title" attribute to iframe in Vimeo shortcode [7b3edc29](https://github.com/gohugoio/hugo/commit/7b3edc293144dd450e87ca32f238221c21eb1b47) [@zbayoff](https://github.com/zbayoff) 
* Modify error messages of after, first, and last [65b7d422](https://github.com/gohugoio/hugo/commit/65b7d4221b90445bfc089873092411cf7e322933) [@BaibhaVatsa](https://github.com/BaibhaVatsa) [#6415](https://github.com/gohugoio/hugo/issues/6415)
* Last now accepts 0 as limit [0e75af74](https://github.com/gohugoio/hugo/commit/0e75af74db30259ec355a7b79a1e257d5fe00eef) [@BaibhaVatsa](https://github.com/BaibhaVatsa) [#6419](https://github.com/gohugoio/hugo/issues/6419)
* After now accepts 0 as index [096a4b67](https://github.com/gohugoio/hugo/commit/096a4b67b98259dabff5ebfbfd879a41999a1ed2) [@BaibhaVatsa](https://github.com/BaibhaVatsa) [#6388](https://github.com/gohugoio/hugo/issues/6388)
* Make getJSON/getCVS accept non-string args [0d7b05be](https://github.com/gohugoio/hugo/commit/0d7b05be4cb2391cbd280f6109c01ec2d3d7e0c6) [@bep](https://github.com/bep) [#6382](https://github.com/gohugoio/hugo/issues/6382)
* Add `rel="noopener"` for external links [34dc06b0](https://github.com/gohugoio/hugo/commit/34dc06b032741abac342d7a2a77510ded9b72ae8) [@XhmikosR](https://github.com/XhmikosR) 
* Remove unneeded space [2b1814ee](https://github.com/gohugoio/hugo/commit/2b1814ee580f3149f9fe0a4cf30b754bac9f0c90) [@XhmikosR](https://github.com/XhmikosR) 
* Remove eq argument limitation [5e660947](https://github.com/gohugoio/hugo/commit/5e660947757023434dd7a1ec8b8239c0577fd501) [@vazrupe](https://github.com/vazrupe) [#6237](https://github.com/gohugoio/hugo/issues/6237)

### Output

* Add common video media types [689f647b](https://github.com/gohugoio/hugo/commit/689f647baf96af078186f0cdc45199f7d0995d22) [@martignoni](https://github.com/martignoni) 
* Simplify test output to simplify diffing [339ee371](https://github.com/gohugoio/hugo/commit/339ee37143ca5a6bb22bbc1b0468d785f450cfb7) [@bep](https://github.com/bep) 
* Use + to create the Type string [64ec8c89](https://github.com/gohugoio/hugo/commit/64ec8c89049461c4731b23c491fb41e00a09a8b2) [@bep](https://github.com/bep) 
* Support output image format in image operations [e5856e61](https://github.com/gohugoio/hugo/commit/e5856e61d88ef5149582851b00e06b5b93dce9f8) [@jansorg](https://github.com/jansorg) [#6298](https://github.com/gohugoio/hugo/issues/6298)

### Other

* Replace /docs [39121de4](https://github.com/gohugoio/hugo/commit/39121de4d991bdcf5f202da4d8d81a8ac6c149fc) [@bep](https://github.com/bep) 
* Recover from file corruption [180195aa](https://github.com/gohugoio/hugo/commit/180195aa342777fece1b29a08ec89456d7996c61) [@bep](https://github.com/bep) [#6401](https://github.com/gohugoio/hugo/issues/6401)
* Allow to set background fill colour [4b286b9d](https://github.com/gohugoio/hugo/commit/4b286b9d2722909d0682e50eeecdfe16c1f47fd8) [@bep](https://github.com/bep) [#6298](https://github.com/gohugoio/hugo/issues/6298)
* Replace .RSSLink [46cafdba](https://github.com/gohugoio/hugo/commit/46cafdbaca13866f32db04c0cc28374e30ec5914) [@bep](https://github.com/bep) [#6037](https://github.com/gohugoio/hugo/issues/6037)
* Use binary search in Pages.Prev/Next if possible [653e6856](https://github.com/gohugoio/hugo/commit/653e6856ea1cfc60cc16733807d23b302dbe4bd5) [@bep](https://github.com/bep) [#4500](https://github.com/gohugoio/hugo/issues/4500)
* Make Pages.Prev/Next work like the other Prev/Next methods [f4f566ed](https://github.com/gohugoio/hugo/commit/f4f566edf4bd6a590cf9cdbd5cfc0026ecd93b14) [@bep](https://github.com/bep) [#4500](https://github.com/gohugoio/hugo/issues/4500)
* Update feature_request.md [5f1aafaf](https://github.com/gohugoio/hugo/commit/5f1aafafb40299bb4c8aebf71e05843431eb84c5) [@bep](https://github.com/bep) 
* Update to Go 1.12.10 and 1.13.1 [71b18a07](https://github.com/gohugoio/hugo/commit/71b18a0786894893eafa01263a0915149ed303ec) [@bep](https://github.com/bep) [#6406](https://github.com/gohugoio/hugo/issues/6406)
* Add FileMeta.String [f10db101](https://github.com/gohugoio/hugo/commit/f10db101a18f5cad332c9398136f77e35a169d52) [@bep](https://github.com/bep) 
* Update minify to v2.5.2 [b401858e](https://github.com/gohugoio/hugo/commit/b401858ebd346c433dd69a260eba7098bded5a30) [@anthonyfok](https://github.com/anthonyfok) 
* Add BaseFs to RenderingContext [020a6fbd](https://github.com/gohugoio/hugo/commit/020a6fbd7f6996ed84d80ba6c37fe0d8c2536806) [@niklasfasching](https://github.com/niklasfasching) 
* Update go-org [b152216d](https://github.com/gohugoio/hugo/commit/b152216d5c8adbf1bfa4c6fb7b2a50b6866c685e) [@niklasfasching](https://github.com/niklasfasching) 
* Upgrade to latest version of emoji dependency [c466b88c](https://github.com/gohugoio/hugo/commit/c466b88c998bc99e5d26e41cb67d87e1d4b976f5) [@jamietanna](https://github.com/jamietanna) [#6391](https://github.com/gohugoio/hugo/issues/6391)
* Upgrade to latest version of emoji dependency [170f18d9](https://github.com/gohugoio/hugo/commit/170f18d9352d39213170dd9d5e947eb45854c84b) [@jamietanna](https://github.com/jamietanna) 
* Update Architectures [15a0364d](https://github.com/gohugoio/hugo/commit/15a0364d39741da34b8661f9a8386b54016049d6) [@bep](https://github.com/bep) 
* Add ability to invalidate Google Cloud CDN [674e81ae](https://github.com/gohugoio/hugo/commit/674e81ae8700bdd00d3e5e47ff930d42d25bc68b) [@gkelly](https://github.com/gkelly) 
* Ensure same dirinfos sort order in TestImageOperationsGolden [298092d5](https://github.com/gohugoio/hugo/commit/298092d516f623cc20051f506d460fb7625cdc84) [@anthonyfok](https://github.com/anthonyfok) 
* Update bug_report.md [019ae384](https://github.com/gohugoio/hugo/commit/019ae384835446266b951875aa0870d245382cf2) [@bep](https://github.com/bep) 
* Support typed bool, int and float in shortcode params [329e88db](https://github.com/gohugoio/hugo/commit/329e88db1f6d043d32c7083570773dccfd4f11fc) [@bep](https://github.com/bep) [#6371](https://github.com/gohugoio/hugo/issues/6371)
* Update Chroma [e073f4ef](https://github.com/gohugoio/hugo/commit/e073f4efb1345f6408000ef3f389873f8cf7179e) [@bep](https://github.com/bep) [#6279](https://github.com/gohugoio/hugo/issues/6279)
* Add issue templates and action [454a033d](https://github.com/gohugoio/hugo/commit/454a033dc5bc9b3db626fe1533d7e8494d79f472) [@bmackinney](https://github.com/bmackinney) 
* Add some more resource transform tests [c262a95a](https://github.com/gohugoio/hugo/commit/c262a95a5c5a9304c82b9d9e39701bc471916851) [@bep](https://github.com/bep) [#6348](https://github.com/gohugoio/hugo/issues/6348)
* Do not compile in Azure on Solaris [c0d7188e](https://github.com/gohugoio/hugo/commit/c0d7188ec85e7a4b61489e38896108d877f6d902) [@fazalmajid](https://github.com/fazalmajid) [#6324](https://github.com/gohugoio/hugo/issues/6324)
* Ignore "does not exist" errors in prune [fcfa6f33](https://github.com/gohugoio/hugo/commit/fcfa6f33bbebc128a3f9bc3162173bc3780c5f50) [@bep](https://github.com/bep) [#6326](https://github.com/gohugoio/hugo/issues/6326)[#5745](https://github.com/gohugoio/hugo/issues/5745)
* Avoid writing the same processed image to /public twice [9442937d](https://github.com/gohugoio/hugo/commit/9442937d82005b369780edcc557e0d15d6bf0bad) [@bep](https://github.com/bep) [#6307](https://github.com/gohugoio/hugo/issues/6307)
* Update github.com/bep/gitmap [24ad4295](https://github.com/gohugoio/hugo/commit/24ad4295718341dcae12b72bf52fef312d1036ed) [@bep](https://github.com/bep) [#6313](https://github.com/gohugoio/hugo/issues/6313)

## Fixes

### Core

* Fix broken bundle live reload logic [901077c0](https://github.com/gohugoio/hugo/commit/901077c0364eaf3fe4f997c3026aa18cfc7781ed) [@bep](https://github.com/bep) [#6315](https://github.com/gohugoio/hugo/issues/6315)[#6308](https://github.com/gohugoio/hugo/issues/6308)

### Other

* Fix elements are doubling when append a not assignable type [a9762b5c](https://github.com/gohugoio/hugo/commit/a9762b5c48054e036332eff541a8fd32e54ada13) [@vazrupe](https://github.com/vazrupe) [#6188](https://github.com/gohugoio/hugo/issues/6188)
* Fix data race in global logger init [bc70f2bf](https://github.com/gohugoio/hugo/commit/bc70f2bf123d94fc3226754ec9f1f44748e98162) [@bep](https://github.com/bep) [#6409](https://github.com/gohugoio/hugo/issues/6409)
* Fix image test error on s390x, ppc64* and arm64 [39ed33fc](https://github.com/gohugoio/hugo/commit/39ed33fcebcde91605e645fd28fd94020b442d97) [@anthonyfok](https://github.com/anthonyfok) [#6387](https://github.com/gohugoio/hugo/issues/6387)
* Fix cache key transformed resources [6dec671f](https://github.com/gohugoio/hugo/commit/6dec671fb930029e18ba9aa5135b3a27adcddb21) [@bep](https://github.com/bep) [#6348](https://github.com/gohugoio/hugo/issues/6348)
* Fix cache keys for bundled resoures in transform.Unmarshal [c0d75736](https://github.com/gohugoio/hugo/commit/c0d7573677e9726c14749ccd432dccb75e0d194d) [@bep](https://github.com/bep) [#6327](https://github.com/gohugoio/hugo/issues/6327)
* Fix concat with fingerprint regression [3be2c253](https://github.com/gohugoio/hugo/commit/3be2c25351b421a26ee1ff2a38cbab00280c0583) [@bep](https://github.com/bep) [#6309](https://github.com/gohugoio/hugo/issues/6309)





