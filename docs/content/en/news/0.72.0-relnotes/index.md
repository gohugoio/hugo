
---
date: 2020-05-31
title: URL rewrites in dev server
description: "Hugo 0.72.0 comes with dev server redirects and URL rewrites, Goldmark typography extension fixes, Scratch.Values."
categories: ["Releases"]
---

This is a rather small release, its probably main motivation being the fixes in Goldmark's [Typographer extension](https://github.com/gohugoio/hugo/commit/432885c499849efb29d3e50196f377fe0e908333).

This release also adds [redirect and URL rewrite support](https://gohugo.io/getting-started/configuration/#configure-server) to the development server, with mostly Netlify-compatible configuration syntax. This is especially useful if you're building a [SPA](https://en.wikipedia.org/wiki/Single-page_application) with client-side routing.


This release represents **13 contributions by 3 contributors** to the main Hugo code base.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **9 contributions by 6 contributors**. A special thanks to [@faraixyz](https://github.com/faraixyz), [@bep](https://github.com/bep), [@coliff](https://github.com/coliff), and [@Leon0824](https://github.com/Leon0824) for their work on the documentation site.


Hugo now has:

* 44383+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 437+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 327+ [themes](http://themes.gohugo.io/)

## Enhancements

### Other

* Add Scratch.Values [2919a6a5](https://github.com/gohugoio/hugo/commit/2919a6a503f7b369154d6eb787023a1fe58a9ad4) [@bep](https://github.com/bep) [#7335](https://github.com/gohugoio/hugo/issues/7335)
* Update Goldmark to improve Typographer [432885c4](https://github.com/gohugoio/hugo/commit/432885c499849efb29d3e50196f377fe0e908333) [@bep](https://github.com/bep) [#7289](https://github.com/gohugoio/hugo/issues/7289)
* Add redirect support to the server [6a3e8974](https://github.com/gohugoio/hugo/commit/6a3e89743ccad58097a6dd203a63448946a2304d) [@bep](https://github.com/bep) [#7323](https://github.com/gohugoio/hugo/issues/7323)

## Fixes

### Other

* Fix tag collector for nested table elements [c950c86b](https://github.com/gohugoio/hugo/commit/c950c86b4e5fb93f787ec78ca823bded9ef9fa3a) [@bep](https://github.com/bep) [#7318](https://github.com/gohugoio/hugo/issues/7318)
* Fix build error: my previous commits did not fix it [91520249](https://github.com/gohugoio/hugo/commit/915202494b140882d594e0542153531f6afada02) [@anthonyfok](https://github.com/anthonyfok) 


