---
title: collections.Intersect
description: Returns the common elements of two arrays or slices, in the same order as the first array.
categories: []
keywords: []
action:
  aliases: [intersect]
  related:
    - functions/collections/Complement
    - functions/collections/SymDiff
    - functions/collections/Union
  returnType: any
  signatures: [collections.Intersect SET1 SET2]
aliases: [/functions/intersect]
---

A useful example is to use it as `AND` filters when combined with where:

```go-html-template
{{ $pages := where .Site.RegularPages "Type" "not in" (slice "page" "about") }}
{{ $pages := $pages | union (where .Site.RegularPages "Params.pinned" true) }}
{{ $pages := $pages | intersect (where .Site.RegularPages "Params.images" "!=" nil) }}
```

The above fetches regular pages not of `page` or `about` type unless they are pinned. And finally, we exclude all pages with no `images` set in Page parameters.

See [union](/functions/collections/union) for `OR`.
