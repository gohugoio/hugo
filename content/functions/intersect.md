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
#tags: []
signature: ["intersect SET1 SET2"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Given two arrays (or slices), `intersect` returns the common elements. The elements supported are strings, integers, and floats (only float64).

A useful example of `intersect` functionality is a "related posts" block. `isset` allows us to create a list of links to other posts that have tags that intersect with the tags in the current post.

The following is an example of a "related posts" [partial template][partials] that could be added to a [single page template][single]:

{{% code file="layouts/partials/related-posts.html" download="related-posts.html" %}}
```html
<ul>
{{ $page_link := .Permalink }}
{{ $tags := .Params.tags }}
{{ range .Site.Pages }}
    {{ $page := . }}
    {{ $has_common_tags := intersect $tags .Params.tags | len | lt 0 }}
    {{ if and $has_common_tags (ne $page_link $page.Permalink) }}
        <li><a href="{{ $page.Permalink }}">{{ $page.Title }}</a></li>
    {{ end }}
{{ end }}
</ul>
```
{{% /code %}}

[partials]: /templates/partials/
[single]: /templates/single-page-templates/
