
---
date: 2019-09-04
title: "Image Processing Galore!"
description: "Hugo 0.58 adds the long sought after Exif method plus many useful image filters. And it's faster ..."
categories: ["Releases"]
---

**Hugo 0.58** adds the long sought after [Exif (docs)](https://gohugo.io/content-management/image-processing/#exif)  method on image and a bunch of useful [image filters (docs)](https://gohugo.io/functions/images/#image-filters), courtesy of [@disintegration](https://github.com/disintegration)'s great [Gift](https://github.com/disintegration/gift) image library.

This means that you now can do variations of this:

```go-html-template
{{ $blurryGrayscale := $myimage.Resize "300x200" | images.Filter images.Grayscale (images.GaussianBlur 8) }}
{{ $exif := $myimg.Exif }}
```

It's worth noting that the issue that enabled/triggered the implementation of the above was the simplifications needed to fix [#5903](https://github.com/gohugoio/hugo/issues/5903), which makes sure that type information is preserved when processed via **Hugo Pipes**. E.g. you can now do:

```go-html-template
{{ ($myimg | fingerprint ).Width }}
```

And it works as expected.

This release is also built with the brand new **Go 1.13** which means that it's also the [fastest Hugo version](https://discourse.gohugo.io/t/hugo-benchmarks-go-1-12-vs-go-1-13/20572/5) to date.

This release represents **39 contributions by 5 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@niklasfasching](https://github.com/niklasfasching), [@vazrupe](https://github.com/vazrupe), and [@jakejarvis](https://github.com/jakejarvis) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **8 contributions by 8 contributors**. A special thanks to [@jacebenson](https://github.com/jacebenson), [@digitalcraftsman](https://github.com/digitalcraftsman), [@jernst](https://github.com/jernst), and [@rgwood](https://github.com/rgwood) for their work on the documentation site.


Hugo now has:

* 37859+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 440+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 317+ [themes](http://themes.gohugo.io/)

## Notes

* `home.Pages` now behaves like all the other sections, see [#6240](https://github.com/gohugoio/hugo/issues/6240). If you want to list all the regular pages, use `.Site.RegularPages`.
* We have added some new image filters to Hugo's image processing. This also means that we have consolidated the resize operations to use the one `gift` library (from the same developer as the one we used before). The operations work as before, but one difference is that we no longer embed color profile information in PNG images, but this should also be a more portable solution. Software that supports color profiles will assume that images without an embedded profile are in the sRGB profile. Software that doesn't support color profiles will use the monitor's profile, which is most likely to be sRGB as well.
* We have improved the file cache logic for processed images and only stores them once when the same image is bundled in multiple languages. This means that you may want to run `hugo --gc` to clean your image cache.

## Enhancements

### Templates

* Migrate last shortcodes (YouTube and Vimeo) to HTTPS embeds [00297085](https://github.com/gohugoio/hugo/commit/00297085db48cbb7949c9867012f6df38817fc29) [@jakejarvis](https://github.com/jakejarvis) 
* Use RegularPages for RSS template [88d69936](https://github.com/gohugoio/hugo/commit/88d69936122f82fffc02850516bdb37be3d0892b) [@bep](https://github.com/bep) [#6238](https://github.com/gohugoio/hugo/issues/6238)
* Avoid "home page warning" in RSS template [564cf1bb](https://github.com/gohugoio/hugo/commit/564cf1bb11e100891992e9131b271a79ea7fc528) [@bep](https://github.com/bep) [#6238](https://github.com/gohugoio/hugo/issues/6238)

### Core

* Adjust Go version specific test [dc3f3df2](https://github.com/gohugoio/hugo/commit/dc3f3df29d2b65532cedc9d321db7c4a38a28d7d) [@bep](https://github.com/bep) [#6304](https://github.com/gohugoio/hugo/issues/6304)
* Remove the old and slow site benchmarks [28501ceb](https://github.com/gohugoio/hugo/commit/28501ceb93613729c5971105010dd3c22cfa0f7f) [@bep](https://github.com/bep) 
* Add a Sass includePaths test [1b5c7e32](https://github.com/gohugoio/hugo/commit/1b5c7e327c7f98cf8e9fff920f3328198f67a598) [@bep](https://github.com/bep) [#6274](https://github.com/gohugoio/hugo/issues/6274)
* Change to output non-panic error message if missing shortcode template [fd3d90ce](https://github.com/gohugoio/hugo/commit/fd3d90ced85baaf6941be45b2fe29c25ff755c18) [@vazrupe](https://github.com/vazrupe) [#6075](https://github.com/gohugoio/hugo/issues/6075)
* Don't use the global warning logger [ea681603](https://github.com/gohugoio/hugo/commit/ea6816030081b2cffa6c0ae9ca5429a2c6fe2fa5) [@bep](https://github.com/bep) [#6238](https://github.com/gohugoio/hugo/issues/6238)
* Allow index.md inside bundles [4b4bdcfe](https://github.com/gohugoio/hugo/commit/4b4bdcfe740d988e4cfb4fee53eced6985576abd) [@bep](https://github.com/bep) [#6208](https://github.com/gohugoio/hugo/issues/6208)
* Add a site benchmark [416493b5](https://github.com/gohugoio/hugo/commit/416493b548a9bbaa27758fba9bab50a22b680e9d) [@bep](https://github.com/bep) 
* Recover and log panics in content init [7f3aab5a](https://github.com/gohugoio/hugo/commit/7f3aab5ac283ecfc7029b680d4c0a34920e728c8) [@bep](https://github.com/bep) [#6210](https://github.com/gohugoio/hugo/issues/6210)
* Add some outputs tests [028b9926](https://github.com/gohugoio/hugo/commit/028b992611209b241b1f55def8d47f9188038dc3) [@bep](https://github.com/bep) [#6210](https://github.com/gohugoio/hugo/issues/6210)

### Other

* Update to Go 1.13 [b4313011](https://github.com/gohugoio/hugo/commit/b43130115d9e3888d94df9e6f5fc72eba662632f) [@bep](https://github.com/bep) [#6304](https://github.com/gohugoio/hugo/issues/6304)
* Cache processed images by their source path [8624b9fe](https://github.com/gohugoio/hugo/commit/8624b9fe9eb81aeb884d36311fb6f85fed98aa43) [@bep](https://github.com/bep) [#6269](https://github.com/gohugoio/hugo/issues/6269)
* Remove test artifact [018494f3](https://github.com/gohugoio/hugo/commit/018494f363a32b9e4d3622da6842bc3e59b420b2) [@bep](https://github.com/bep) 
* Make the "is this a Hugo Module" logic more lenient [43298f02](https://github.com/gohugoio/hugo/commit/43298f028ccdf38e949b573d03d328bf96b998a3) [@bep](https://github.com/bep) [#6299](https://github.com/gohugoio/hugo/issues/6299)
* Update to Go 1.11.13 and 1.12.9 [05d83b6c](https://github.com/gohugoio/hugo/commit/05d83b6c08089c20ca1d99bcd224188ed5d127d4) [@bep](https://github.com/bep) [#6228](https://github.com/gohugoio/hugo/issues/6228)
* Make home.Pages work like any other section [4898fb3d](https://github.com/gohugoio/hugo/commit/4898fb3d64c856c5e0f324e0dfbf3b60da1d1d3a) [@bep](https://github.com/bep) [#6240](https://github.com/gohugoio/hugo/issues/6240)
* Add some fingerprint tests [45d7988f](https://github.com/gohugoio/hugo/commit/45d7988f2d0aa95d1a56f4c66342574075cf2963) [@bep](https://github.com/bep) [#6284](https://github.com/gohugoio/hugo/issues/6284)[#6280](https://github.com/gohugoio/hugo/issues/6280)
* Cache Exif data to disk [ce47c21a](https://github.com/gohugoio/hugo/commit/ce47c21a2998630f8edcbd056983d9c59a80b676) [@bep](https://github.com/bep) [#6291](https://github.com/gohugoio/hugo/issues/6291)
* Remove metaDataFormat setting [de9cbf61](https://github.com/gohugoio/hugo/commit/de9cbf61954201943a7b170a7d0a8b34afb5942c) [@bep](https://github.com/bep) 
* Make the Exif benchmark filenames distinct [4f501169](https://github.com/gohugoio/hugo/commit/4f5011692a22762e213e872fd9e39d015141083f) [@bep](https://github.com/bep) 
* Add Exif benchmark [3becba7a](https://github.com/gohugoio/hugo/commit/3becba7a982f39f67c7ee7cff411eae50931c8cd) [@bep](https://github.com/bep) [#6291](https://github.com/gohugoio/hugo/issues/6291)
* Remove unused map type [20bdc69a](https://github.com/gohugoio/hugo/commit/20bdc69a47b851871bdc4d9be6366fa7f51f25db) [@bep](https://github.com/bep) 
* Add image.Exif [28143397](https://github.com/gohugoio/hugo/commit/28143397d625cce1f89f4161cba97c0dddd9004c) [@bep](https://github.com/bep) [#4600](https://github.com/gohugoio/hugo/issues/4600)
* Add a set of image filters [823f53c8](https://github.com/gohugoio/hugo/commit/823f53c861bb49aecc6104e0add39fc3b0729025) [@bep](https://github.com/bep) [#6255](https://github.com/gohugoio/hugo/issues/6255)
* Image resource refactor [f9978ed1](https://github.com/gohugoio/hugo/commit/f9978ed16476ca6d233a89669c62c798cdf9db9d) [@bep](https://github.com/bep) [#5903](https://github.com/gohugoio/hugo/issues/5903)[#6234](https://github.com/gohugoio/hugo/issues/6234)[#6266](https://github.com/gohugoio/hugo/issues/6266)
* Remove debug check left during development [ad1d6d64](https://github.com/gohugoio/hugo/commit/ad1d6d6406c9b208d4fd4e09d6ad9ef19aa65dbb) [@bep](https://github.com/bep) [#6249](https://github.com/gohugoio/hugo/issues/6249)
* Adjust the default paginator for sections [18836a71](https://github.com/gohugoio/hugo/commit/18836a71ce7b671fa71dd1318b99fc661755e94d) [@bep](https://github.com/bep) [#6231](https://github.com/gohugoio/hugo/issues/6231)
* Update to Go 1.11.13 and 1.12.9 [f28efd35](https://github.com/gohugoio/hugo/commit/f28efd35820dc4909832c14dfd8ea6812ecead31) [@bep](https://github.com/bep) [#6228](https://github.com/gohugoio/hugo/issues/6228)
* Disable "auto tidy" for now [321418f2](https://github.com/gohugoio/hugo/commit/321418f22a4a94b87f01e1403a2f4a71106461fb) [@bep](https://github.com/bep) [#6115](https://github.com/gohugoio/hugo/issues/6115)
* Make sure the hugo field is always initialized before it's used [ea9261e8](https://github.com/gohugoio/hugo/commit/ea9261e856c13c1d4ae05fcca08766d410b4b65c) [@vazrupe](https://github.com/vazrupe) [#6193](https://github.com/gohugoio/hugo/issues/6193)

## Fixes

### Core

* Fix draft etc. handling of _index.md pages [6ccf50ea](https://github.com/gohugoio/hugo/commit/6ccf50ea7bb291bcbe1d56a4d697a6fd57a9c629) [@bep](https://github.com/bep) [#6222](https://github.com/gohugoio/hugo/issues/6222)[#6210](https://github.com/gohugoio/hugo/issues/6210)
* Fix taxonomies vs expired [9475f61a](https://github.com/gohugoio/hugo/commit/9475f61a377fcf23f910cbfd4ddca59261326665) [@bep](https://github.com/bep) [#6213](https://github.com/gohugoio/hugo/issues/6213)

### Other

* Update go-org (fix descriptive lists) [8a8d4a6d](https://github.com/gohugoio/hugo/commit/8a8d4a6d97d181f1aaee639d35b198a27bb788e2) [@niklasfasching](https://github.com/niklasfasching) 
* Update go-org (fix footnotes in headlines) [58d4c0a8](https://github.com/gohugoio/hugo/commit/58d4c0a8be8beefbd7437b17bf7a9a381164d09b) [@niklasfasching](https://github.com/niklasfasching) 
* Discrepancy typo fix [c5319db9](https://github.com/gohugoio/hugo/commit/c5319db9f13f1dee97db5fbbeae38429a074c7d0) [@coliff](https://github.com/coliff) 
* Fix mainSections logic [67524c99](https://github.com/gohugoio/hugo/commit/67524c993623871626f0f22e6a2ac705a816a959) [@bep](https://github.com/bep) [#6217](https://github.com/gohugoio/hugo/issues/6217)
* Fix live reload mount logic with sub paths [952a3194](https://github.com/gohugoio/hugo/commit/952a3194962dd91f87e5bd227a1591b00c39ff05) [@bep](https://github.com/bep) [#6209](https://github.com/gohugoio/hugo/issues/6209)





