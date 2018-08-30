
---
date: 2018-07-13
title: "Hugo 0.44: Friday the 13th Edition"
description: "A sequel to the very popular Hugo Pipes Edition; bug-fixes and enhancements …"
categories: ["Releases"]
---

	
Hugo `0.44` is the follow-up release, or **The Sequel**, of the very well received `0.43` only days ago. That release added **Hugo Pipes**, with SCSS/SASS support, assets bundling and minification, ad-hoc image processing and much more.

This is mostly a bug-fix release, but it also includes several important improvements.

Many complained that their SVG images vanished when browsed from the `hugo server`. With **Hugo Pipes** MIME types suddenly got really important, but Hugo's use of `Suffix` was ambiguous. This became visible when we redefined the `image/svg+xml` to work with **Hugo Pipes**. We have now added a `Suffixes` field on the MIME type definition in Hugo, which is a list of one or more filename suffixes the MIME type is identified with. If you need to add a custom MIME type definition, this means that you also need to specify the full MIME type as the key, e.g. `image/svg+xml`.

Hugo now has:

* 27120+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 443+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 239+ [themes](http://themes.gohugo.io/)

## Notes
* `MediaType.Suffix` is deprecated and replaced with a plural version,  `MediaType.Suffixes`, with a more specific definition. You will get a detailed WARNING in the console if you need to do anything.

## Enhancements
* Allow multiple file suffixes per media type [b874a1ba](https://github.com/gohugoio/hugo/commit/b874a1ba7ab8394dc741c8c70303a30a35b63e43) [@bep](https://github.com/bep) [#4920](https://github.com/gohugoio/hugo/issues/4920)
* Clean up the in-memory Resource reader usage [47d38628](https://github.com/gohugoio/hugo/commit/47d38628ec0f4e72ff17661f13456b2a1511fe13) [@bep](https://github.com/bep) [#4936](https://github.com/gohugoio/hugo/issues/4936)
* Move opening of the transformed resources after cache check [0024dcfe](https://github.com/gohugoio/hugo/commit/0024dcfe3e016c67046de06d1dac5e7f5235f9e1) [@bep](https://github.com/bep) 
* Improve type support in `resources.Concat` [306573de](https://github.com/gohugoio/hugo/commit/306573def0e20ec16ee5c447981cc09ed8bb7ec7) [@bep](https://github.com/bep) [#4934](https://github.com/gohugoio/hugo/issues/4934)
* Flush `partialCached` cache on rebuilds [6b6dcb44](https://github.com/gohugoio/hugo/commit/6b6dcb44a014699c289bf32fe57d4c4216777be0) [@bep](https://github.com/bep) [#4931](https://github.com/gohugoio/hugo/issues/4931)
* Include the transformation step in the error message [d96f2a46](https://github.com/gohugoio/hugo/commit/d96f2a460f58e91d8f6253a489d4879acfec6916) [@bep](https://github.com/bep) [#4924](https://github.com/gohugoio/hugo/issues/4924)
* Exclude *.svg from CRLF/LF conversion [9c1e8208](https://github.com/gohugoio/hugo/commit/9c1e82085eb07d5b4dcdacbe82d5bafd26e08631) [@anthonyfok](https://github.com/anthonyfok) 

## Fixes

* Fix `resources.Concat` for transformed resources [beec1fc9](https://github.com/gohugoio/hugo/commit/beec1fc98e5d37bba742d6bc2a0ff7c344b469f8) [@bep](https://github.com/bep) [#4936](https://github.com/gohugoio/hugo/issues/4936)
* Fix static filesystem for themed multihost sites [80c8f3b8](https://github.com/gohugoio/hugo/commit/80c8f3b81a9849080e64bf877288ede28d960d3f) [@bep](https://github.com/bep) [#4929](https://github.com/gohugoio/hugo/issues/4929)
* Set permission of embedded templates to 0644 [2b73e89d](https://github.com/gohugoio/hugo/commit/2b73e89d6d2822e86360a6c92c87f539677c119b) [@anthonyfok](https://github.com/anthonyfok) 

