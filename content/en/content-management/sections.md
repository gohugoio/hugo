---
title: Sections
description: Organize content into sections.

categories: [content management]
keywords: [lists,sections,content types,organization]
menu:
  docs:
    parent: content-management
    weight: 120
weight: 120
toc: true
aliases: [/content/sections/]
---

## Overview

A section is a top-level content directory, or any content directory with an&nbsp;_index.md file. A content directory with an _index.md file is also known as a [branch bundle](/getting-started/glossary/#branch-bundle). Section templates receive one or more page [collections](/getting-started/glossary/#collection) in [context](/getting-started/glossary/#context).

{{% note %}}
Although top-level directories without _index.md files are sections, we recommend creating _index.md files in _all_ sections.
{{% /note %}}

A typical site consists of one or more sections. For example:

```text
content/
├── articles/             <-- section (top-level directory)
│   ├── 2022/
│   │   ├── article-1/
│   │   │   ├── cover.jpg
│   │   │   └── index.md
│   │   └── article-2.md
│   └── 2023/
│       ├── article-3.md
│       └── article-4.md
├── products/             <-- section (top-level directory)
│   ├── product-1/        <-- section (has _index.md file)
│   │   ├── benefits/     <-- section (has _index.md file)
│   │   │   ├── _index.md
│   │   │   ├── benefit-1.md
│   │   │   └── benefit-2.md
│   │   ├── features/     <-- section (has _index.md file)
│   │   │   ├── _index.md
│   │   │   ├── feature-1.md
│   │   │   └── feature-2.md
│   │   └── _index.md
│   └── product-2/        <-- section (has _index.md file)
│       ├── benefits/     <-- section (has _index.md file)
│       │   ├── _index.md
│       │   ├── benefit-1.md
│       │   └── benefit-2.md
│       ├── features/     <-- section (has _index.md file)
│       │   ├── _index.md
│       │   ├── feature-1.md
│       │   └── feature-2.md
│       └── _index.md
├── _index.md
└── about.md
```

The example above has two top-level sections: articles and products. None of the directories under articles are sections, while all of the directories under products are sections. A section within a section is a known as a nested section or subsection.

## Explanation

Sections and non-sections behave differently.

||Sections|Non-sections
:--|:-:|:-:
Directory names become URL segments|:heavy_check_mark:|:heavy_check_mark:
Have logical ancestors and descendants|:heavy_check_mark:|:x:
Have list pages|:heavy_check_mark:|:x:

With the file structure from the [example above](#overview):

1. The list page for the articles section includes all articles, regardless of directory structure; none of the subdirectories are sections.

1. The articles/2022 and articles/2023 directories do not have list pages; they are not sections.

1. The list page for the products section, by default, includes product-1 and product-2, but not their descendant pages. To include descendant pages, use the `RegularPagesRecursive` method instead of the `Pages` method in the list template.

[`Pages`]: /methods/page/pages/
[`RegularPagesRecursive`]: /methods/page/regularpagesrecursive/

1. All directories in the products section have list pages; each directory is a section.

## Template selection

Hugo has a defined [lookup order] to determine which template to use when rendering a page. The [lookup rules] consider the top-level section name; subsection names are not considered when selecting a template.

With the file structure from the [example above](#overview):

Content directory|Section template
:--|:--
content/products|layouts/products/list.html
content/products/product-1|layouts/products/list.html
content/products/product-1/benefits|layouts/products/list.html

Content directory|Single template
:--|:--
content/products|layouts/products/single.html
content/products/product-1|layouts/products/single.html
content/products/product-1/benefits|layouts/products/single.html

If you need to use a different template for a subsection, specify `type` and/or `layout` in front matter.

[lookup rules]: /templates/lookup-order/#lookup-rules
[lookup order]: /templates/lookup-order/

## Ancestors and descendants

A section has one or more ancestors (including the home page), and zero or more descendants. With the file structure from the [example above](#overview):

```text
content/products/product-1/benefits/benefit-1.md
```

The content file (benefit-1.md) has four ancestors: benefits, product-1, products, and the home page. This logical relationship allows us to use the `.Parent` and `.Ancestors` methods to traverse the site structure.

For example, use the `.Ancestors` method to render breadcrumb navigation.

{{< code file=layouts/partials/breadcrumb.html >}}
<nav aria-label="breadcrumb" class="breadcrumb">
  <ol>
    {{ range .Ancestors.Reverse }}
      <li>
        <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
      </li>
    {{ end }}
    <li class="active">
      <a aria-current="page" href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
    </li>
  </ol>
</nav>
{{< /code >}}

With this CSS:

```css
.breadcrumb ol {
  padding-left: 0;
}

.breadcrumb li {
  display: inline;
}

.breadcrumb li:not(:last-child)::after {
  content: "»";
}
```

Hugo renders this, where each breadcrumb is a link to the corresponding page:

```text
Home » Products » Product 1 » Benefits » Benefit 1
```

[archetype]: /content-management/archetypes/
[content type]: /content-management/types/
[directory structure]: /getting-started/directory-structure/
[section templates]: /templates/types/#section
[leaf bundles]: /content-management/page-bundles/#leaf-bundles
[branch bundles]: /content-management/page-bundles/#branch-bundles
