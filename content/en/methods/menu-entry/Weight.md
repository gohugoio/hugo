---
title: Weight
description: Returns the `weight` property of the given menu entry.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [MENUENTRY.Weight]
---

If you define the menu entry [automatically], the `Weight` method returns the page's [`Weight`].

If you define the menu entry in [front matter] or in your [project configuration], the `Weight` method returns the `weight` property, falling back to the page's `Weight`.

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

[`Weight`]: /methods/page/weight/
[automatically]: /content-management/menus/#define-automatically
[front matter]: /content-management/menus/#define-in-front-matter
[project configuration]: /content-management/menus/#define-in-project-configuration
