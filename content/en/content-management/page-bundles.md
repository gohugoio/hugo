---
title : "Page Bundles"
description : "Content organization using Page Bundles"
date : 2018-01-24T13:09:00-05:00
lastmod : 2018-01-28T22:26:40-05:00
linktitle : "Page Bundles"
keywords : ["page", "bundle", "leaf", "branch"]
categories : ["content management"]
toc : true
menu :
  docs:
    identifier : "page-bundles"
    parent : "content-management"
    weight : 11
---

Page Bundles are a way to group [Page Resources](/content-management/page-resources/).

A Page Bundle can be one of:

-   Leaf Bundle (leaf means it has no children)
-   Branch Bundle (home page, section, taxonomy terms, taxonomy list)

|                                     | Leaf Bundle                                              | Branch Bundle                                                                                                                                                                                                      |
|-------------------------------------|----------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------  |
| Usage                               | Collection of content and attachments for single pages   | Collection of attachments for section pages (home page, section, taxonomy terms, taxonomy list)                                                                                                                    |
| Index file name                     | `index.md` [^fn:1]                                       | `_index.md` [^fn:1]                                                                                                                                                                                                |
| Allowed Resources                   | Page and non-page (like images, pdf, etc.) types         | Only non-page (like images, pdf, etc.) types                                                                                                                                                                       |
| Where can the Resources live?       | At any directory level within the leaf bundle directory. | Only in the directory level **of** the branch bundle directory i.e. the directory containing the `_index.md` ([ref](https://discourse.gohugo.io/t/question-about-content-folder-structure/11822/4?u=kaushalmodi)). |
| Layout type                         | `single`                                                 | `list`                                                                                                                                                                                                             |
| Nesting                             | Does not allow nesting of more bundles under it          | Allows nesting of leaf or branch bundles under it                                                                                                                                                                  |
| Example                             | `content/posts/my-post/index.md`                         | `content/posts/_index.md`                                                                                                                                                                                          |
| Content from non-index page files... | Accessed only as page resources                          | Accessed only as regular pages                                                                                                                                                                                     |


## Leaf Bundles {#leaf-bundles}

A _Leaf Bundle_ is a directory at any hierarchy within the `content/`
directory, that contains an **`index.md`** file.

### Examples of Leaf Bundle organization {#examples-of-leaf-bundle-organization}

```text
content/
├── about
│   ├── index.md
├── posts
│   ├── my-post
│   │   ├── content1.md
│   │   ├── content2.md
│   │   ├── image1.jpg
│   │   ├── image2.png
│   │   └── index.md
│   └── my-other-post
│       └── index.md
│
└── another-section
    ├── ..
    └── not-a-leaf-bundle
        ├── ..
        └── another-leaf-bundle
            └── index.md
```

In the above example `content/` directory, there are four leaf
bundles:

about
: This leaf bundle is at the root level (directly under
    `content` directory) and has only the `index.md`.

my-post
: This leaf bundle has the `index.md`, two other content
    Markdown files and two image files.

my-other-post
: This leaf bundle has only the `index.md`.

another-leaf-bundle
: This leaf bundle is nested under couple of
    directories. This bundle also has only the `index.md`.

{{% note %}}
The hierarchy depth at which a leaf bundle is created does not matter,
as long as it is not inside another **leaf** bundle.
{{% /note %}}


### Headless Bundle {#headless-bundle}

A headless bundle is a bundle that is configured to not get published
anywhere:

-   It will have no `Permalink` and no rendered HTML in `public/`.
-   It will not be part of `.Site.RegularPages`, etc.

But you can get it by `.Site.GetPage`. Here is an example:

```go-html-template
{{ $headless := .Site.GetPage "/some-headless-bundle" }}
{{ $reusablePages := $headless.Resources.Match "author*" }}
<h2>Authors</h2>
{{ range $reusablePages }}
    <h3>{{ .Title }}</h3>
    {{ .Content }}
{{ end }}
```

_In this example, we are assuming the `some-headless-bundle` to be a headless
   bundle containing one or more **page** resources whose `.Name` matches
   `"author*"`._

Explanation of the above example:

1. Get the `some-headless-bundle` Page "object".
2. Collect a *slice* of resources in this *Page Bundle* that matches
   `"author*"` using `.Resources.Match`.
3. Loop through that *slice* of nested pages, and output their `.Title` and
   `.Content`.

---

A leaf bundle can be made headless by adding below in the Front Matter
(in the `index.md`):

```toml
headless = true
```

{{% note %}}
Only leaf bundles can be made headless.
{{% /note %}}

There are many use cases of such headless page bundles:

-   Shared media galleries
-   Reusable page content "snippets"


## Branch Bundles {#branch-bundles}

A _Branch Bundle_ is any directory at any hierarchy within the
`content/` directory, that contains at least an **`_index.md`** file.

This `_index.md` can also be directly under the `content/` directory.

{{% note %}}
Here `md` (markdown) is used just as an example. You can use any file
type as a content resource as long as it is a content type recognized by Hugo.
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
The hierarchy depth at which a branch bundle is created does not
matter.
{{% /note %}}

[^fn:1]: The `.md` extension is just an example. The extension can be `.html`, `.json` or any of any valid MIME type.
