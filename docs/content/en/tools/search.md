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

A static website with a dynamic search function? Yes, Hugo provides an alternative to embeddable scripts from Google or other search engines for static websites. Hugo allows you to provide your visitors with a custom search function by indexing your content files directly.

* [GitHub Gist for Hugo Workflow](https://gist.github.com/sebz/efddfc8fdcb6b480f567). This gist contains a simple workflow to create a search index for your static website. It uses a simple Grunt script to index all your content files and [lunr.js](https://lunrjs.com/) to serve the search results.
* [hugo-elasticsearch](https://www.npmjs.com/package/hugo-elasticsearch). Generate [Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html) indexes for Hugo static sites by parsing front matter. Hugo-Elasticsearch will generate a newline delimited JSON (NDJSON) file that can be bulk uploaded into Elasticsearch using any one of the available [clients](https://www.elastic.co/guide/en/elasticsearch/client/index.html).
* [hugo-lunr](https://www.npmjs.com/package/hugo-lunr). A simple way to add site search to your static Hugo site using [lunr.js](https://lunrjs.com/). Hugo-lunr will create an index file of any html and markdown documents in your Hugo project.
* [hugo-lunr-zh](https://www.npmjs.com/package/hugo-lunr-zh). A bit like Hugo-lunr, but Hugo-lunr-zh can help you separate the Chinese keywords.
* [Github Gist for Fuse.js integration](https://gist.github.com/eddiewebb/735feb48f50f0ddd65ae5606a1cb41ae). This gist demonstrates how to leverage Hugo's existing build time processing to generate a searchable JSON index used by [Fuse.js](https://fusejs.io/) on the client side. Although this gist uses Fuse.js for fuzzy matching, any client side search tool capable of reading JSON indexes will work. Does not require npm, grunt or other build-time tools except Hugo!
* [hugo-search-index](https://www.npmjs.com/package/hugo-search-index). A library containing Gulp tasks and a prebuilt browser script that implements search. Gulp generates a search index from project markdown files.

## Commercial Search Services

* [Algolia](https://www.algolia.com/)'s Search API makes it easy to deliver a great search experience in your apps and websites. Algolia Search provides hosted full-text, numerical, faceted, and geolocalized search.
* [Bonsai](https://www.bonsai.io) is a fully-managed hosted Elasticsearch service that is fast, reliable, and simple to set up. Easily ingest your docs from Hugo into Elasticsearch following [this guide from the docs](https://docs.bonsai.io/docs/hugo).
* [ExpertRec](https://www.expertrec.com/) is a hosted search-as-a-service solution that is fast and scalable. Set-up and integration is extremely easy and takes only a few minutes. The search settings can be modified without coding using a dashboard.
