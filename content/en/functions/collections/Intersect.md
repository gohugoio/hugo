---
title: collections.Intersect
description: Returns a slice containing the common elements found in two given slices, in the same order as the first slice.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [intersect]
    returnType: '[]any'
    signatures: [collections.Intersect SLICE1 SLICE2]
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
