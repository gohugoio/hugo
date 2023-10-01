---
title: crypto.MD5
linkTitle: md5
description: hashes the given input and returns its MD5 checksum.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [md5]
  returnType: string
  signatures: [crypto.MD5 INPUT]
relatedFunctions:
  - crypto.FNV32a
  - crypto.HMAC
  - crypto.MD5
  - crypto.SHA1
  - crypto.SHA256
aliases: [/functions/md5]
---

```go-html-template
{{ md5 "Hello world" }} → 3e25960a79dbc69b674cd4ec67a72c62

```

This can be useful if you want to use [Gravatar](https://en.gravatar.com/) for generating a unique avatar:

```html
<img src="https://www.gravatar.com/avatar/{{ md5 "your@email.com" }}?s=100&d=identicon">
```
