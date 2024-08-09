---
title: hash.FNV32a
description: Returns the 32-bit FNV (Fowler–Noll–Vo) non-cryptographic hash of the given string.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/hash/Xxhash
    - functions/crypto/HMAC
    - functions/crypto/MD5
    - functions/crypto/SHA1
    - functions/crypto/SHA256
  returnType: int
  signatures: [hash.FNV32a STRING]
aliases: [/functions/crypto.fnv32a]
---

```go-html-template
{{ hash.FNV32a "Hello world" }} → 1498229191
```
