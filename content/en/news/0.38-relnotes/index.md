
---
date: 2018-04-02
title: "Hugo 0.38: The Easter Egg Edition"
description: "Hugo 0.38: Date and slug from filenames, multiple content dirs, config from themes, language merge func …"
categories: ["Releases"]
---

Hugo `0.38` is an **Easter egg** filled with good stuff. We now support fetching **date and slug from the content filename**, making the move from Jekyll even easier. And you can now set `contentDir` per language with intelligent merging, and themes can now provide configuration ...  Also worth mentioning is several improvements in the [Chroma](https://github.com/alecthomas/chroma) highlighter, most notable support for Go templates.

We are working hard to get the documentation up-to-date with the new features, but you can also see them in action with the full source at [hugotest.bep.is](http://hugotest.bep.is/).

This release represents **39 contributions by 4 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@felicianotech](https://github.com/felicianotech), and [@paulcmal](https://github.com/paulcmal) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Also, a shoutout to [@regisphilibert](https://github.com/regisphilibert) for his work on the new [Code Toggle Shortcode](https://gohugo.io/getting-started/code-toggle/) on the Hugo docs site, which we will put to good use to improve all the configuration samples.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **55 contributions by 18 contributors**. A special thanks to [@kaushalmodi](https://github.com/kaushalmodi), [@bep](https://github.com/bep), [@xa0082249956](https://github.com/xa0082249956), and [@paulcmal](https://github.com/paulcmal) for their work on the documentation site.


Hugo now has:

* 24547+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 447+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 213+ [themes](http://themes.gohugo.io/)

## Notes

* Hugo now allows partial redefinition `outputs` in your site configuration. This is what most people would expect, but it is still a change in behaviour. For details, see [#4487](https://github.com/gohugoio/hugo/issues/4487)
* Before this release, Hugo flattened URLs of processed images in sub-folders. This worked fine but was not intentional. See [#4502](https://github.com/gohugoio/hugo/issues/4502).

## Enhancements

* Allow themes to define output formats, media types and params [e9c7b620](https://github.com/gohugoio/hugo/commit/e9c7b6205f94a7edac0e0df2cd18d1456cb26a06) [@bep](https://github.com/bep) [#4490](https://github.com/gohugoio/hugo/issues/4490)
* Allow partial redefinition of the `ouputs` config [f8dc47ee](https://github.com/gohugoio/hugo/commit/f8dc47eeffa847fd0b51e376da355e3d957848a6) [@bep](https://github.com/bep) [#4487](https://github.com/gohugoio/hugo/issues/4487)
* Add a way to merge pages by language [ffaec4ca](https://github.com/gohugoio/hugo/commit/ffaec4ca8c4c6fd05b195879ccd65acf2fd5a6ac) [@bep](https://github.com/bep) [#4463](https://github.com/gohugoio/hugo/issues/4463)
* Extract `date` and `slug` from filename [68bf1511](https://github.com/gohugoio/hugo/commit/68bf1511f2be39b6576d882d071196e477c72c9f) [@bep](https://github.com/bep) [#285](https://github.com/gohugoio/hugo/issues/285)[#3310](https://github.com/gohugoio/hugo/issues/3310)[#3762](https://github.com/gohugoio/hugo/issues/3762)[#4340](https://github.com/gohugoio/hugo/issues/4340)
* Add `Delete` method to delete key from `Scratch` [e46ab29b](https://github.com/gohugoio/hugo/commit/e46ab29bd24caa9e2cfa51f24ba15037750850d6) [@paulcmal](https://github.com/paulcmal) 
* Simplify Prev/Next [79dd7cb3](https://github.com/gohugoio/hugo/commit/79dd7cb31a941d7545df33b938ca3ed46593ddfd) [@bep](https://github.com/bep) 
* List Chroma lexers [2c54f1ad](https://github.com/gohugoio/hugo/commit/2c54f1ad48fe2a2f7504117d351d45abc89dcb1f) [@bep](https://github.com/bep) [#4554](https://github.com/gohugoio/hugo/issues/4554)
* Add support for a `contentDir` set per language [eb42774e](https://github.com/gohugoio/hugo/commit/eb42774e587816b1fbcafbcea59ed65df703882a) [@bep](https://github.com/bep) [#4523](https://github.com/gohugoio/hugo/issues/4523)[#4552](https://github.com/gohugoio/hugo/issues/4552)[#4553](https://github.com/gohugoio/hugo/issues/4553)
* Update Chroma [7a634898](https://github.com/gohugoio/hugo/commit/7a634898c359a6af0da52be17df07cae97c7937c) [@bep](https://github.com/bep) [#4549](https://github.com/gohugoio/hugo/issues/4549)
* Add `.Site.IsServer` [1823c053](https://github.com/gohugoio/hugo/commit/1823c053c8900cb6ee53b8e5c02939c7398e34dd) [@felicianotech](https://github.com/felicianotech) [#4478](https://github.com/gohugoio/hugo/issues/4478)
* Move to Ubuntu Trusty image [511d5d3b](https://github.com/gohugoio/hugo/commit/511d5d3b7681cb76822098f430ed6862232ca529) [@anthonyfok](https://github.com/anthonyfok) 
* Bump some deprecations [b6798ee8](https://github.com/gohugoio/hugo/commit/b6798ee8676c48f86b0bd8581ea244f4be4ef3fa) [@bep](https://github.com/bep) 
* Update Chroma to get `Go template support` [904a3d9d](https://github.com/gohugoio/hugo/commit/904a3d9ddf523d452d04d0b5814503e0ff17bd2e) [@bep](https://github.com/bep) [#4515](https://github.com/gohugoio/hugo/issues/4515)
* Recover from error in server [f0052b6d](https://github.com/gohugoio/hugo/commit/f0052b6d0f8e113a50aeb6cd7bd34555dbf34a00) [@bep](https://github.com/bep) [#4516](https://github.com/gohugoio/hugo/issues/4516)
* Spring test cleaning, take 2 [da880157](https://github.com/gohugoio/hugo/commit/da88015776645cc68b96e8b94030c95905df53ae) [@bep](https://github.com/bep) 
* Add docs for `lang.Merge` [70005364](https://github.com/gohugoio/hugo/commit/70005364a245ea3bc59c74192e1f4c56cb6879cf) [@bep](https://github.com/bep) 
* Remove archetype title/date warning [ac12d51e](https://github.com/gohugoio/hugo/commit/ac12d51e7ea3a0ffb7d8053a10b6bf6acf1235ae) [@bep](https://github.com/bep) [#4504](https://github.com/gohugoio/hugo/issues/4504)
* Add docs on the new front matter configuration [0dbf79c2](https://github.com/gohugoio/hugo/commit/0dbf79c2f8cd5b1a5c91c04a8d677f956b0b8fe8) [@bep](https://github.com/bep) [#4495](https://github.com/gohugoio/hugo/issues/4495)
* Refactor the GitInfo into the date handlers [ce6e4310](https://github.com/gohugoio/hugo/commit/ce6e4310febf5659392a41b543594382441f3681) [@bep](https://github.com/bep) [#4495](https://github.com/gohugoio/hugo/issues/4495)
* Do not print build total when `--quiet` is set [50a03a5a](https://github.com/gohugoio/hugo/commit/50a03a5acc7c200c795590c3f4b964fdc56085f2) [@bep](https://github.com/bep) [#4456](https://github.com/gohugoio/hugo/issues/4456)

## Fixes

* Fix freeze in invalid front matter error case [93e24a03](https://github.com/gohugoio/hugo/commit/93e24a03ce98d3212a2d49ad04739141229d0809) [@bep](https://github.com/bep) [#4526](https://github.com/gohugoio/hugo/issues/4526)
* Fix path duplication/flattening in processed images [3fbc7553](https://github.com/gohugoio/hugo/commit/3fbc75534d1acda2be1c597aa77c919d3a02659d) [@bep](https://github.com/bep) [#4502](https://github.com/gohugoio/hugo/issues/4502)[#4501](https://github.com/gohugoio/hugo/issues/4501)
* Fix SVG and similar resource handling [ba94abbf](https://github.com/gohugoio/hugo/commit/ba94abbf5dd90f989242af8a7027d67a572a6128) [@bep](https://github.com/bep) [#4455](https://github.com/gohugoio/hugo/issues/4455)




