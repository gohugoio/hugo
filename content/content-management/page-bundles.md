---
title : "Page Bundles"
description : "Content organization using Page Bundles"
date : 2018-01-24T13:09:00-05:00
lastmod : 2018-01-26T10:58:36-05:00
linktitle : "Page Bundles"
keywords : ["page", "bundle", "leaf", "branch"]
categories : ["content management", "bundle"]
draft : false
toc : true
menu :
  docs:
    identifier : "page-bundles"
    parent : "content-management"
    weight : 11
---

Page Bundles are a way to organize the content files. It's useful for
cases where a page or section's content needs to be split into
multiple content pages for convenience or has associated attachments
like documents or images.

A Page Bundle can be one of two types:

-   Leaf Bundle
-   Branch Bundle

|                 | Leaf Bundle                                            | Branch Bundle                                           |
|-----------------|--------------------------------------------------------|---------------------------------------------------------|
| Usage           | Collection of content and attachments for single pages | Collection of content and attachments for section pages |
| Index file name | `index.md` [^fn:1]                                     | `_index.md` [^fn:1]                                     |
| Layout type     | `single`                                               | `list`                                                  |
| Example         | `content/posts/my-post/index.md`                       | `content/posts/_index.md`                               |


## Leaf Bundles {#leaf-bundles}

A _Leaf Bundle_ is a directory at any hierarchy within a [Section](/content-management/sections/)
directory, that contains at least an **`index.md`** file.

{{% note %}}
Here `md` (markdown) is used just as an example. You can use any file
type as a content resource as long as it is a MIME type recognized by
Hugo (`json` files will, as one example, work fine). If you want to
get exotic, you can define your own media type.
{{% /note %}}


### Examples of Leaf Bundle organization {#examples-of-leaf-bundle-organization}

```text
content/
├── posts
│   ├── my-post
│   │   ├── content1.md
│   │   ├── content2.md
│   │   ├── image1.jpg
│   │   ├── image2.png
│   │   └── index.md
│   └── my-another-post
│       └── index.md
│
└── another-section
    ├── ..
    └── not-a-leaf-bundle
        ├── ..
        └── another-leaf-bundle
            └── index.md
```

In the above example `content/` directory, there are three leaf
bundles:

my-post
: This leaf bundle has the `index.md`, two other content
    Markdown files and two image files.

my-another-post
: This leaf bundle has only the `index.md`.

another-leaf-bundle
: This leaf bundle is nested under couple of
    directories. This bundle also has only the `index.md`.

{{% note %}}
The hierarchy depth at which a leaf bundle is created does not matter,
as long as (1) it's not inside another **leaf** bundle (2) it's not
directly under the `content/` directory.
{{% /note %}}


### Headless Leaf Bundle {#headless-leaf-bundle}

**TODO**


## Branch Bundles {#branch-bundles}

A _Branch Bundle_ is any directory at any hierarchy within the
`content/` directory, that contains at least an **`_index.md`** file.

This `_index.md` can also be directly under the `content/` directory.

{{% note %}}
Here `md` (markdown) is used just as an example. You can use any file
type as a content resource as long as it is a MIME type recognized by
Hugo (`json` files will, as one example, work fine). If you want to
get exotic, you can define your own media type.
{{% /note %}}


### Examples of Branch Bundle organization {#examples-of-branch-bundle-organization}

```text
content/
├── branch-bundle-1
│   ├── branch-content1.md
│   ├── branch-content2.md
│   ├── image1.jpg
│   ├── image2.png
│   └── _index.md
└── branch-bundle-2
    ├── _index.md
    └── a-leaf-bundle
        └── index.md
```

In the above example `content/` directory, there are two branch
bundles (and a leaf bundle):

`branch-bundle-1`
: This branch bundle has the `_index.md`, two
    other content Markdown files and two image files.

`branch-bundle-2`
: This branch bundle has the `_index.md` and a
    nested leaf bundle.

{{% note %}}
The hierarchy depth at which a branch bundle is created does not matter.
{{% /note %}}

[^fn:1]: The `.md` extension is just an example. The extension can be `.html`, `.json` or any of any valid MIME type.
