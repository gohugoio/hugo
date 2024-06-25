---
title: Page bundles
description: Use page bundles to logically associate one or more resources with content.
categories: [content management]
keywords: [page,bundle,leaf,branch]
menu :
  docs:
    parent: content-management
    weight: 30
weight: 30
toc: true
---

## Introduction

A page bundle is a directory that encapsulates both content and associated resources.

By way of example, this site has an "about" page and a "privacy" page:

```text
content/
├── about/
│   ├── index.md
│   └── welcome.jpg
└── privacy.md
```

The "about" page is a page bundle. It logically associates a resource with content by bundling them together. Resources within a page bundle are [page resources], accessible with the [`Resources`] method on the `Page` object.

Page bundles are either _leaf bundles_ or _branch bundles_.

leaf bundle
: A _leaf bundle_ is a directory that contains an index.md file and zero or more resources. Analogous to a physical leaf, a leaf bundle is at the end of a branch. It has no descendants.

branch bundle
: A _branch bundle_ is a directory that contains an _index.md file and zero or more resources. Analogous to a physical branch, a branch bundle may have descendants including leaf bundles and other branch bundles. Top level directories with or without _index.md files are also branch bundles. This includes the home page.

{{% note %}}
In the definitions above and the examples below, the extension of the index file depends on the [content format]. For example, use index.md for Markdown content, index.html for HTML content, index.adoc for AsciiDoc content, etc.

[content format]: /getting-started/glossary/#content-format
{{% /note %}}

## Comparison

Page bundle characteristics vary by bundle type.

|                     | Leaf bundle                                             | Branch bundle                                           |
|---------------------|---------------------------------------------------------|---------------------------------------------------------|
| Index file          | index.md                                                | _index.md                                               |
| Example             | content/about/index.md                                  | content/posts/_index.md                                 |
| [Page kinds]        | `page`                                                  | `home`, `section`, `taxonomy`, or `term`                |
| Template types       | [single]                                                | [home], [section], [taxonomy], or [term]           |
| Descendant pages    | None                                                    | Zero or more                                            |
| Resource location   | Adjacent to the index file or in a nested subdirectory  | Same as a leaf bundles, but excludes descendant bundles |
| [Resource types]    | `page`, `image`, `video`, etc.                          | all but `page`                                          |

[single]: /templates/types/#single
[home]: /templates/types/#home
[section]: /templates/types/#section
[taxonomy]: /templates/types/#taxonomy
[term]: /templates/types/#term

Files with [resource type] `page` include content written in Markdown, HTML, AsciiDoc, Pandoc, reStructuredText, and Emacs Org Mode. In a leaf bundle, excluding the index file, these files are only accessible as page resources. In a branch bundle, these files are only accessible as content pages.

## Leaf bundles

A _leaf bundle_ is a directory that contains an index.md file and zero or more resources. Analogous to a physical leaf, a leaf bundle is at the end of a branch. It has no descendants.

```text
content/
├── about
│   └── index.md
├── posts
│   ├── my-post
│   │   ├── content-1.md
│   │   ├── content-2.md
│   │   ├── image-1.jpg
│   │   ├── image-2.png
│   │   └── index.md
│   └── my-other-post
│       └── index.md
└── another-section
    ├── foo.md
    └── not-a-leaf-bundle
        ├── bar.md
        └── another-leaf-bundle
            └── index.md
```

There are four leaf bundles in the example above:

about
: This leaf bundle does not contain any page resources.

my-post
: This leaf bundle contains an index file, two resources of [resource type] `page`, and two resources of resource type `image`.

- content-1, content-2

  These are resources of resource type `page`, accessible via the [`Resources`] method on the `Page` object. Hugo will not render these as individual pages.

- image-1, image-2

  These are resources of resource type `image`, accessible via the `Resources` method on the `Page` object

my-other-post
: This leaf bundle does not contain any page resources.

another-leaf-bundle
: This leaf bundle does not contain any page resources.

{{% note %}}
Create leaf bundles at any depth within the content directory, but a leaf bundle may not contain another bundle. Leaf bundles do not have descendants.
{{% /note %}}

## Branch bundles

A _branch bundle_ is a directory that contains an _index.md file and zero or more resources. Analogous to a physical branch, a branch bundle may have descendants including leaf bundles and other branch bundles. Top level directories with or without _index.md files are also branch bundles. This includes the home page.

```text
content/
├── branch-bundle-1/
│   ├── _index.md
│   ├── content-1.md
│   ├── content-2.md
│   ├── image-1.jpg
│   └── image-2.png
├── branch-bundle-2/
│   ├── a-leaf-bundle/
│   │   └── index.md
│   └── _index.md
└── _index.md
```

There are three branch bundles in the example above:

home page
: This branch bundle contains an index file, two descendant branch bundles, and no resources.

branch-bundle-1
:  This branch bundle contains an index file, two resources of [resource type] `page`, and two resources of resource type `image`.

branch-bundle-2
: This branch bundle contains an index file and a leaf bundle.

{{% note %}}
Create branch bundles at any depth within the content directory, but a leaf bundle may not contain another bundle. Leaf bundles do not have descendants.
{{% /note %}}


## Headless bundles

Use [build options] in front matter to create an unpublished leaf or branch bundle whose content and resources you can include in other pages.

[`Resources`]: /methods/page/resources/
[build options]: content-management/build-options/
[page kinds]: /getting-started/glossary/#page-kind
[page resources]: /content-management/page-resources/
[resource type]: /getting-started/glossary/#resource-type
[resource types]: /getting-started/glossary/#resource-type
