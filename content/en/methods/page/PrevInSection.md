---
title: PrevInSection
description: Returns the previous page within a section, relative to the given page.  
categories: []
keywords: []
action:
  related:
    - methods/page/NextInSection
    - methods/page/Next
    - methods/pages/Next
    - methods/page/Prev
    - methods/pages/Prev
  returnType: page.Page
  signatures: [PAGE.PrevInSection]
---


The behavior of the `PrevInSection` and `NextInSection` methods on a `Page` object is probably the reverse of what you expect.

With this content structure:

```text
content/
├── books/
│   ├── _index.md
│   ├── book-1.md
│   ├── book-2.md
│   └── book-3.md
├── films/
│   ├── _index.md
│   ├── film-1.md
│   ├── film-2.md
│   └── film-3.md
└── _index.md
```

When you visit book-2:

- The `PrevInSection` method points to book-3
- The `NextInSection` method points to book-1

{{% note %}}
Use the opposite label in your navigation links as shown in the example below.
{{% /note %}}

```go-html-template
{{ with .NextInSection }}
  <a href="{{ .RelPermalink }}">Previous in section</a>
{{ end }}

{{ with .PrevInSection }}
  <a href="{{ .RelPermalink }}">Next in section</a>
{{ end }}
```


The `PrevInSection` and `NextInSection` methods uses the default page sort. You can change the sort direction in [Page Config](getting-started/configuration/#configure-page). For more flexibility, use the [Next] and [Prev] methods on the Page object.

[date]: /methods/page/date/
[weight]: /methods/page/weight/
[linkTitle]: /methods/page/linktitle/
[title]: /methods/page/title/
[Next]: /methods/page/next/
[Prev]: /methods/page/prev/
