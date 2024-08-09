---
title: Links and cross references
description: Shortcodes for creating links to documents.
categories: [content management]
keywords: [cross references,references,anchors,urls]
menu:
  docs:
    parent: content-management
    weight: 170
weight: 170
toc: true
aliases: [/extras/crossreferences/]
---

The `ref` and `relref` shortcodes display the absolute and relative permalinks to a document, respectively.

## Use of `ref` and `relref`

The `ref` and `relref` shortcodes require a single argument: the path to a content document, with or without a file extension, with or without an anchor. Paths without a leading `/` are first resolved relative to the current page, then to the remainder of the site.

```text
.
└── content
    ├── about
    |   ├── _index.md
    |   └── credits.md
    ├── pages
    |   ├── document1.md
    |   └── document2.md    // has anchor #anchor
    ├── products
    |   └── index.md
    └── blog
        └── my-post.md
```

The pages can be referenced as follows:

```text
{{</* ref "document2" */>}}             <-- From pages/document1.md, relative path
{{</* ref "document2#anchor" */>}}      
{{</* ref "document2.md" */>}}          
{{</* ref "document2.md#anchor" */>}}   
{{</* ref "#anchor" */>}}               <-- From pages/document2.md
{{</* ref "/blog/my-post" */>}}         <-- From anywhere, absolute path
{{</* ref "/blog/my-post.md" */>}}
{{</* relref "document" */>}}
{{</* relref "document.md" */>}}
{{</* relref "#anchor" */>}}
{{</* relref "/blog/my-post.md" */>}}
```

index.md can be reference either by its path or by its containing folder without the ending `/`. \_index.md can be referenced only by its containing folder:

```text
{{</* ref "/about" */>}}             <-- References /about/_index.md
{{</* ref "/about/_index" */>}}      <-- Raises REF_NOT_FOUND error
{{</* ref "/about/credits.md" */>}}  <-- References /about/credits.md

{{</* ref "/products" */>}}          <-- References /products/index.md
{{</* ref "/products/index" */>}}    <-- References /products/index.md
```

To generate a hyperlink using `ref` or `relref` in Markdown:

```text
[About]({{</* ref "/about" */>}} "About Us")
```

Hugo emits an error or warning if a document cannot be uniquely resolved. The error behavior is configurable; see below.

### Link to another language version

Using `ref` or `relref` without specifying a language, will make the reference resolve to the language of the current content page.

To link to another language version of a document, use this syntax:

```text
{{</* relref path="document.md" lang="ja" */>}}
```

### Get another output format

To link to another Output Format of a document, use this syntax:

```text
{{</* relref path="document.md" outputFormat="rss" */>}}
```

### Heading IDs

When using Markdown document types, Hugo generates element IDs for every heading on a page. For example:

```text
## Reference
```

produces this HTML:

```html
<h2 id="reference">Reference</h2>
```

Get the permalink to a heading by appending the ID to the path when using the `ref` or `relref` shortcodes:

```text
{{</* ref "document.md#reference" */>}}
{{</* relref "document.md#reference" */>}}
```

Generate a custom heading ID by including an attribute. For example:

```text
## Reference A {#foo}
## Reference B {id="bar"}
```

produces this HTML:

```html
<h2 id="foo">Reference A</h2>
<h2 id="bar">Reference B</h2>
```

Hugo will generate unique element IDs if the same heading appears more than once on a page. For example:

```text
## Reference
## Reference
## Reference
```

produces this HTML:

```html
<h2 id="reference">Reference</h2>
<h2 id="reference-1">Reference</h2>
<h2 id="reference-2">Reference</h2>
```

## Ref and RelRef Configuration

The behavior can be configured in `hugo.toml`:

refLinksErrorLevel ("ERROR")
: When using `ref` or `relref` to resolve page links and a link cannot resolved, it will be logged with this log level. Valid values are `ERROR` (default) or `WARNING`. Any `ERROR` will fail the build (`exit -1`).

refLinksNotFoundURL
: URL to be used as a placeholder when a page reference cannot be found in `ref` or `relref`. Is used as-is.
