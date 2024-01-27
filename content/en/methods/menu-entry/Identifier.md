---
title: Identifier
description: Returns the `identifier` property of the given menu entry. 
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [MENUENTRY.Identifier]
---

The `Identifier` method returns the `identifier` property of the menu entry. If you define the menu entry [automatically], it returns the page's section.

[automatically]: /content-management/menus/#define-automatically

{{< code-toggle file=hugo >}}
[[menus.main]]
identifier = 'about'
name = 'About'
pageRef = '/about'
weight = 10

[[menus.main]]
identifier = 'contact'
name = 'Contact'
pageRef = '/contact'
weight = 20
{{< /code-toggle >}}

This example uses the `Identifier` method when querying the translation table on a multilingual site, falling back the `name` property if a matching key in the translation table does not exist:

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ .URL }}">{{ or (T .Identifier) .Name }}</a></li>
  {{ end }}
</ul>
```

{{% note %}}
In the menu definition above, note that the `identifier` property is only required when two or more menu entries have the same name, or when localizing the name using translation tables.

[details]: /content-management/menus/#properties-front-matter
{{% /note %}}
