---
title: collections.Union
description: Returns a slice containing the unique elements from two given slices.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [union]
    returnType: '[]any'
    signatures: [collections.Union SLICE1 SLICE2]
aliases: [/functions/union] 
---



```go-html-template
{{ union (slice 1 2 3) (slice 3 4 5) }} → [1 2 3 4 5]
{{ union (slice 1 2 3) nil }}           → [1 2 3]
{{ union nil (slice 1 2 3) }}           → [1 2 3]
{{ union nil nil }}                     → []
```

## OR filter in where query

This is also very useful to use as `OR` filters when combined with where:

```go-html-template
{{ $pages := where .Site.RegularPages "Type" "not in" (slice "page" "about") }}
{{ $pages = $pages | union (where .Site.RegularPages "Params.pinned" true) }}
{{ $pages = $pages | intersect (where .Site.RegularPages "Params.images" "!=" nil) }}
```

The above fetches regular pages not of `page` or `about` type unless they are pinned. And finally, we exclude all pages with no `images` set in Page parameters.

See [intersect](/functions/collections/intersect) for `AND`.
