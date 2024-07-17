---
title: hash.XxHash
description: Returns the 64-bit xxHash non-cryptographic hash of the given string.
categories: []
keywords: []
action:
  aliases: [xxhash]
  related:
    - functions/hash/FNV32a
    - functions/crypto/HMAC
    - functions/crypto/MD5
    - functions/crypto/SHA1
    - functions/crypto/SHA256
  returnType: string
  signatures: [hash.XxHash STRING]
---

```go-html-template
{{ hash.XxHash "Hello world" }} → c500b0c912b376d8
```
