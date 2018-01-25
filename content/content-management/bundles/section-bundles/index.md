---
title : "Section Bundles"
description : "Organization of sections as bundles"
lastmod : 2018-01-25T16:39:52-05:00
linktitle : "Section Bundles"
keywords : ["section", "bundles"]
categories : ["content management", "bundles"]
weight : 4002
draft : false
toc : true
---

A _Section Bundle_ is any directory at any hierarchy within the
`content/` directory, that contains at least an `_index.md` file (**not
`index.md`**). This `_index.md` can also be directly under the
`content/` directory.

{{% note %}}
Here `md` (markdown) is used just as an example. You can use any file
type as a content resource as long as it is a MIME type recognized by
Hugo (`json` files will, as one example, work fine). If you want to
get exotic, you can define your own media type.
{{% /note %}}


## Examples of Section Bundle organization {#examples-of-section-bundle-organization}

```text
content/
├── section-bundle-1
│   ├── section-content1.md
│   ├── section-content2.md
│   ├── image1.jpg
│   ├── image2.png
│   └── _index.md
└── section-bundle-2
    ├── _index.md
    └── page-bundle-1
        └── index.md
```

In the above example `content/` directory, there are two section
bundles (and a page bundle):

`section-bundle-1`
: This section bundle has the `_index.md`, two
    other content Markdown files and two image files.

`section-bundle-2`
: This section bundle has the `_index.md` and a
    nested page bundle.

{{% note %}}
The hierarchy depth at which a section bundle is created does not matter,
as long as it's not inside another **section** bundle.
{{% /note %}}
