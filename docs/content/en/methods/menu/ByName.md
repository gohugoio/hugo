---
title: ByName
description: Returns the given menu with its entries sorted by name.
categories: []
keywords: []
action:
  related: []
  returnType: navigation.Menu
  signatures: [MENU.ByName]
---

The `Sort` method returns the given menu with its entries sorted by `name`.

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

To sort the entries by `name`:

```go-html-template
<ul>
  {{ range .Site.Menus.main.ByName }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```

Hugo renders this to:

```html
<ul>
  <li><a href="/about/">About</a></li>
  <li><a href="/contact">Contact</a></li>
  <li><a href="/services/">Services</a></li>
</ul>
```

You can also sort menu entries using the [`sort`] function. For example, to sort by `name` in descending order:

```go-html-template
<ul>
  {{ range sort .Site.Menus.main "Name" "desc" }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```

When using the sort function with menu entries, specify any of the following keys: `Identifier`, `Name`, `Parent`, `Post`, `Pre`, `Title`, `URL`, or `Weight`.

[`sort`]: /functions/collections/sort/
