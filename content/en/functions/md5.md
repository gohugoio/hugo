---
title: md5
description: hashes the given input and returns its MD5 checksum.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
signature: ["md5 INPUT"]
relatedfuncs: [sha]
---

```go-html-template
{{ md5 "Hello world, gophers!" }}
<!-- returns the string "b3029f756f98f79e7f1b7f1d1f0dd53b" -->
```

This can be useful if you want to use [Gravatar](https://en.gravatar.com/) for generating a unique avatar:

```html
<img src="https://www.gravatar.com/avatar/{{ md5 "your@email.com" }}?s=100&d=identicon">
```
