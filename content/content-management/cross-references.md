---
title: Links and Cross References
description: Hugo makes it easy to link documents together.
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


 The `ref` and `relref` shortcodes link documents together, both of which are [built-in Hugo shortcodes][]. These shortcodes are also used to provide links to headings inside of your content, whether across documents or within a document. The only difference between `ref` and `relref` is whether the resulting URL is absolute (`http://1.com/about/`) or relative (`/about/`), respectively.

## Use `ref` and `relref`

```
{{</* ref "document.md" */>}}
{{</* ref "#anchor" */>}}
{{</* ref "document.md#anchor" */>}}
{{</* relref "document.md" */>}}
{{</* relref "#anchor" */>}}
{{</* relref "document.md#anchor" */>}}
```

The single parameter to `ref` is a string with a content `documentname` (e.g., `about.md`) with or without an appended in-document `anchor` (`#who`) without spaces.

### Document Names

The `documentname` is the name of a document, including the format extension; this may be just the filename, or the relative path from the `content/` directory. With a document `content/blog/post.md`, either format will produce the same result:

```
{{</* relref "blog/post.md" */>}} => `/blog/post/`
{{</* relref "post.md" */>}} => `/blog/post/`
```

If you have the same filename used across multiple sections, you should only use the relative path format; otherwise, the behavior will be `undefined`. This is best illustrated with an example `content` directory:

```
.
└── content
    ├── events
    │   └── my-birthday.md
    ├── galleries
    │   └── my-birthday.md
    ├── meta
    │   └── my-article.md
    └── posts
        └── my-birthday.md
```

To be sure to get the correct reference in this case, use the full path: 

{{< code file="content/meta/my-article.md" copy="false" >}}
{{</* relref "events/my-birthday.md" */>}} => /events/my-birthday/
{{< /code >}}

### With Multiple Output Formats

If the page exists in multiple [output formats][], `ref` or `relref` can be used with a output format name:

```
 [Neat]({{</* ref "blog/neat.md" "amp" */>}})
```

### Anchors

When an `anchor` is provided by itself, the current page’s unique identifier will be appended; when an `anchor` is provided appended to `documentname`, the found page's unique identifier will be appended:

```
{{</* relref "#anchors" */>}} => #anchors:9decaf7
{{</* relref "about-hugo/hugo-features.md#content" */>}} => /blog/post/#who:badcafe
```

The above examples render as follows for this very page as well as a reference to the "Content" heading in the Hugo docs features pageyoursite

```
{{</* relref "#who" */>}} => #who:9decaf7
{{</* relref "blog/post.md#who" */>}} => /blog/post/#who:badcafe
```

More information about document unique identifiers and headings can be found [below]({{< ref "#hugo-heading-anchors" >}}).

### Examples

* `{{</* ref "blog/post.md" */>}}` => `https://example.com/blog/post/`
* `{{</* ref "post.md#tldr" */>}}` => `https://example.com/blog/post/#tldr:caffebad`
* `{{</* relref "post.md" */>}}` => `/blog/post/`
* `{{</* relref "blog/post.md#tldr" */>}}` => `/blog/post/#tldr:caffebad`
* `{{</* ref "#tldr" */>}}` => `#tldr:badcaffe`
* `{{</* relref "#tldr" */>}}` => `#tldr:badcaffe`

## Hugo Heading Anchors

When using Markdown document types, Hugo generates heading anchors automatically. The generated anchor for this section is `hugo-heading-anchors`. Because the heading anchors are generated automatically, Hugo takes some effort to ensure that heading anchors are unique both inside a document and across the entire site.

Ensuring heading uniqueness across the site is accomplished with a unique identifier for each document based on its path. Unless a document is renamed or moved between sections *in the filesystem*, the unique identifier for the document will not change: `blog/post.md` will always have a unique identifier of `81df004c333b392d34a49fd3a91ba720`.

`ref` and `relref` were added so you can make these reference links without having to know the document’s unique identifier. (The links in document tables of contents are automatically up-to-date with this value.)

```
{{</* relref "content-management/cross-references.md#hugo-heading-anchors" */>}}
/content-management/cross-references/#hugo-heading-anchors:77cd9ea530577debf4ce0f28c8dca242
```

### Manually Specifying Anchors

For Markdown content files, if the `headerIds` [Blackfriday extension][bfext] is
enabled (which it is by default), user can manually specify the anchor for any
heading.

Few examples:

```
## Alpha 101 {#alpha}

## Version 1.0 {#version-1-dot-0}
```

[built-in Hugo shortcodes]: /content-management/shortcodes/#using-the-built-in-shortcodes
[lists]: /templates/lists/
[output formats]: /templates/output-formats/
[shortcode]: /content-management/shortcodes/
[bfext]: /content-management/formats/#blackfriday-extensions
