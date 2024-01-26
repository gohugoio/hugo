---
title: collections.Querify
description: Takes a set or slice of key-value pairs and returns a query string to be appended to URLs.
categories: []
keywords: []
action:
  aliases: [querify]
  related:
    - functions/go-template/urlquery.md
  returnType: string
  signatures:
    - collections.Querify VALUE [VALUE...]
    - collections.Querify COLLECTION
aliases: [/functions/querify]
---

`querify` takes a set or slice of key-value pairs and returns a [query string](https://en.wikipedia.org/wiki/Query_string) that can be appended to a URL.

The following examples create a link to a search results page on Google.

```go-html-template
<a href="https://www.google.com?{{ (querify "q" "test" "page" 3) | safeURL }}">Search</a>

{{ $qs := slice "q" "test" "page" 3 }}
<a href="https://www.google.com?{{ (querify $qs) | safeURL }}">Search</a>
```

Both of these examples render the following HTML:

```html
<a href="https://www.google.com?page=3&q=test">Search</a>
```
