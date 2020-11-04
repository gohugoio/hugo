---
title: hexEncode
description: Encode to hex.
godocref:
date: 2020-11-03
publishdate: 2020-11-03
lastmod: 2020-11-03
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
relatedfuncs: [hexDecode]
signature: ["hexEncode INPUT"]
workson: []
hugoversion: "v0.79.0"
deprecated: false
draft: false
aliases: []
---

`hexEncode` encodes a given input into a hex (base 16) string.

For example:

```
{{ "Hello world" | hexEncode }} → "48656c6c6f20776f726c64"
```
