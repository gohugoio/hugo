---
title : "Page Bundles"
description : "Organization of individual pages as bundles"
lastmod : 2018-01-25T16:38:34-05:00
linktitle : "Page Bundles"
keywords : ["page", "bundles"]
categories : ["content management", "bundles"]
weight : 4001
draft : false
toc : true
---

A _Page Bundle_ is any directory at any hierarchy within the
`content/` directory, that contains at least an `index.md` file (**not
`_index.md`**).

{{% note %}}
Here `md` (markdown) is used just as an example. You can use any file
type as a content resource as long as it is a MIME type recognized by
Hugo (`json` files will, as one example, work fine). If you want to
get exotic, you can define your own media type.
{{% /note %}}


## Examples of Page Bundle organization {#examples-of-page-bundle-organization}

```text
content/
├── page-bundle-1
│   ├── content1.md
│   ├── content2.md
│   ├── image1.jpg
│   ├── image2.png
│   └── index.md
└── some-section
    ├── ..
    ├── ..
    └── page-bundle-2
        └── index.md
```

In the above example `content/` directory, there are two page bundles:

`page-bundle-1`
: This page bundle has the `index.md`, two other
    content Markdown files and two image files.

`page-bundle-2`
: This page bundle is nested in a section. This
    bundle has only the `index.md`.

{{% note %}}
The hierarchy depth at which a page bundle is created does not matter,
as long as it's not inside another **page** bundle.
{{% /note %}}


## <span class="todo TODO_">TODO </span> Headless Page Bundle {#headless-page-bundle}
