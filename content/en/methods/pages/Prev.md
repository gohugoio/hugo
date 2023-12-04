---
title: Prev
description: Returns the previous page in a local page collection, relative to the given page.
categories: []
keywords: []
action:
  related:
    - methods/pages/Next
    - methods/page/Next
    - methods/page/NextInSection
    - methods/page/Prev
    - methods/page/PrevInSection
  returnType: hugolib.pageStates
  signatures: [PAGES.Prev PAGE]
toc: true
---

The behavior of the `Prev` and `Next` methods on a `Pages` objects is probably the reverse of what you expect.

With this content structure and the page collection sorted by weight in ascending order:

```text
content/
├── pages/
│   ├── _index.md
│   ├── page-1.md   <-- front matter: weight = 10
│   ├── page-2.md   <-- front matter: weight = 20
│   └── page-3.md   <-- front matter: weight = 30
└── _index.md
```

When you visit page-2:

- The `Prev` method points to page-3
- The `Next` method points to page-1

{{% note %}}
Use the opposite label in your navigation links as shown in the example below.
{{% /note %}}

```go-html-template
{{ $pages := where .Site.RegularPages.ByWeight "Section" "pages" }}

{{ with $pages.Next . }}
  <a href="{{ .RelPermalink }}">Previous</a>
{{ end }}

{{ with $pages.Prev . }}
  <a href="{{ .RelPermalink }}">Next</a>
{{ end }}
```

## Compare to Page methods

{{% include "methods/_common/next-prev-on-page-vs-next-prev-on-pages.md" %}}
