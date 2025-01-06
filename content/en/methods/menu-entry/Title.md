---
title: Title
description: Returns the `title` property of the given menu entry.  
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [MENUENTRY.Title]
---

The `Title` method returns the `title` property of the given menu entry. If the `title` is not defined, and the menu entry resolves to a page, the `Title`  returns the page [`Title`].

[`Title`]: /methods/page/title/

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ .URL }}" title="{{ .Title }}>{{ .Name }}</a></li>
  {{ end }}
</ul>
```
