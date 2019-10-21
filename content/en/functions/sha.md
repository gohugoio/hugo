---
title: sha
# linktitle: sha
description: Hashes the given input and returns either an SHA1 or SHA256 checksum.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [sha,checksum]
signature: ["sha1 INPUT", "sha256 INPUT"]
workson: []
hugoversion:
relatedfuncs: [md5]
deprecated: false
aliases: [sha1, sha256]
---

`sha1` hashes the given input and returns its SHA1 checksum.

```
{{ sha1 "Hello world, gophers!" }}
<!-- returns the string "c8b5b0e33d408246e30f53e32b8f7627a7a649d4" -->
```

`sha256` hashes the given input and returns its SHA256 checksum.

```
{{ sha256 "Hello world, gophers!" }}
<!-- returns the string "6ec43b78da9669f50e4e422575c54bf87536954ccd58280219c393f2ce352b46" -->
```
