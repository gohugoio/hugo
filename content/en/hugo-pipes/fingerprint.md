---
title: Fingerprint
linkTitle: Fingerprinting and SRI
description: Process a given resource, adding a hash string of the resource's content.
date: 2018-07-14
publishdate: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 70
weight: 70
sections_weight: 70
signature: ["resources.Fingerprint RESOURCE [ALGORITHM]", "fingerprint RESOURCE [ALGORITHM]"]
---

## Usage

Fingerprinting and [SRI](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity) can be applied to any asset file using `resources.Fingerprint` which takes two arguments, the resource object and an optional [hash algorithm](https://en.wikipedia.org/wiki/Secure_Hash_Algorithms).

The default hash algorithm is `sha256`. Other available algorithms are `sha384` and (as of Hugo `0.55`) `sha512` and `md5`.

Any so processed asset will bear a `.Data.Integrity` property containing an integrity string, which is made up of the name of the hash algorithm, one hyphen and the base64-encoded hash sum.

```go-html-template
{{ $js := resources.Get "js/global.js" }}
{{ $secureJS := $js | resources.Fingerprint "sha512" }}
<script src="{{ $secureJS.Permalink }}" integrity="{{ $secureJS.Data.Integrity }}"></script>
```
