---
title: urls.Anchorize
description: Returns a sanitized string to use as an id attribute within an anchor element.
categories: []
keywords: []
action:
  aliases: [anchorize]
  related:
    - functions/urls/URLize
  returnType: string
  signatures: [urls.Anchorize INPUT]
aliases: [/functions/anchorize]
---

If [Goldmark](/getting-started/configuration-markup#goldmark) is set as `defaultMarkdownHandler`, the sanitizing logic adheres to the setting [`markup.goldmark.parser.autoHeadingIDType`](/getting-started/configuration-markup#goldmark).

Since the `defaultMarkdownHandler` and this template function use the same sanitizing logic, you can use the latter to determine the ID of a header for linking with anchor tags.

```go-html-template
{{ anchorize "This is a header" }} → this-is-a-header
{{ anchorize "This is also    a header" }} → this-is-also----a-header
{{ anchorize "main.go" }} → maingo
{{ anchorize "Article 123" }} → article-123
{{ anchorize "<- Let's try this, shall we?" }} → --lets-try-this-shall-we
{{ anchorize "Hello, 世界" }} → hello-世界
```
