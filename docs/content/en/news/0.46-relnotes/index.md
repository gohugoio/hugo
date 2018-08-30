
---
date: 2018-08-01
title: "The Summer of Hugo"
description: "Hugo 0.46: Full SCSS/SASS import inheritance support, rework of Hugo Pipes vs. multilingual, and more …"
categories: ["Releases"]
---

	**Hugo 0.46** is the closing credits to **The Summer of Hugo**. While most people have been relaxing on the beach, Hugo has been really busy: 

* **0.42, June 12, 2018**: [Theme Composition and Inheritance!](https://gohugo.io/news/0.42-relnotes/)
* **0.43, July 9, 2018**: [Hugo Pipes!](https://gohugo.io/news/0.43-relnotes/)

This was followed by some more technical follow-up releases. And today, when July has turned into August, we come with another one. It's not a big release. But with the big interest in **Hugo Pipes**, we felt that it was important to get this out there sooner rather than later.

There are two main items in this release:

1. We have added a custom SCSS/SASS import resolver that respects Hugo's project/themes filesystem hierarchy anywhere in `/assets`. Using the LibSass' resolver alone only made this work for the entry folder.
2. Resources fetched via `resources.Get` and similar are now language agnostic. The thought behind the original implementation was maximum flexibility with support for `assetDir` per language. In practice, this was a bad idea, as it made some CSS imports hard to get working in multilingual setups, and you got duplication of identical content for no good reason, with added processing time.

This release represents **12 contributions by 2 contributors** to the main Hugo code base.

A special thanks in this release goes to [@onedrawingperday](https://github.com/onedrawingperday) for his excellent work maintaining the fast-growing [Themes Site](https://themes.gohugo.io/).

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **5 contributions by 3 contributors**.

Hugo now has:

* 27596+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 442+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 245+ [themes](http://themes.gohugo.io/)

## Notes

* Resources fetched via `resources.Get` and similar are now considered global across languages. If you need, as an example, to create a CSS per language, you need to set the path yourself.

## Enhancements

* Add `templates.Exists` template function. This can be used to check if a template, e.g. a partial, exists in the project or any of the themes in use. [0ba19c57](https://github.com/gohugoio/hugo/commit/0ba19c57f180c33b41c64335ea1d1c89335d34c0) [@bep](https://github.com/bep) [#5010](https://github.com/gohugoio/hugo/issues/5010)
* Remove superflous loop [0afa2897](https://github.com/gohugoio/hugo/commit/0afa2897a0cf90f4348929ef432202efddc183a0) [@bep](https://github.com/bep) 
* Update Chroma [b5d13ca1](https://github.com/gohugoio/hugo/commit/b5d13ca16bf106c1bc29c2a5295cd231d1bf13fd) [@bep](https://github.com/bep) [#5019](https://github.com/gohugoio/hugo/issues/5019)
* Make resources fetched via `resources.Get` and similar language agnostic [6b02f5c0](https://github.com/gohugoio/hugo/commit/6b02f5c0f4e0ba1730aebc5a590a111548233bd5) [@bep](https://github.com/bep) [#5017](https://github.com/gohugoio/hugo/issues/5017)
* Improve SCSS project vs themes import resolution [f219ac09](https://github.com/gohugoio/hugo/commit/f219ac09f6b7e26d84599401512233d77c1bdb4c) [@bep](https://github.com/bep) [#5008](https://github.com/gohugoio/hugo/issues/5008)
* Improve _ prefix handling in SCSS imports [88e447c4](https://github.com/gohugoio/hugo/commit/88e447c449608523d87c517396bde31a62f392b6) [@bep](https://github.com/bep) [#5008](https://github.com/gohugoio/hugo/issues/5008)

## Fixes
* Fix file paths for uncached transformed images [b718d743](https://github.com/gohugoio/hugo/commit/b718d743b7a2eff3bea74ced57147825294a629f) [@RJacksonm1](https://github.com/RJacksonm1) [#5012](https://github.com/gohugoio/hugo/issues/5012)
* Fix image cache eviction for sites with subdir in baseURL [786f7230](https://github.com/gohugoio/hugo/commit/786f72302f65580ca8d1df2132a7756584539ea0) [@bep](https://github.com/bep) [#5006](https://github.com/gohugoio/hugo/issues/5006)
