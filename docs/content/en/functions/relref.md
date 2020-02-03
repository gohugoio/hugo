---
title: relref
# linktitle: relref
description: Looks up a content page by relative path.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2019-12-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [cross references, anchors]
signature: ["relref . CONTENT"]
workson: []
hugoversion:
relatedfuncs: [ref]
deprecated: false
aliases: []
---

`ref` and `relref` look up a content page by logical name (`ref`) or relative path (`relref`) to return the permalink:

```
{{ relref . "about.md" }}
```

{{% note "Usage Note" %}}
`relref` looks up Hugo "Regular Pages" only. It can't be used for the homepage, section pages, etc.
{{% /note %}}

It is also possible to pass additional arguments to link to another language or an alternative output format. Therefore, pass a map of arguments instead of just the path. 

``` 
{{ relref . (dict "path" "about.md" "lang" "ja" "outputFormat" "rss") }}
```

These functions are used in two of Hugo's built-in shortcodes. You can see basic usage examples of both `ref` and `relref` in the [shortcode documentation](/content-management/shortcodes/#ref-and-relref).

For an extensive explanation of how to leverage `ref` and `relref` for content management, see [Cross References](/content-management/cross-references/).
