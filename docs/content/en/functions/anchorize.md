---
title: anchorize
description: Takes a string and sanitizes it the same way as the [`defaultMarkdownHandler`](https://gohugo.io/getting-started/configuration-markup#configure-markup) does for markdown headers.
date: 2018-10-13
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [markdown,strings]
signature: ["anchorize INPUT"]
hugoversion: "0.39"
workson: []
relatedfuncs: [humanize]
---

If [Goldmark](https://gohugo.io/getting-started/configuration-markup#goldmark) is set as `defaultMarkdownHandler`, the sanitizing logic adheres to the setting [`markup.goldmark.parser.autoHeadingIDType`](https://gohugo.io/getting-started/configuration-markup#goldmark). If [Blackfriday](https://gohugo.io/getting-started/configuration-markup#blackfriday) is set as `defaultMarkdownHandler`, this template function uses the [`SanitizedAnchorName` logic from Blackfriday](https://github.com/russross/blackfriday#sanitized-anchor-names) (the same applies when `markup.goldmark.parser.autoHeadingIDType` is set to `blackfriday`).

Since the `defaultMarkdownHandler` and this template function use the same sanitizing logic, you can use the latter to determine the ID of a header for linking with anchor tags.

```
{{anchorize "This is a header"}} → "this-is-a-header"
{{anchorize "This is also          a header"}} → "this-is-also-a-header"
{{anchorize "main.go"}} → "main-go"
{{anchorize "Article 123"}} → "article-123"
{{anchorize "<- Let's try this, shall we?"}} → "let-s-try-this-shall-we"
{{anchorize "Hello, 世界"}} → "hello-世界"
```
