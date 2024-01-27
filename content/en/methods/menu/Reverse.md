---
title: Reverse
description: Returns the given menu, reversing the sort order of its entries.
categories: []
keywords: []
action:
  related: []
  returnType: navigation.Menu
  signatures: [MENU.Reverse]
---

The `Reverse` method returns the given menu, reversing the sort order of its entries.

Consider this menu definition:

{{< code-toggle file=hugo >}}
[[menus.main]]
name = 'Services'
pageRef = '/services'
weight = 10

[[menus.main]]
name = 'About'
pageRef = '/about'
weight = 20

[[menus.main]]
name = 'Contact'
pageRef = '/contact'
weight = 30
{{< /code-toggle >}}

To sort the entries by name in descending order:

```go-html-template
<ul>
  {{ range .Site.Menus.main.ByName.Reverse }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```

Hugo renders this to:

```html
<ul>
  <li><a href="/services/">Services</a></li>
  <li><a href="/contact">Contact</a></li>
  <li><a href="/about/">About</a></li>
</ul>
```
