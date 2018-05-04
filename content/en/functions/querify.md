---
title: querify
linktitle: querify
description: Takes a set of key-value pairs and returns a query string to be appended to URLs.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [urls]
godocref:
signature: ["querify KEY VALUE [KEY VALUE]..."]
hugoversion:
deprecated: false
workson: []
relatedfuncs: []
aliases: []
---

`querify` takes a set of key-value pairs and returns a [query string](https://en.wikipedia.org/wiki/Query_string) that can be appended to a URL. E.g.

The following example creates a link to a search results page on Google.

```
<a href="https://www.google.com?{{ (querify "q" "test" "page" 3) | safeURL }}">Search</a>
```

This example renders the following HTML:

```
<a href="https://www.google.com?page=3&q=test">Search</a>
```
