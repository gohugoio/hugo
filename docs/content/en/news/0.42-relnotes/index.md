
---
date: 2018-06-12
title: "Hugo 0.42: Theme Composition and Inheritance!"
description: "Hugo 0.42 adds Theme Components support, a new and powerful way of composing your Hugo sites."
categories: ["Releases"]
---

	Hugo `0.42` adds **Theme Components**. This is a really powerful new way of building your Hugo sites with reusable components. This is both **Theme Composition** and **Theme Inheritance**.

>Read more about Theme Components in the [Hugo Documentation](https://gohugo.io/themes/theme-components/).

The feature above was implemented by [@bep](https://github.com/bep), the funny Norwegian, with great design help from the Hugo community. But that implementation would not have been possible without [Afero](https://github.com/spf13/afero), [Steve Francia's](https://github.com/spf13) virtual file system. Hugo is built on top of Afero and many other fast and solid Go projects, and if you look at the combined contribution log of the Hugo project and its many open source dependencies, the total amount of contributions is staggering. A big thank you to them all!

This release represents **27 contributions by 7 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@onedrawingperday](https://github.com/onedrawingperday), [@anthonyfok](https://github.com/anthonyfok), and [@stefanneuhaus](https://github.com/stefanneuhaus) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **8 contributions by 5 contributors**. A special thanks to [@bep](https://github.com/bep), [@LorenzCK](https://github.com/LorenzCK), [@gavinwray](https://github.com/gavinwray), and [@deyton](https://github.com/deyton) for their work on the documentation site.


Hugo now has:

* 26286+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 444+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 235+ [themes](http://themes.gohugo.io/)

## Notes
* The `speakerdeck` shortcode is removed. It didn't work properly, so is isn't expected to be in wide use. If you use it, you will get a build error. You will either have to remove its usage or add a `speakerdeck` shortcode to your own project or theme.
* We have now virtualized the filesystems for project and theme files. This makes everything simpler, faster and more powerful. But it also means that template lookups on the form `{{ template “theme/partials/pagination.html” . }}` will not work anymore. That syntax has never been documented, so it's not expected to be in wide use. `{{ partial "pagination.html" . }}` will give you the most specific version of that partial.

## Enhancements


* Always load GA script over HTTPS [2e6712e2](https://github.com/gohugoio/hugo/commit/2e6712e2814f333caa807888c1d8a9a5a3c03709) [@coliff](https://github.com/coliff) 
* Remove speakerdeck shortcode [65deb72d](https://github.com/gohugoio/hugo/commit/65deb72dc4c9299416cf2d9defddb96dba4101fd) [@onedrawingperday](https://github.com/onedrawingperday) [#4830](https://github.com/gohugoio/hugo/issues/4830)
* Add `strings.RuneCount` [019bd557](https://github.com/gohugoio/hugo/commit/019bd5576be87c9f06b6a928ede1a5e78677f7b3) [@theory](https://github.com/theory) 
* Add `strings.Repeat` [13435a6f](https://github.com/gohugoio/hugo/commit/13435a6f608306c5094fdcd72a1d9538727f91b2) [@theory](https://github.com/theory) 
* Reset Page's main output on server rebuilds [dc4226a8](https://github.com/gohugoio/hugo/commit/dc4226a8b27e03e31068fc945daab885d3819d04) [@bep](https://github.com/bep) [#4819](https://github.com/gohugoio/hugo/issues/4819)
* Create LICENSE rather than LICENSE.md in "new theme" [ed4a345e](https://github.com/gohugoio/hugo/commit/ed4a345efeaa19eef2c1c6360d22f75c24abc31a) [@anthonyfok](https://github.com/anthonyfok) [#4623](https://github.com/gohugoio/hugo/issues/4623)
* Create `_default/baseof.html` in "new theme" [9717ac7d](https://github.com/gohugoio/hugo/commit/9717ac7dce84d004afde4edb32ad81319c7dd8a7) [@anthonyfok](https://github.com/anthonyfok) [#3576](https://github.com/gohugoio/hugo/issues/3576)
* Add `tweet_simple`shortcode [07b96d16](https://github.com/gohugoio/hugo/commit/07b96d16e8679c40e289c9076ef4414ed6eb7f81) [@onedrawingperday](https://github.com/onedrawingperday) 
* Make "new theme" feedback more intuitive [692ec008](https://github.com/gohugoio/hugo/commit/692ec008726b570c9b30ac3391774cbb622730cb) [@anthonyfok](https://github.com/anthonyfok) 
* Move nextStepsText() to new_site.go [d3dd74fd](https://github.com/gohugoio/hugo/commit/d3dd74fd655c22f21e91e38edb1d377a1357e3be) [@anthonyfok](https://github.com/anthonyfok) 
* Add support for theme composition and inheritance [80230f26](https://github.com/gohugoio/hugo/commit/80230f26a3020ff33bac2bef01b2c0e314b89f86) [@bep](https://github.com/bep) [#4460](https://github.com/gohugoio/hugo/issues/4460)[#4450](https://github.com/gohugoio/hugo/issues/4450)
* Add vimeo_simple [8de53244](https://github.com/gohugoio/hugo/commit/8de53244799f0d2d0343056d348d810343cf7aa5) [@onedrawingperday](https://github.com/onedrawingperday) [#4749](https://github.com/gohugoio/hugo/issues/4749)
* Add a BlackFriday option for rel="noreferrer" on external links [20cbc2c7](https://github.com/gohugoio/hugo/commit/20cbc2c7856a9b07d45648d940276374db35e425) [@stefanneuhaus](https://github.com/stefanneuhaus) [#4722](https://github.com/gohugoio/hugo/issues/4722)
* Add a BlackFriday option for rel="nofollow" on external links [7a619264](https://github.com/gohugoio/hugo/commit/7a6192647a4b383cd539df2063388ea380371de6) [@stefanneuhaus](https://github.com/stefanneuhaus) [#4722](https://github.com/gohugoio/hugo/issues/4722)
* Update Chroma [b5b36e32](https://github.com/gohugoio/hugo/commit/b5b36e32008bc8ea779ae06bf249b537f6d5c336) [@bep](https://github.com/bep) [#4779](https://github.com/gohugoio/hugo/issues/4779)
* Enhance Page and Resource String() [4f0665f4](https://github.com/gohugoio/hugo/commit/4f0665f476e06e9707621c18f7422fdeb776e0d1) [@vassudanagunta](https://github.com/vassudanagunta) 
* Update theme documentation [c74b0f8f](https://github.com/gohugoio/hugo/commit/c74b0f8f9b30866e09efac8235cc5e0093ab3d50) [@bep](https://github.com/bep) [#4460](https://github.com/gohugoio/hugo/issues/4460)

## Fixes

* Make sure that `.Site.Taxonomies` is always set on rebuilds [6464981a](https://github.com/gohugoio/hugo/commit/6464981adb4d7d0f41e8e2c987342082982210a1) [@bep](https://github.com/bep) [#4838](https://github.com/gohugoio/hugo/issues/4838)
* Reset the "distinct error logger" on rebuilds [bf5f10fa](https://github.com/gohugoio/hugo/commit/bf5f10faa9fd445c4dd21839aa7d73cd2acbfb85) [@bep](https://github.com/bep) [#4818](https://github.com/gohugoio/hugo/issues/4818)
* Remove frameborder attr YT iframe + CSS fixes [ceaff7ca](https://github.com/gohugoio/hugo/commit/ceaff7cafc5357274e546984ae02a4cbdf305f81) [@onedrawingperday](https://github.com/onedrawingperday) 
* Fix vimeo_simple thumb scaling [b84389c5](https://github.com/gohugoio/hugo/commit/b84389c5e0e1ef15449b24d488bbbcbc41245c59) [@onedrawingperday](https://github.com/onedrawingperday) 
* Fix typo instagram_simple [d68367cb](https://github.com/gohugoio/hugo/commit/d68367cbe76cbc02adb5b778e8be98bed6319368) [@onedrawingperday](https://github.com/onedrawingperday) 
* Prevent isBaseTemplate() from matching "baseof" in dir [c3115292](https://github.com/gohugoio/hugo/commit/c3115292a7f2d2623cb45054a361e997ad9330c9) [@anthonyfok](https://github.com/anthonyfok) [#4809](https://github.com/gohugoio/hugo/issues/4809)





