---
title: ref and relref
linktitle: ref and relref
description: Looks up a content page by relative path or logical name to return the content page's permalink.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [cross references, anchors]
signature:
workson: []
hugoversion:
relatedfuncs: [relref]
deprecated: false
aliases: [/functions/ref/,/functions/relref/]
---

These two functions looks up a content page by relative path (`relref`) or logical name (`ref`) to return the permalink. Both functions require a `Page` object (usually satisfied with a "`.`"):

```golang
{{ relref . "about.md" }}
```

These functions are used in two of Hugo's built-in shortcodes. You can see basic usage examples of both `ref` and `relref` in the [shortcode documentation](/content-management/shortcodes/#ref-and-relref).

For an extensive explanation of how to leverage `ref` and `relref` for content management, see [Cross References](/content-management/cross-references/).