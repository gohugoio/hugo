---
title: collections.Shuffle
description: Returns a random permutation of a given array or slice.
categories: []
keywords: [random]
params:
  functions_and_methods:
    aliases: [shuffle]
    returnType: any
    signatures: [collections.Shuffle COLLECTION]
aliases: [/functions/shuffle]
---

```go-html-template
{{ collections.Shuffle (slice "a" "b" "c") }} → [b a c] 
```

The result will vary from one build to the next.

To render an unordered list of 5 random pages from a page collection:

```go-html-template
<ul>
  {{ $p := site.RegularPages }}
  {{ range $p | collections.Shuffle | first 5 }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

{{< new-in 0.149.0 />}}

Using the [`collections.D`][] function for the same task is significantly faster.

[`collections.D`]: /functions/collections/D/
