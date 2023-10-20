---
title: urls.Anchorize
linkTitle: anchorize
description: Takes a string and sanitizes it the same way as the [`defaultMarkdownHandler`](/getting-started/configuration-markup#default-configuration) does for markdown headers.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [anchorize]
  returnType: string
  signatures: [urls.Anchorize INPUT]
relatedFunctions:
  - urls.Anchorize
  - urls.URLize
aliases: [/functions/anchorize]
---

If [Goldmark](/getting-started/configuration-markup#goldmark) is set as `defaultMarkdownHandler`, the sanitizing logic adheres to the setting [`markup.goldmark.parser.autoHeadingIDType`](/getting-started/configuration-markup#goldmark).

Since the `defaultMarkdownHandler` and this template function use the same sanitizing logic, you can use the latter to determine the ID of a header for linking with anchor tags.

```go-html-template
{{ anchorize "This is a header" }} → "this-is-a-header"
{{ anchorize "This is also    a header" }} → "this-is-also----a-header"
{{ anchorize "main.go" }} → "maingo"
{{ anchorize "Article 123" }} → "article-123"
{{ anchorize "<- Let's try this, shall we?" }} → "--lets-try-this-shall-we"
{{ anchorize "Hello, 世界" }} → "hello-世界"
```
