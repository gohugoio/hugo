---
title: Search tools
linkTitle: Search
description: See some of the open-source and commercial search options for your newly created Hugo website.
categories: []
keywords: []
weight: 30
---

A static website with a dynamic search function? Yes, Hugo provides an alternative to embeddable scripts from Google or other search engines for static websites. Hugo allows you to provide your visitors with a custom search function by indexing your content files directly.

## Open-source

[Pagefind][]
: A fully static search library that aims to perform well on large sites, while using as little of your users' bandwidth as possible.

[GitHub Gist for Hugo Workflow][]
: This gist contains a simple workflow to create a search index for your static website. It uses a simple Grunt script to index all your content files and [lunr.js][] to serve the search results.

[hugo-lunr][]
: A simple way to add site search to your static Hugo site using [lunr.js][]. Hugo-lunr will create an index file of any HTML and Markdown documents in your Hugo project.

[hugo-lunr-zh][]
: A bit like Hugo-lunr, but Hugo-lunr-zh can help you separate the Chinese keywords.

[GitHub Gist for Fuse.js integration][]
: This gist demonstrates how to leverage Hugo's existing build time processing to generate a searchable JSON index used by [Fuse.js][] on the client side. Although this gist uses Fuse.js for fuzzy matching, any client-side search tool capable of reading JSON indexes will work. Does not require npm, grunt, or other build-time tools except Hugo!

[hugo-search-index][]
: A library containing Gulp tasks and a prebuilt browser script that implements search. Gulp generates a search index from project Markdown files.

[hugofastsearch][]
: A usability and speed update to "GitHub Gist for Fuse.js integration" — global, keyboard-optimized search.

[JS & Fuse.js tutorial][]
: A simple client-side search solution, using FuseJS (does not require jQuery).

[hugo-zbsearch][]
: A Hugo Module providing client-side full-text search powered by [ZBSearch][]. Supports 27 languages with stemming, typo tolerance, field boosting, and works with any theme. Distributed as a Hugo Module with zero configuration needed beyond `hugo mod get`.

[INFINI Pizza for WebAssembly][]
: Pizza is a super-lightweight yet fully featured search engine written in Rust. You can quickly add offline search functionality to your Hugo website in just five minutes with only three lines of code. For a step-by-step guide on integrating it with Hugo, check out [this blog tutorial][].

[searchmysite.net][]
: An ad-free privacy-aware alternative to Google Programmable Search focussed on personal websites: add your site to searchmysite.net (if it isn't there already), and add a site search box to your Hugo site with one line of code. The free tier shows the search box on your site and results on searchmysite.net, and the paid tier shows results embedded in your site (and lets you reindex on demand etc.). All [open source][] with self-hosting an option. See [Website Search Tool][] and a [blog post with a Hugo-specific example][].

## Commercial

[Algolia DocSearch][]
: Algolia DocSearch is free for public technical documentation sites and easy to set up. For other use cases, [Algolia's Search API][] makes it easy to deliver a great search experience in your apps and websites. Algolia Search provides hosted full-text, numerical, faceted, and geolocalized search.

[Bonsai][]
: Bonsai is a fully-managed hosted Elasticsearch service that is fast, reliable, and simple to set up. Easily ingest your docs from Hugo into Elasticsearch following [this guide from the docs][].

[ExpertRec][]
: ExpertRec is a hosted search-as-a-service solution that is fast and scalable. Set-up and integration is extremely easy and takes only a few minutes. The search settings can be modified without coding using a dashboard.

[Algolia DocSearch]: https://docsearch.algolia.com/
[Algolia's Search API]: https://www.algolia.com
[Bonsai]: https://www.bonsai.io
[ExpertRec]: https://www.expertrec.com/
[Fuse.js]: https://fusejs.io/
[GitHub Gist for Fuse.js integration]: https://gist.github.com/eddiewebb/735feb48f50f0ddd65ae5606a1cb41ae
[GitHub Gist for Hugo Workflow]: https://gist.github.com/sebz/efddfc8fdcb6b480f567
[INFINI Pizza for WebAssembly]: https://github.com/infinilabs/pizza-docsearch
[JS & Fuse.js tutorial]: https://makewithhugo.com/add-search-to-a-hugo-site/
[Pagefind]: https://github.com/cloudcannon/pagefind
[hugo-lunr-zh]: https://www.npmjs.com/package/hugo-lunr-zh
[hugo-lunr]: https://www.npmjs.com/package/hugo-lunr
[hugo-search-index]: https://www.npmjs.com/package/hugo-search-index
[hugo-zbsearch]: https://github.com/serenasensini/hugo-zbsearch
[hugofastsearch]: https://gist.github.com/cmod/5410eae147e4318164258742dd053993
[lunr.js]: https://lunrjs.com/
[this blog tutorial]: https://dev.to/medcl/adding-search-functionality-to-a-hugo-static-site-based-on-infini-pizza-for-webassembly-4h5e
[this guide from the docs]: https://bonsai.io/docs/hugo
[ZBSearch]: https://github.com/micheleriva/zbsearch
[searchmysite.net]: https://searchmysite.net/
[open source]: https://github.com/searchmysite/searchmysite.net
[Website Search Tool]: https://searchmysite.net/pages/website-search-tool/
[blog post with a Hugo-specific example]: https://blog.searchmysite.net/posts/one-line-to-add-a-site-specific-search-to-your-site/
