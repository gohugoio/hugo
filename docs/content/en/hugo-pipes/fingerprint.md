---
title: Fingerprint
linkTitle: Fingerprinting and SRI hashing
description: Process a given resource, adding a hash string of the resource's content.
categories: [asset management]
keywords: []
menu:
  docs:
    parent: hugo-pipes
    weight: 100
weight: 100
action:
  aliases: [fingerprint]
  returnType: resource.Resource
  signatures: ['resources.Fingerprint [ALGORITHM] RESOURCE']
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
