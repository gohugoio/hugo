---
title: IndexOf
description: Returns the zero-based index of the given page within the given page collection.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [PAGES.IndexOf PAGE]
---

{{< new-in 0.166.0 />}}

If the given page is not in the page collection, the `IndexOf` method returns `-1`.

With this content structure:

```tree
content/
├── posts/
│   ├── _index.md
│   ├── post-1.md   <-- front matter: weight = 10
│   ├── post-2.md   <-- front matter: weight = 20
│   └── post-3.md   <-- front matter: weight = 30
└── _index.md
```

And this template:

```go-html-template {file="layouts/posts/page.html"}
{{ $pages := .CurrentSection.Pages.ByWeight }}
{{ $index := add ($pages.IndexOf .) 1 }}
<p>This is post {{ $index }} of {{ $pages.Len }} in {{ .CurrentSection.LinkTitle }}.</p>
```

When you visit post-2, Hugo renders:

```html
<p>This is post 2 of 3 in Posts.</p>
```

You can also use the `IndexOf` method to build previous and next navigation links. Combined with the [`index`][] and [`add`][]/[`sub`][] functions, use it to get the previous and next page in the same page collection:

```go-html-template
{{ $pages := .CurrentSection.Pages.ByWeight }}
{{ $index := $pages.IndexOf . }}

{{ if ge $index 0 }}
  {{ with index $pages (add $index 1) }}
    <a href="{{ .RelPermalink }}">Previous</a>
  {{ end }}

  {{ with index $pages (sub $index 1) }}
    <a href="{{ .RelPermalink }}">Next</a>
  {{ end }}
{{ end }}
```

This is equivalent to using the [`Prev`][] and [`Next`][] methods:

```go-html-template
{{ $pages := .CurrentSection.Pages.ByWeight }}

{{ with $pages.Prev . }}
  <a href="{{ .RelPermalink }}">Previous</a>
{{ end }}

{{ with $pages.Next . }}
  <a href="{{ .RelPermalink }}">Next</a>
{{ end }}
```

> [!TIP]
> Unlike `Prev` and `Next`, the `IndexOf` approach also gives you the page's position within the collection, as shown in the example at the beginning of this page.

[`Next`]: /methods/pages/next/
[`Prev`]: /methods/pages/prev/
[`add`]: /functions/math/add/
[`index`]: /functions/collections/indexfunction/
[`sub`]: /functions/math/sub/
