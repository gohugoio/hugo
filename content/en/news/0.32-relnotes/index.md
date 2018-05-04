
---
date: 2017-12-31
title: "Hugo 0.32: Page Bundles and Image Processing!"
description: "Images and other resources with page-relative links, resize, scale and crop images, and much more."
categories: ["Releases"]
images:
- images/blog/hugo-32-poster.png
---

	Hugo `0.32` features **Page Bundles and Image Processing** by [@bep](https://github.com/bep), which is very cool and useful on so many levels. Read about it in more detail in the [Hugo documentation](https://gohugo.io/about/new-in-032/), but some of the highlights include:

* Automatic bundling of a content page with its resources. Resources can be anything: Images, `JSON` files ... and also other content pages.
* A `Resource` will have its `RelPermalink` and `Permalink` relative to the "owning page". This makes the complete article with both text and images portable (just send a ZIP file with a folder to your editor), and it can be previewed directly on GitHub.
* Powerful and simple to use image processing with the new `.Resize`, `.Fill`, and `.Fit` methods on the new `Image` resource.
* Full support for symbolic links inside `/content`, both for regular files and directories.

The built-in benchmarks in Hugo show that this is also the [fastest and most memory effective](https://gist.github.com/bep/2a9bbd221de2da5d39c8b32085c658f7) Hugo version to date. But note that the build time total reported in the console is now adjusted to be the *real total*, including the copy of static files. So, if it reports more milliseconds, it is still most likely faster ...

This release represents **30 contributions by 7 contributors** to the main Hugo code base.

[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@betaveros](https://github.com/betaveros), [@chaseadamsio](https://github.com/chaseadamsio), and [@kropp](https://github.com/kropp). And as always big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **17 contributions by 7 contributors**. A special thanks to [@bep](https://github.com/bep), [@felicianotech](https://github.com/felicianotech), [@maiki](https://github.com/maiki), and [@carlchengli](https://github.com/carlchengli) for their work on the documentation site.

Hugo now has:

* 22061+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 454+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 193+ [themes](http://themes.gohugo.io/)

Today is **New Year's Eve.** It is the last day of 2017, a year that have seen a **string of pearls of Hugo releases**, making Hugo _the_ top choice for website development:

* 0.32, December 2017: **Page Bundles and Image Processing** edition.
* 0.31, November 2017: The Language **Multihost Edition!** with one `baseURL` per language.
* 0.30, October 2017: The Race Car Edition with the **Fast Render Mode**.
* 0.29, September 2017: Added **Template Metrics**.
*  0.28, September 2017:  **Blistering fast and native syntax highlighting** from [Chroma](https://github.com/alecthomas/chroma).
* 0.27, September 2017: Fast and flexible **Related Content.**
*  0.26, August 2017: The **Language Style Edition**  with AP Style or Chicago Style Title Case and « French Guillemets ».
* 0.25, July 2017: The **Kinder Surprise** edition added, among other cool things, `hugo server --navigateToChanged` which navigates to the content page you start editing.
* 0.24, June 2017: Was **The Revival of the Archetypes!** Now archetype files, i.e. the content file templates, can include template syntax with all of Hugo's functions and variables.
* 0.23, June 2017: Hugo moved to it's own GitHub organization, **gohugoio**.
* 0.22, June 2017: Added **nested sections**, a long sought after feature.
* 0.21, May 2017: Full support for shortcodes per output format (think **AMP**).
* 0.20, April 2017: Was all about **Custom Output Formats**.
* 0.19, February 2017: Native Emacs Org-mode content support and lots of internal upgrades.

## Notes

* The build total in the console is now the ... total (i.e. it now includes both the copy of the static files and the Hugo build). So if your Hugo site seems to build slightly slower, it is in reality probably slightly faster than before this release.
* Images and other static resources in folders with "_index.md" will have its `RelPermalink` relative to its page.
* Images and other static resources in or below "index.md" folders will have its `RelPermalink` relative to its page (respecting permalink settings etc.)
* Content pages in or below "index.md" will not get their own `URL`, but will be part of the `.Resources` collection of its page.
* `.Site.Files` is deprecated.
* Hugo no longer minfies CSS files inside `/content`. This was an undocumented "proof of concept feature". We may revisit the "assets handling" in a future release.	
* `Page.GetParam`does not lowercase your result anymore. If you really want to lowercase your params, do it with `.GetParam "myparam" | lower` or similar.

Previously deprecated that will now `ERROR`:

* `disable404`: Use `disableKinds=["404"]`
* `disableRSS`:  Use `disableKinds=["RSS"]`
* `disableSitemap`:  Use `disableKinds=["sitemap"]`
* `disableRobotsTXT`: Use `disableKinds=["robotsTXT"]`

## Enhancements

* Add `.Title` and `.Page` to `MenuEntry` [9df3736f](https://github.com/gohugoio/hugo/commit/9df3736fec164c51d819797416dc263f2869be77) [@rmetzler](https://github.com/rmetzler) [#2784](https://github.com/gohugoio/hugo/issues/2784)
* Add `Pandoc` support [e69da7a4](https://github.com/gohugoio/hugo/commit/e69da7a4cb725987f153707bf2fc59c135007e2a) [@betaveros](https://github.com/betaveros) [#234](https://github.com/gohugoio/hugo/issues/234)
* Implement Page bundling and image handling [3cdf19e9](https://github.com/gohugoio/hugo/commit/3cdf19e9b7e46c57a9bb43ff02199177feb55768) [@bep](https://github.com/bep) [#3651](https://github.com/gohugoio/hugo/issues/3651)[#3158](https://github.com/gohugoio/hugo/issues/3158)[#1014](https://github.com/gohugoio/hugo/issues/1014)[#2021](https://github.com/gohugoio/hugo/issues/2021)[#1240](https://github.com/gohugoio/hugo/issues/1240)[#3757](https://github.com/gohugoio/hugo/issues/3757)
* Make `chomp` return the type it receives [22cd89ad](https://github.com/gohugoio/hugo/commit/22cd89adc4792a3b55389d38acd4acfae3786775) [@kropp](https://github.com/kropp) [#2187](https://github.com/gohugoio/hugo/issues/2187) 
* Reuse the `BlackFriday` config instance when possible [db4b7a5c](https://github.com/gohugoio/hugo/commit/db4b7a5c6742c75f9cd9627d3b054d3a72802ec8) [@bep](https://github.com/bep) 
* Remove the goroutines from the shortcode lexer [24369410](https://github.com/gohugoio/hugo/commit/243694102a60da2fb1050020f68384539f9f9ef5) [@bep](https://github.com/bep) 
* Improve site benchmarks [051fa343](https://github.com/gohugoio/hugo/commit/051fa343d06d6c070df742f7cbd125432fcab665) [@bep](https://github.com/bep) 
* Update `Chroma` to `v0.2.0` [79892101](https://github.com/gohugoio/hugo/commit/7989210120dbde78da3741e2ef01b13f4aa78692) [@bep](https://github.com/bep) [#4087](https://github.com/gohugoio/hugo/issues/4087)
* Update `goorgeous` to `v1.1.0` [7f2ae3ef](https://github.com/gohugoio/hugo/commit/7f2ae3ef39f27a9bd26ddb9258b073a840faf491) [@chaseadamsio](https://github.com/chaseadamsio) 
* Add test for homepage content for all rendering engines [407c2402](https://github.com/gohugoio/hugo/commit/407c24020ef2db90cf33fd07e7522b2257013722) [@bep](https://github.com/bep) [#4166](https://github.com/gohugoio/hugo/issues/4166)
* Add output formats definition to benchmarks [a2d81ce9](https://github.com/gohugoio/hugo/commit/a2d81ce983d45b5742c93bd472503c88286f099a) [@bep](https://github.com/bep) 

## Fixes

### Templates

* Do not unescape input to `highlight` [c067f345](https://github.com/gohugoio/hugo/commit/c067f34558b82455b63b9ce8f5983b4b4849c7cf) [@bep](https://github.com/bep) [#4179](https://github.com/gohugoio/hugo/issues/4179)
* Properly close image file in `imageConfig` [6d79beb5](https://github.com/gohugoio/hugo/commit/6d79beb5f67dbb54d7714c3195addf9d8e3924e8) [@bep](https://github.com/bep) 
 * Fix  `opengraph` video range template [23f69efb](https://github.com/gohugoio/hugo/commit/23f69efb3914946b39ce673fcc0f2e3a9ed9d878) [@drlogout](https://github.com/drlogout) [#4136](https://github.com/gohugoio/hugo/issues/4136)
* Fix `humanize` for multi-byte runes [e7652180](https://github.com/gohugoio/hugo/commit/e7652180a13ce149041c48a1c2754c471df569c8) [@bep](https://github.com/bep) [#4133](https://github.com/gohugoio/hugo/issues/4133)

### Other

* Fix broken live reload without a server port. [25114986](https://github.com/gohugoio/hugo/commit/25114986086e5877a0b4108d8cf5e4e95f377241) [@sainaen](https://github.com/sainaen) [#4141](https://github.com/gohugoio/hugo/issues/4141)
* Make sure all language homes are always re-rendered in fast render mode [72903be5](https://github.com/gohugoio/hugo/commit/72903be587e9c4e3644f60b11e26238ec03da2db) [@bep](https://github.com/bep) [#4125](https://github.com/gohugoio/hugo/issues/4125)
* Do not `tolower` result from Page.GetParam [1c114d53](https://github.com/gohugoio/hugo/commit/1c114d539b0755724443fe28c90b12fe2a19085a) [@bep](https://github.com/bep) [#4187](https://github.com/gohugoio/hugo/issues/4187)
