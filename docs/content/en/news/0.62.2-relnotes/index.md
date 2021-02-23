
---
date: 2020-01-05
title: "Hugo 0.62.2: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.62.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

The main driving force behind these patch releases in the new year has been getting a good story with [portable links](https://github.com/bep/portable-hugo-links/) between GitHub and Hugo, using the new render hooks introduced in [Hugo 0.62.0](https://gohugo.io/news/0.62.0-relnotes/). And all was mostly well until a Hugo user asked about anchor links. Which is, when you look into it, a slightly sad Markdown story. They have been [talking about anchors in Markdown](https://talk.commonmark.org/t/anchors-in-markdown/247) over at the CommonMark forum for the last six years, but it has come to nothing. The current situation is that some Markdown engines support the `{#anchorName}` attribute syntax (Hugo's Goldmark does, GitHub does not) and for those that support auto-generation of IDs, the implementation often varies. And this makes for poor portability.

To improve this, Hugo has now reverse-engineered GitHub's implementation and made that the default strategy for generation or header IDs. We understand that this isn't everyone's cup of tea, so you can [configure the behaviour](https://gohugo.io/getting-started/configuration-markup#goldmark) to be one of `github`, `github-ascii` (some client-libraries have Unicode issues) or `blackfriday` (which will match how it behaved before Hugo 0.60).

* hugolib: Fix relative .Page.GetPage from bundle [196a9df5](https://github.com/gohugoio/hugo/commit/196a9df585c4744e3280f37c1c24e469fce14b8c) [@bep](https://github.com/bep) [#6705](https://github.com/gohugoio/hugo/issues/6705)
* markup/goldmark: Adjust auto ID space handling [9b6e6146](https://github.com/gohugoio/hugo/commit/9b6e61464b09ffe3423fb8d7c72bddb7a9ed5b98) [@bep](https://github.com/bep) [#6710](https://github.com/gohugoio/hugo/issues/6710)
* docs: Document the new autoHeadingIDType setting [d62ede8e](https://github.com/gohugoio/hugo/commit/d62ede8e9e5883e7ebb023e49b82f07b45edc1c7) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)[#6616](https://github.com/gohugoio/hugo/issues/6616)
* docs: Regenerate docshelper [81b7e48a](https://github.com/gohugoio/hugo/commit/81b7e48a55092203aeee8785799e6fed3928760e) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)[#6616](https://github.com/gohugoio/hugo/issues/6616)
* markup/goldmark: Add an optional Blackfriday auto ID strategy [16e7c112](https://github.com/gohugoio/hugo/commit/16e7c1120346bd853cf6510ffac8e94824bf2c7f) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)
* markup/goldmark: Make the autoID type config a string [8f071fc1](https://github.com/gohugoio/hugo/commit/8f071fc159ce9a0fc0ea14a73bde8f299bedd109) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)
* markup/goldmark: Simplify code [5ee1f087](https://github.com/gohugoio/hugo/commit/5ee1f0876f3ec8b79d6305298185dc821ead2d28) [@bep](https://github.com/bep) 
* markup/goldmark: Make auto IDs GitHub compatible [a82d2700](https://github.com/gohugoio/hugo/commit/a82d2700fcc772aada15d65b8f76913ca23f7404) [@bep](https://github.com/bep) [#6616](https://github.com/gohugoio/hugo/issues/6616)



