
---
date: 2018-07-09
title: "And Now: Hugo Pipes!"
description: "Hugo 0.43 adds a powerful and simple to use assets pipeline with SASS/SCSS and much, much more …"
categories: ["Releases"]
---

Hugo `0.43` adds a powerful and very simple to use **Assets Pipeline** with **SASS and SCSS** with source map support, **PostCSS** and **minification** and **fingerprinting** and **Subresource Integrity** and ... much more. Oh, did we mention that you can now do **ad-hoc image processing** and execute text resources as Go templates?

An example pipeline:

```go-html-template
{{ $styles := resources.Get "scss/main.scss" | toCSS | postCSS | minify | fingerprint }}
<link rel="stylesheet" href="{{ $styles.Permalink }}" integrity="{{ $styles.Data.Integrity }}" media="screen">
```

To me, the above is beautiful in its speed and simplicity. It could be printed on a t-shirt. I wrote in the [Hugo Birthday Post](https://gohugo.io/news/lets-celebrate-hugos-5th-birthday/) some days ago about the value of a single binary with native and fast implementations. I should have added _simplicity_ as a keyword. There seem to be a misconception that all of this needs to be hard and painful.

New functions to create `Resource` objects:

* `resources.Get`
* `resources.FromString`: Create a Resource from a string.

New `Resource` transformation funcs:

* `resources.ToCSS`: Compile `SCSS` or `SASS` into `CSS`.
* `resources.PostCSS`: Process your CSS with PostCSS. Config file support (project or theme or passed as an option).
* `resources.Minify`: Currently supports `css`, `js`, `json`, `html`, `svg`, `xml`.
* `resources.Fingerprint`: Creates a fingerprinted version of the given Resource with Subresource Integrity.
* `resources.Concat`: Concatenates a list of Resource objects. Think of this as a poor man's bundler.
* `resources.ExecuteAsTemplate`: Parses and executes the given Resource and data context (e.g. .Site) as a Go template.


I, [@bep](https://github.com/bep), implemented this in [dea71670](https://github.com/gohugoio/hugo/commit/dea71670c059ab4d5a42bd22503f18c087dd22d4). We will work hard to get the documentation up to date, but follow the links above for details, and also see this [demo project](https://github.com/bep/hugo-sass-test).


This release represents **35 contributions by 7 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@openscript](https://github.com/openscript), and [@caarlos0](https://github.com/caarlos0) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **11 contributions by 5 contributors**. A special thanks to [@bep](https://github.com/bep), [@danrl](https://github.com/danrl), [@regisphilibert](https://github.com/regisphilibert), and [@digitalcraftsman](https://github.com/digitalcraftsman) for their work on the documentation site.

Hugo now has:

* 26968+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 443+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 238+ [themes](http://themes.gohugo.io/)

## Notes

* Replace deprecated {Get,}ByPrefix with {Get,}Match [42ed6025](https://github.com/gohugoio/hugo/commit/42ed602580a672e420e1d860384e812f4871ff67) [@anthonyfok](https://github.com/anthonyfok)
* Hugo is now released with two binary version: One with and one without SCSS/SASS support. At the time of writing, this is only available in the binaries on the GitHub release page. Brew, Snap builds etc. will come. But note that you **only need the extended version if you want to edit SCSS**. For your CI server, or if you don't use SCSS, you will most likely want the non-extended version.

## Enhancements

### Templates

* Return en empty slice in `after` instead of error [f8212d20](https://github.com/gohugoio/hugo/commit/f8212d20009c4b5cc6e1ec733d09531eb6525d9f) [@bep](https://github.com/bep) [#4894](https://github.com/gohugoio/hugo/issues/4894)
* Update internal pagination template to support Bootstrap 4 [ca1e46ef](https://github.com/gohugoio/hugo/commit/ca1e46efb94e3f3d2c8482cb9434d2f38ffd2683) [@bep](https://github.com/bep) [#4881](https://github.com/gohugoio/hugo/issues/4881)
* Support text/template/parse API change in go1.11 [9f27091e](https://github.com/gohugoio/hugo/commit/9f27091e1067875e2577c331acc60adaef5bb234) [@anthonyfok](https://github.com/anthonyfok) [#4784](https://github.com/gohugoio/hugo/issues/4784)

### Core

* Allow forward slash in shortcode names [de37455e](https://github.com/gohugoio/hugo/commit/de37455ec73cffd039b44e8f6c62d2884b1d6bbd) [@bep](https://github.com/bep) [#4886](https://github.com/gohugoio/hugo/issues/4886)
* Reset the global pages cache on server rebuilds [128f14ef](https://github.com/gohugoio/hugo/commit/128f14efad90886ffef37c01ac1e20436a732f97) [@bep](https://github.com/bep) [#4845](https://github.com/gohugoio/hugo/issues/4845)

### Other

* Bump CircleCI image [e3df6478](https://github.com/gohugoio/hugo/commit/e3df6478f09a7a5fed96aced791fa94fd2c35d1a) [@bep](https://github.com/bep)
* Add Goreleaser extended config [626afc98](https://github.com/gohugoio/hugo/commit/626afc98254421f5a5edc97c541b10bd81d5bbbb) [@bep](https://github.com/bep) [#4908](https://github.com/gohugoio/hugo/issues/4908)
* Build both hugo and hugo.extended for 0.43 [e1027c58](https://github.com/gohugoio/hugo/commit/e1027c5846b48c4ad450f6cc27e2654c9e0dae39) [@anthonyfok](https://github.com/anthonyfok) [#4908](https://github.com/gohugoio/hugo/issues/4908)
* Add temporary build script [bfc3488b](https://github.com/gohugoio/hugo/commit/bfc3488b8e8b3dc1ffc6a339ee2dac8dcbdb55a9) [@bep](https://github.com/bep)
* Add "extended" to "hugo version" [ce84b524](https://github.com/gohugoio/hugo/commit/ce84b524f4e94299b5b66afe7ce1a9bd4a9959fc) [@anthonyfok](https://github.com/anthonyfok) [#4913](https://github.com/gohugoio/hugo/issues/4913)
* Add a newScratch template func [2b8d907a](https://github.com/gohugoio/hugo/commit/2b8d907ab731627f4e2a30442cd729064516c8bb) [@bep](https://github.com/bep) [#4685](https://github.com/gohugoio/hugo/issues/4685)
* Add Hugo Piper with SCSS support and much more [dea71670](https://github.com/gohugoio/hugo/commit/dea71670c059ab4d5a42bd22503f18c087dd22d4) [@bep](https://github.com/bep) [#4381](https://github.com/gohugoio/hugo/issues/4381)[#4903](https://github.com/gohugoio/hugo/issues/4903)[#4858](https://github.com/gohugoio/hugo/issues/4858)
* Consider root and current section's content type if set in front matter [c790029e](https://github.com/gohugoio/hugo/commit/c790029e1dbb0b66af18d05764bd6045deb2e180) [@bep](https://github.com/bep) [#4891](https://github.com/gohugoio/hugo/issues/4891)
* Update docker image [554553c0](https://github.com/gohugoio/hugo/commit/554553c09c7657d28681e1fa0638806a452737a0) [@bep](https://github.com/bep)
* Merge branch 'release-0.42.2' [282f6035](https://github.com/gohugoio/hugo/commit/282f6035e7c36f8550d91033e3a66718468c6c8b) [@bep](https://github.com/bep)
* Release 0.42.2 [1637d12e](https://github.com/gohugoio/hugo/commit/1637d12e3762fc1ebab4cd675f75afaf25f59cdb) [@bep](https://github.com/bep)
* Update GoReleaser config [1f0c4e1f](https://github.com/gohugoio/hugo/commit/1f0c4e1fb347bb233f3312c424fbf5a013c03604) [@caarlos0](https://github.com/caarlos0)
* Create missing head.html partial on new theme generation [fd71fa89](https://github.com/gohugoio/hugo/commit/fd71fa89bd6c197402582c87b2b76d4b96d562bf) [@openscript](https://github.com/openscript)
* Add html doctype to baseof.html template for new themes [b5a3aa70](https://github.com/gohugoio/hugo/commit/b5a3aa7082135d0a573f4fbb00f798e26b67b902) [@openscript](https://github.com/openscript)
* Adds .gitattributes to force Go files to LF [6a2968fd](https://github.com/gohugoio/hugo/commit/6a2968fd5c0116d93de0f379ac615e9076821899) [@neurocline](https://github.com/neurocline)
* Update to Go 1.9.7 and Go 1.10.3 [23d5fc82](https://github.com/gohugoio/hugo/commit/23d5fc82ee01d56440d0991c899acd31e9b63e27) [@anthonyfok](https://github.com/anthonyfok)
* Update Dockerfile to a multi-stage build [8531ec7c](https://github.com/gohugoio/hugo/commit/8531ec7ca36fd35a57fba06bbb06a65c94dfd3ed) [@skoblenick](https://github.com/skoblenick) [#4154](https://github.com/gohugoio/hugo/issues/4154)
* Release 0.42.1 [d67e843c](https://github.com/gohugoio/hugo/commit/d67e843c1212e1f53933556b5f946c8541188d9a) [@bep](https://github.com/bep)
* Do not fail server build when /static is missing [34ee27a7](https://github.com/gohugoio/hugo/commit/34ee27a78b9e2b5f475d44253ae234067b76cc6e) [@bep](https://github.com/bep) [#4846](https://github.com/gohugoio/hugo/issues/4846)

## Fixes

* Do not create paginator pages for the other output formats [43338c3a](https://github.com/gohugoio/hugo/commit/43338c3a99769eb7d0df0c12559b8b3d42b67dba) [@bep](https://github.com/bep) [#4890](https://github.com/gohugoio/hugo/issues/4890)
* Fix the shortcodes/partials vs base template detection [a5d0a57e](https://github.com/gohugoio/hugo/commit/a5d0a57e6bdab583134a68c035aac9b3007f006a) [@bep](https://github.com/bep) [#4897](https://github.com/gohugoio/hugo/issues/4897)
* nfpm replacements [e1a052ec](https://github.com/gohugoio/hugo/commit/e1a052ecb823c688406a8af97dfaaf52a75231da) [@caarlos0](https://github.com/caarlos0)
* Fix typos [3cea2932](https://github.com/gohugoio/hugo/commit/3cea2932e17a08ebc19cd05f3079d9379bc8fba5) [@idealhack](https://github.com/idealhack)
