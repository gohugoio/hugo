---
title: Name
description: Returns the `name` property of the given menu entry.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [MENUENTRY.Name]
---

If you define the menu entry [automatically], the `Name` method returns the page’s [`LinkTitle`], falling back to its [`Title`].

If you define the menu entry [in front matter] or [in site configuration], the `Name` method returns the `name` property, falling back to the page’s `LinkTitle`, then to its `Title`.

[`LinkTitle`]: /methods/page/linktitle/
[`Title`]: /methods/page/title/
[automatically]: /content-management/menus/#define-automatically
[in front matter]: /content-management/menus/#define-in-front-matter
[in site configuration]: /content-management/menus/#define-in-site-configuration

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```
