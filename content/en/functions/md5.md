---
title: md5
linktitle: md5
description: hashes the given input and returns its MD5 checksum.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
signature: ["md5 INPUT"]
workson: []
hugoversion:
relatedfuncs: [sha]
deprecated: false
aliases: []
---



```
{{ md5 "Hello world, gophers!" }}
<!-- returns the string "b3029f756f98f79e7f1b7f1d1f0dd53b" -->
```

This can be useful if you want to use [Gravatar](https://en.gravatar.com/) for generating a unique avatar:

```
<img src="https://www.gravatar.com/avatar/{{ md5 "your@email.com" }}?s=100&d=identicon">
```
