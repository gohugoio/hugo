---
title: Search for your Hugo Website
linktitle: Search
description: See some of the open-source and commercial search options for your newly created Hugo website.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-26
categories: [developer tools]
keywords: [search,tools]
menu:
  docs:
    parent: "tools"
    weight: 60
weight: 60
sections_weight: 60
draft: false
aliases: []
toc: true
---

A static website with a dynamic search function? Yes. As alternatives to embeddable scripts from Google or other search engines, you can provide your visitors a custom search by indexing your content files directly.

* [Hugoidx](https://github.com/blevesearch/hugoidx) is an experimental application to create a search index. It's built on top of [Bleve](http://www.blevesearch.com/).
* [GitHub Gist for Hugo Workflow](https://gist.github.com/sebz/efddfc8fdcb6b480f567). This gist contains a simple workflow to create a search index for your static website. It uses a simple Grunt script to index all your content files and [lunr.js](http://lunrjs.com/) to serve the search results.
* [hugo-lunr](https://www.npmjs.com/package/hugo-lunr). A simple way to add site search to your static Hugo site using [lunr.js](http://lunrjs.com/). Hugo-lunr will create an index file of any html and markdown documents in your Hugo project.
* [hugo-lunr-zh](https://www.npmjs.com/package/hugo-lunr-zh). A bit like Hugo-lunr, but Hugo-lunr-zh can help you seperate the Chinese keywords.

## Commercial Search Services

* [Algolia](https://www.algolia.com/)'s Search API makes it easy to deliver a great search experience in your apps and websites. Algolia Search provides hosted full-text, numerical, faceted, and geolocalized search.
