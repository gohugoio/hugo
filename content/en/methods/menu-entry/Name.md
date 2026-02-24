---
title: Name
description: Returns the `name` property of the given menu entry.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [MENUENTRY.Name]
---

If you define the menu entry [automatically], the `Name` method returns the page's [`LinkTitle`], falling back to its [`Title`].

If you define the menu entry in [front matter] or in your [project configuration], the `Name` method returns the `name` property of the given menu entry. If the `name` is not defined, and the menu entry resolves to a page, the `Name` returns the page [`LinkTitle`], falling back to its [`Title`].

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```

[`LinkTitle`]: /methods/page/linktitle/
[`Title`]: /methods/page/title/
[automatically]: /content-management/menus/#define-automatically
[front matter]: /content-management/menus/#define-in-front-matter
[project configuration]: /content-management/menus/#define-in-project-configuration
