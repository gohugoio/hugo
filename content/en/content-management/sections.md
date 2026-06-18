---
title: Sections
description: Organize content into sections.

categories: []
keywords: []
aliases: [/content/sections/]
---

## Overview

{{% glossary-term "section" %}}

```tree
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

&nbsp;|Sections|Non-sections
:--|:-:|:-:
Directory names become URL segments|:heavy_check_mark:|:heavy_check_mark:
Have logical ancestors and descendants|:heavy_check_mark:|:x:
Have list pages|:heavy_check_mark:|:x:

With the file structure from the [example above](#overview):

1. The list page for the articles section includes all articles, regardless of directory structure; none of the subdirectories are sections.
1. The articles/2022 and articles/2023 directories do not have list pages; they are not sections.
1. The list page for the products section, by default, includes product-1 and product-2, but not their descendant pages. To include descendant pages, use the `RegularPagesRecursive` method instead of the `Pages` method in the _section_ template.
1. All directories in the products section have list pages; each directory is a section.

## Template selection

Hugo has a defined [lookup order][] to determine which template to use when rendering a page. The [lookup rules][] consider the top-level section name; subsection names are not considered when selecting a template.

With the file structure from the [example above](#overview):

Content directory|Section template
:--|:--
`content/products`|`layouts/products/section.html`
`content/products/product-1`|`layouts/products/section.html`
`content/products/product-1/benefits`|`layouts/products/section.html`

Content directory|Page template
:--|:--
`content/products`|`layouts/products/page.html`
`content/products/product-1`|`layouts/products/page.html`
`content/products/product-1/benefits`|`layouts/products/page.html`

If you need to use a different template for a subsection, specify `type` and/or `layout` in front matter.

## Ancestors and descendants

A section has one or more ancestors (including the home page), and zero or more descendants. With the file structure from the [example above](#overview):

```text
content/products/product-1/benefits/benefit-1.md
```

The content file (benefit-1.md) has four ancestors: benefits, product-1, products, and the home page. This logical relationship allows us to use the `.Parent` and `.Ancestors` methods to traverse the site structure.

For example, use the `.Ancestors` method to render breadcrumb navigation.

```go-html-template {file="layouts/_partials/breadcrumb.html"}
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
```

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

[lookup order]: /templates/lookup-order/
[lookup rules]: /templates/lookup-order/#lookup-rules
