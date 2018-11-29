---
title: path.Join
description: Join path elements into a single path.
godocref:
date: 2018-11-28
publishdate: 2018-11-28
lastmod: 2018-11-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [path, join]
signature: ["path.Join ELEMENT..."]
workson: []
hugoversion: "0.39"
relatedfuncs: [path.Split]
deprecated: false
---

`path.Join` joins path elements into a single path, adding a separating slash if necessary.
All empty strings are ignored.

**Note:** All path elements on Windows are converted to slash ('/') separators.

```
{{ path.Join "partial" "news.html" }} → "partial/news.html"
{{ path.Join "partial/" "news.html" }} → "partial/news.html"
{{ path.Join "foo/baz" "bar" }} → "foo/baz/bar"
```
