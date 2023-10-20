---
title: crypto.SHA1
linkTitle: sha1
description: Hashes the given input and returns its SHA1 checksum.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [sha1]
  returnType: string
  signatures: [crypto.SHA1 INPUT]
relatedFunctions:
  - crypto.FNV32a
  - crypto.HMAC
  - crypto.MD5
  - crypto.SHA1
  - crypto.SHA256
aliases: [/functions/sha,/functions/sha1]
---

```go-html-template
{{ sha1 "Hello world" }} → 7b502c3a1f48c8609ae212cdfb639dee39673f5e
```
