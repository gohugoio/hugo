---
title: sha1
description: Hashes the given input and returns its SHA1 checksum.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: crypto
relatedFuncs:
  - crypto.FNV32a
  - crypto.HMAC
  - crypto.MD5
  - crypto.SHA1
  - crypto.SHA256
signature:
  - crypto.SHA1 INPUT
  - sha1 INPUT
aliases: [sha]
---

```go-html-template
{{ sha1 "Hello world" }} → 7b502c3a1f48c8609ae212cdfb639dee39673f5e
```
