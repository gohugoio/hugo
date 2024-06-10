---
title: Next
description: Returns the next page in a global page collection, relative to the given page. 
categories: []
keywords: []
action:
  related:
    - methods/page/Prev
    - methods/page/NextInSection
    - methods/page/PrevInSection
    - methods/pages/Next
    - methods/pages/Prev
  returnType: page.Page
  signatures: [PAGE.Next]
toc: true
---

The behavior of the `Prev` and `Next` methods on a `Page` object is probably the reverse of what you expect.

With this content structure:

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
{{ with .Next }}
  <a href="{{ .RelPermalink }}">Prev</a>
{{ end }}

{{ with .Prev }}
  <a href="{{ .RelPermalink }}">Next</a>
{{ end }}
```

## Compare to Pages methods

{{% include "methods/_common/next-prev-on-page-vs-next-prev-on-pages.md" %}}
