---
title: Links and Cross References
description: Shortcodes for creating links to documents.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-31
categories: [content management]
keywords: ["cross references","references", "anchors", "urls"]
menu:
  docs:
    parent: "content-management"
    weight: 100
weight: 100	#rem
aliases: [/extras/crossreferences/]
toc: true
---

The `ref` and `relref` shortcodes display the absolute and relative permalinks to a document, respectively.

## Use `ref` and `relref`

```go-html-template
{{</* ref "document" */>}}
{{</* ref "document#anchor" */>}}
{{</* ref "document.md" */>}}
{{</* ref "document.md#anchor" */>}}
{{</* ref "#anchor" */>}}
{{</* ref "/blog/my-post" */>}}
{{</* ref "/blog/my-post.md" */>}}
{{</* relref "document" */>}}
{{</* relref "document.md" */>}}
{{</* relref "#anchor" */>}}
{{</* relref "/blog/my-post.md" */>}}
```

To generate a hyperlink using `ref` or `relref` in markdown:

```md
[About]({{</* ref "/page/about" */>}} "About Us")
```

The `ref` and `relref` shortcodes require a single parameter: the path to a content document, with or without a file extension, with or without an anchor.

**Paths without a leading `/` are first resolved relative to the current page, then to the remainder of the site.

Hugo emits an error or warning if a document cannot be uniquely resolved. The error behavior is configurable; see below.

### Link to another language version

To link to another language version of a document, use this syntax:

```go-html-template
{{</* relref path="document.md" lang="ja" */>}}
```

### Get another Output Format

To link to another Output Format of a document, use this syntax:

```go-html-template
{{</* relref path="document.md" outputFormat="rss" */>}}
```

### Heading IDs

When using Markdown document types, Hugo generates element IDs for every heading on a page. For example:

```md
## Reference
```

produces this HTML:

```html
<h2 id="reference">Reference</h2>
```

Get the permalink to a heading by appending the ID to the path when using the `ref` or `relref` shortcodes:

```go-html-template
{{</* ref "document.md#reference" */>}}
{{</* relref "document.md#reference" */>}}
```

Generate a custom heading ID by including an attribute. For example:

```md
## Reference A {#foo}
## Reference B {id="bar"}
```

produces this HTML:

```html
<h2 id="foo">Reference A</h2>
<h2 id="bar">Reference B</h2>
```

Hugo will generate unique element IDs if the same heading appears more than once on a page. For example:

```md
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

The behavior can, since Hugo 0.45, be configured in `config.toml`:

refLinksErrorLevel ("ERROR")
: When using `ref` or `relref` to resolve page links and a link cannot resolved, it will be logged with this log level. Valid values are `ERROR` (default) or `WARNING`. Any `ERROR` will fail the build (`exit -1`).

refLinksNotFoundURL
: URL to be used as a placeholder when a page reference cannot be found in `ref` or `relref`. Is used as-is.


[lists]: /templates/lists/
[output formats]: /templates/output-formats/
[shortcode]: /content-management/shortcodes/
