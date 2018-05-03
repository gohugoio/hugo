
---
date: 2017-08-07
title: "Hugo 0.26: Language Style Edition"
description: "Hugo 0.26 brings proper AP Style or Chicago Style Title Case, « French Guillemets » and more."
categories: ["Releases"]
images:
- images/blog/hugo-26-poster.png
---

This release brings a choice of **AP Style or Chicago Style Title Case** ([8fb594bf](https://github.com/gohugoio/hugo/commit/8fb594bfb090c017d4e5cbb2905780221e202c41) [#989](https://github.com/gohugoio/hugo/issues/989)). You can also now configure Blackfriday to render **« French Guillemets »** ([cb9dfc26](https://github.com/gohugoio/hugo/commit/cb9dfc2613ae5125cafa450097fb0f62dd3770e7) [#3725](https://github.com/gohugoio/hugo/issues/3725)). To enable French Guillemets, put this in your site `config.toml`:


```bash
[blackfriday]
angledQuotes = true
smartypantsQuotesNBSP = true
```

Oh, and this release also fixes it so you should see no ugly long crashes no more when you step wrong in your templates ([794ea21e](https://github.com/gohugoio/hugo/commit/794ea21e9449b876c5514f1ce8fe61449bbe4980)).

Hugo `0.26` represents **46 contributions by 11 contributors** to the main Hugo code base.

[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@jorinvo](https://github.com/jorinvo), and [@digitalcraftsman](https://github.com/digitalcraftsman) for their ongoing contributions. And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **838 contributions by 30 contributors**. A special thanks to [@rdwatters](https://github.com/rdwatters), [@bep](https://github.com/bep),  [@digitalcraftsman](https://github.com/digitalcraftsman), and  [@budparr](https://github.com/budparr) for their work on the documentation site.

This may look like a **Waiting Sausage**, a barbecue term used in Norway for that sausage you eat while waiting for the steak to get ready. And it is: We're working on bigger and even more interesting changes behind the scenes. Stay tuned!

Hugo now has:

* 18802+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 457+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 175+ [themes](http://themes.gohugo.io/)

## Notes

* `sourceRelativeLinks` has been deprecated for a while and has now been removed. [9891c0fb](https://github.com/gohugoio/hugo/commit/9891c0fb0eb274b8a95b62c40070a87a6e04088c) [@bep](https://github.com/bep) [#3766](https://github.com/gohugoio/hugo/issues/3766)
* The `title` template function and taxonomy page titles now default to following the [AP Stylebook](https://www.apstylebook.com/) for title casing.  To override this default to use the old behavior, set `titleCaseStyle` to `Go` in your site configuration. [8fb594bf](https://github.com/gohugoio/hugo/commit/8fb594bfb090c017d4e5cbb2905780221e202c41) [@bep](https://github.com/bep) [#989](https://github.com/gohugoio/hugo/issues/989)

## Enhancements

### Templates

* Use hash for cache key [6cd33f69](https://github.com/gohugoio/hugo/commit/6cd33f6953671edb13d42dcb15746bd10df3428b) [@RealOrangeOne](https://github.com/RealOrangeOne) [#3690](https://github.com/gohugoio/hugo/issues/3690)
* Add some empty slice tests to intersect [e0cf2e05](https://github.com/gohugoio/hugo/commit/e0cf2e05bbdcb8b4a3f875df84a878f4ca80e904) [@bep](https://github.com/bep) [#3686](https://github.com/gohugoio/hugo/issues/3686)

### Core

* Support `reflinks` starting with a slash [dbe63970](https://github.com/gohugoio/hugo/commit/dbe63970e09313dec287816ab070b5c2f5a13b1b) [@bep](https://github.com/bep) [#3703](https://github.com/gohugoio/hugo/issues/3703)
* Make template panics into nice error messages [794ea21e](https://github.com/gohugoio/hugo/commit/794ea21e9449b876c5514f1ce8fe61449bbe4980) [@bep](https://github.com/bep)

### Other

* Make the `title` case style guide configurable [8fb594bf](https://github.com/gohugoio/hugo/commit/8fb594bfb090c017d4e5cbb2905780221e202c41) [@bep](https://github.com/bep) [#989](https://github.com/gohugoio/hugo/issues/989)
* Add support for French Guillemets [cb9dfc26](https://github.com/gohugoio/hugo/commit/cb9dfc2613ae5125cafa450097fb0f62dd3770e7) [@bep](https://github.com/bep) [#3725](https://github.com/gohugoio/hugo/issues/3725)
* Add support for French Guillemets [c4a0b6e8](https://github.com/gohugoio/hugo/commit/c4a0b6e8abdf9f800fbd7a7f89e9f736edc60431) [@bep](https://github.com/bep) [#3725](https://github.com/gohugoio/hugo/issues/3725)
* Switch from fork bep/inflect to markbates/inflect [09907d36](https://github.com/gohugoio/hugo/commit/09907d36af586c5b29389312f2ecc2962c06313c) [@jorinvo](https://github.com/jorinvo)
* Remove unused dependencies from vendor.json [9b4170ce](https://github.com/gohugoio/hugo/commit/9b4170ce768717adfbe9d97c46e38ceaec2ce994) [@jorinvo](https://github.com/jorinvo)
* Add `--debug` option to be improved on over time [aee2b067](https://github.com/gohugoio/hugo/commit/aee2b06780858c12d8cb04c7b1ba592543410aa9) [@maxandersen](https://github.com/maxandersen)
* Reduce Docker image size from 277MB to 27MB [bfe0bfbb](https://github.com/gohugoio/hugo/commit/bfe0bfbbd1a59ddadb72a6b07fecce71716088ec) [@ellerbrock](https://github.com/ellerbrock) [#3730](https://github.com/gohugoio/hugo/issues/3730)[#3738](https://github.com/gohugoio/hugo/issues/3738)
* Optimize Docker image size [606d6a8c](https://github.com/gohugoio/hugo/commit/606d6a8c9177dda4551ed198e0aabbe569f0725d) [@ellerbrock](https://github.com/ellerbrock) [#3674](https://github.com/gohugoio/hugo/issues/3674)
* Add `--trace` to asciidoctor args [b60aa1a5](https://github.com/gohugoio/hugo/commit/b60aa1a504f3fbf9c19a6bf2030fdc7a04ab4a5a) [@miltador](https://github.com/miltador) [#3714](https://github.com/gohugoio/hugo/issues/3714)
* Add script to pull in docs changes [ff433f98](https://github.com/gohugoio/hugo/commit/ff433f98133662063cbb16e220fd44c678c82823) [@bep](https://github.com/bep)
* Add `HasShortcode` [deccc540](https://github.com/gohugoio/hugo/commit/deccc54004cbe88ddbf8f3f951d3178dc0693189) [@bep](https://github.com/bep) [#3707](https://github.com/gohugoio/hugo/issues/3707)
* Improve the twitter card template [00b590d7](https://github.com/gohugoio/hugo/commit/00b590d7ab4f3021814acceaf74c4eaf64edb226) [@bep](https://github.com/bep) [#3711](https://github.com/gohugoio/hugo/issues/3711)
* Add `GOEXE` to support building with different versions of `go` [ea5e9e34](https://github.com/gohugoio/hugo/commit/ea5e9e346c93320538c6517b619b5f57473291c8) [@mdhender](https://github.com/mdhender)

## Fixes

### Templates

* Fix intersect on `[]interface{}` handling [55d0b894](https://github.com/gohugoio/hugo/commit/55d0b89417651eba3ae51c96bd9de9e0daa0399e) [@moorereason](https://github.com/moorereason) [#3718](https://github.com/gohugoio/hugo/issues/3718)

### Other

* Fix broken `TaskList` in Markdown [481924b3](https://github.com/gohugoio/hugo/commit/481924b34d23b0ce435778cce7bce77571b22f9d) [@mpcabd](https://github.com/mpcabd) [#3710](https://github.com/gohugoio/hugo/issues/3710)



