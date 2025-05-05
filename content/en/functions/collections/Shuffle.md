---
title: collections.Shuffle
description: Returns a random permutation of a given array or slice.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [shuffle]
    returnType: any
    signatures: [collections.Shuffle COLLECTION]
aliases: [/functions/shuffle]
---

```go-html-template
{{ collections.Shuffle (seq 1 2 3) }} → [3 1 2] 
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

For better performance with large collections, use the [`math.Rand`] and [`collections.Index`] functions instead:

```go-html-template
<ul>
  {{ $p := site.RegularPages }}
  {{ range seq 5 }}
    {{ with math.Rand | mul $p.Len | math.Floor | int }}
      {{ with index $p . }}
        <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
      {{ end }}
    {{ end }}
  {{ end }}
</ul>
```

[`collections.Index`]:/functions/collections/indexfunction/
[`math.Rand`]: /functions/math/rand/
