---
title: intersect
linktitle: intersect
description: Returns the common elements of two arrays or slices.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [collections,intersect,union,complement,symdiff]
signature: ["intersect SET1 SET2"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---
An useful example is to use it as `AND` filters when combined with where:

## AND filter in where query

```
{{ $pages := where .Site.RegularPages "Type" "not in" (slice "page" "about") }}
{{ $pages := $pages | union (where .Site.RegularPages "Params.pinned" true) }}
{{ $pages := $pages | intersect (where .Site.RegularPages "Params.images" "!=" nil) }}
```

The above fetches regular pages not of `page` or `about` type unless they are pinned. And finally, we exclude all pages with no `images` set in Page params.

See [union](/functions/union) for `OR`.


[partials]: /templates/partials/
[single]: /templates/single-page-templates/
