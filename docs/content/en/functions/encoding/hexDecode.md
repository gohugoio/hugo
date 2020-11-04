---
title: hexDecode
description: Decode a hex string.
godocref:
date: 2020-11-03
publishdate: 2020-11-03
lastmod: 2020-11-03
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
relatedfuncs: [hexEncode]
signature: ["hexDecode INPUT"]
workson: []
hugoversion: "v0.79.0"
deprecated: false
draft: false
aliases: []
---

`hexDecode` decodes the given hex (base 16) input string.

For example:

```
{{ 42 | hexEncode | hexDecode }} → "42"
{{ "48656c6c6f20776f726c64" | hexDecode }} → "Hello world"
```

An example of decoding a hex color into an RGB slice:
```
{{ $v := hexDecode "c0d0e0" }}
{{ slice (index $v 0) (index $v 1) (index $v 2) }} → [192, 208, 224]
```
