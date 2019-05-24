---
title: anchorize
description: Takes a string and sanitizes it the same way as Blackfriday does for markdown headers.
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

The template function uses the [`SanitizedAnchorName` logic from Blackfriday](https://github.com/russross/blackfriday#sanitized-anchor-names).
Since the same sanitizing logic is used as the markdown parser, you can determine the ID of a header for linking with anchor tags.

```
{{anchorize "This is a header"}} → "this-is-a-header"
{{anchorize "This is also          a header"}} → "this-is-also-a-header"
{{anchorize "main.go"}} → "main-go"
{{anchorize "Article 123"}} → "article-123"
{{anchorize "<- Let's try this, shall we?"}} → "let-s-try-this-shall-we"
{{anchorize "Hello, 世界"}} → "hello-世界"
```
