---
title: urls.Parse
description: Parse parses a given URL, which may be relative or absolute, into a URL structure.
date: 2017-09-25
publishdate: 2017-09-25
lastmod: 2017-09-25
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [urls]
signature: ["urls.Parse URL"]
workson: []
hugoversion:
deprecated: false
aliases: []
---

`urls.Parse` takes a url as input


```go-html-template
{{ $url := urls.Parse "http://www.gohugo.io" }}
```

and returns a [URL](https://godoc.org/net/url#URL) structure. The struct fields are accessed via the `.` notation:

```go-html-template
{{ $url.Scheme }} → "http"
{{ $url.Host }} → "www.gohugo.io"
```
