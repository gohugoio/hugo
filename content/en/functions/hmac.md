---
title: hmac
linktitle: hmac
description: Compute the cryptographic checksum of a message.
date: 2020-05-29
publishdate: 2020-05-29
lastmod: 2020-05-29
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [hmac,checksum]
signature: ["hmac HASH_TYPE KEY MESSAGE"]
workson: []
hugoversion:
relatedfuncs: [hmac]
deprecated: false
aliases: [hmac]
---

`hmac` returns a cryptographic hash that uses a key to sign a message.

```
{{ hmac "sha256" "Secret key" "Hello world, gophers!"}},
<!-- returns the string "b6d11b6c53830b9d87036272ca9fe9d19306b8f9d8aa07b15da27d89e6e34f40"
```

Supported hash functions:

* md5
* sha1
* sha256
* sha512
