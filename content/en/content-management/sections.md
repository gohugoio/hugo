---
title: Content Sections
linkTitle: Sections
description: Hugo generates a **section tree** that matches your content.
categories: [content management]
keywords: [lists,sections,content types,organization]
menu:
  docs:
    parent: content-management
    weight: 120
toc: true
weight: 120
aliases: [/content/sections/]
---

A **Section** is a collection of pages that gets defined based on the
organization structure under the `content/` directory.

By default, all the **first-level** directories under `content/` form their own
sections (**root sections**) provided they constitute [Branch Bundles][branch bundles].
Directories which are just [Leaf Bundles][leaf bundles] do *not* form
their own sections, despite being first-level directories.

If a user needs to define a section `foo` at a deeper level, they need to create
a directory named `foo` with an `_index.md` file (see [Branch Bundles][branch bundles]
for more information).


{{% note %}}
A **section** cannot be defined or overridden by a front matter parameter -- it
is strictly derived from the content organization structure.
{{% /note %}}

## Nested Sections

The sections can be nested as deeply as you need.

```bash
content
└── blog        <-- Section, because first-level dir under content/
    ├── funny-cats
    │   ├── mypost.md
    │   └── kittens         <-- Section, because contains _index.md
    │       └── _index.md
    └── tech                <-- Section, because contains _index.md
        └── _index.md
```

**The important part to understand is, that to make the section tree fully navigational, at least the lower-most section needs a content file. (e.g. `_index.md`).**

{{% note %}}
When we talk about a **section** in correlation with template selection, it is
currently always the *root section* only (`/blog/funny-cats/mypost/ => blog`).

If you need a specific template for a sub-section, you need to adjust either the `type` or `layout` in front matter.
{{% /note %}}

## Example: Breadcrumb Navigation

With the available [section variables and methods](#section-page-variables-and-methods) you can build powerful navigation. One common example would be a partial to show Breadcrumb navigation:

{{< code file="layouts/partials/breadcrumb.html" download="breadcrumb.html" >}}
<nav aria-label="breadcrumb">
  <ol>
    {{ range .Ancestors.Reverse }}
      <li>
        <a href="{{ .Permalink }}">{{ .LinkTitle }}</a>
      </li>
    {{ end }}
    <li class="active">
      <a aria-current="page" href="{{ .Permalink }}">{{ .LinkTitle }}</a>
    </li>
  </ol>
</nav>
{{< /code >}}

## Section Page Variables and Methods

Also see [Page Variables](/variables/page/).

{{< readfile file="/content/en/readfiles/sectionvars.md" markdown="true" >}}

## Content Section Lists

Hugo will automatically create a page for each *root section* that lists all the content in that section. See the documentation on [section templates] for details on customizing the way these pages are rendered.

## Content *Section* vs Content *Type*

By default, everything created within a section will use the [content `type`][content type] that matches the *root section* name. For example, Hugo will assume that `posts/post-1.md` has a `posts` content `type`. If you are using an [archetype] for your `posts` section, Hugo will generate front matter according to what it finds in `archetypes/posts.md`.

[archetype]: /content-management/archetypes/
[content type]: /content-management/types/
[directory structure]: /getting-started/directory-structure/
[section templates]: /templates/section-templates/
[leaf bundles]: /content-management/page-bundles/#leaf-bundles
[branch bundles]: /content-management/page-bundles/#branch-bundles
