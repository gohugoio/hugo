---
title: resources.Fingerprint
description: Cryptographically hashes the content of the given resource.
categories: []
keywords: []
action:
  aliases: [fingerprint]
  related:
    - functions/js/Build
    - functions/resources/Babel
    - functions/resources/Minify
    - functions/resources/PostCSS
    - functions/resources/PostProcess
    - functions/resources/ToCSS
  returnType: resource.Resource
  signatures: ['resources.Fingerprint [ALGORITHM] RESOURCE']
---

```go-html-template
{{ with resources.Get "js/main.js" }}
  {{ with . | fingerprint "sha256" }}
    <script src="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous"></script>
  {{ end }}
{{ end }}
```

Hugo renders this to something like:

```html
<script src="/js/main.62e...df1.js" integrity="sha256-Yuh...rfE=" crossorigin="anonymous"></script>
```

Although most commonly used with CSS and JavaScript resources, you can use the `resources.Fingerprint` function with any resource type.

The hash algorithm may be one of `md5`, `sha256` (default), `sha384`, or `sha512`.

After cryptographically hashing the resource content:

1. The values returned by the `.Permalink` and `.RelPermalink` methods include the hash sum
2. The resource's `.Data.Integrity` method returns a [Subresource Integrity] (SRI) value consisting of the name of the hash algorithm, one hyphen, and the base64-encoded hash sum

[Subresource Integrity]: https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
