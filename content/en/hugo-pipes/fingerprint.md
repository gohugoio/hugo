---
title: Fingerprinting and SRI
description: Hugo Pipes allows Fingerprinting and Subresource Integrity.
date: 2018-07-14
publishdate: 2018-07-14
lastmod: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 70
weight: 70
sections_weight: 70
draft: false
---


Fingerprinting and [SRI](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity) can be applied to any asset file using `resources.Fingerprint` which takes two arguments, the resource object and a [hash function](https://en.wikipedia.org/wiki/Cryptographic_hash_function). 

The default hash function is `sha256`. Other available functions are `sha384` (from Hugo `0.55`), `sha512` and `md5`.

Any so processed asset will bear a `.Data.Integrity` property containing an integrity string, which is made up of the name of the hash function, one hyphen and the base64-encoded hash sum.

```go-html-template
{{ $js := resources.Get "js/global.js" }}
{{ $secureJS := $js | resources.Fingerprint "sha512" }}
<script type="text/javascript" src="{{ $secureJS.Permalink }}" integrity="{{ $secureJS.Data.Integrity }}"></script>
```
