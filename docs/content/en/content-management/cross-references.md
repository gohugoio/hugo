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


The `ref` and `relref` shortcode resolves the absolute or relative permalink given a path to a document.

## Use `ref` and `relref`

```go-html-template
{{</* ref "document.md" */>}}
{{</* ref "#anchor" */>}}
{{</* ref "document.md#anchor" */>}}
{{</* ref "/blog/my-post" */>}}
{{</* ref "/blog/my-post.md" */>}}
{{</* relref "document.md" */>}}
{{</* relref "#anchor" */>}}
{{</* relref "document.md#anchor" */>}}
```

The single parameter to `ref` is a string with a content `documentname` (e.g., `about.md`) with or without an appended in-document `anchor` (`#who`) without spaces. Hugo is flexible in how we search for documents, so the file suffix may be omitted.

**Paths without a leading `/` will first  be tried resolved relative to the current page.**

You will get an error if your document could not be uniquely resolved. The error behaviour can be configured, see below.

### Link to another language version

Link to another language version of a document, you need to use this syntax:

```go-html-template
{{</* relref path="document.md" lang="ja" */>}}
```

### Get another Output Format

To link to a given Output Format of a document, you can use this syntax:

```go-html-template
{{</* relref path="document.md" outputFormat="rss" */>}}
```

### Anchors

When an `anchor` is provided by itself, the current page’s unique identifier will be appended; when an `anchor` is provided appended to `documentname`, the found page's unique identifier will be appended:

```go-html-template
{{</* relref "#anchors" */>}} => #anchors:9decaf7
```

The above examples render as follows for this very page as well as a reference to the "Content" heading in the Hugo docs features pageyoursite

```go-html-template
{{</* relref "#who" */>}} => #who:9decaf7
{{</* relref "/blog/post.md#who" */>}} => /blog/post/#who:badcafe
```

More information about document unique identifiers and headings can be found [below]({{< ref "#hugo-heading-anchors" >}}).

## Hugo Heading Anchors

When using Markdown document types, Hugo generates heading anchors automatically. The generated anchor for this section is `hugo-heading-anchors`. Because the heading anchors are generated automatically, Hugo takes some effort to ensure that heading anchors are unique both inside a document and across the entire site.

Ensuring heading uniqueness across the site is accomplished with a unique identifier for each document based on its path. Unless a document is renamed or moved between sections *in the filesystem*, the unique identifier for the document will not change: `blog/post.md` will always have a unique identifier of `81df004c333b392d34a49fd3a91ba720`.

`ref` and `relref` were added so you can make these reference links without having to know the document’s unique identifier. (The links in document tables of contents are automatically up-to-date with this value.)

```
{{</* relref "content-management/cross-references.md#hugo-heading-anchors" */>}}
/content-management/cross-references/#hugo-heading-anchors:77cd9ea530577debf4ce0f28c8dca242
```

## Ref and RelRef Configuration

The behaviour can, since Hugo 0.45, be configured in `config.toml`:

refLinksErrorLevel ("ERROR") 
: When using `ref` or `relref` to resolve page links and a link cannot resolved, it will be logged with this log level. Valid values are `ERROR` (default) or `WARNING`. Any `ERROR` will fail the build (`exit -1`).

refLinksNotFoundURL
: URL to be used as a placeholder when a page reference cannot be found in `ref` or `relref`. Is used as-is.


[lists]: /templates/lists/
[output formats]: /templates/output-formats/
[shortcode]: /content-management/shortcodes/
[bfext]: /content-management/formats/#blackfriday-extensions
