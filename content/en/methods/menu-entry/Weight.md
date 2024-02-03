---
title: Weight
description:  Returns the `weight` property of the given menu entry.   
categories: []
keywords: []
action:
  related: []
  returnType: int
  signatures: [MENUENTRY.Weight]
---

If you define the menu entry [automatically], the `Weight` method returns the page’s [`Weight`].

If you define the menu entry [in front matter] or [in site configuration], the `Weight` method returns the `weight` property, falling back to the page’s `Weight`.

[`Weight`]: /methods/page/weight/
[automatically]: /content-management/menus/#define-automatically
[in front matter]: /content-management/menus/#define-in-front-matter
[in site configuration]: /content-management/menus/#define-in-site-configuration

In this contrived example, we limit the number of menu entries based on weight:

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    {{ if le .Weight 42 }}
      <li><a href="{{ .URL }}">{{ .Name }}</a></li>
    {{ end }}
  {{ end }}
</ul>
```
