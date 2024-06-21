---
title: ByWeight
description: Returns the given menu with its entries sorted by weight, then by name, then by identifier.
categories: []
keywords: []
action:
  related: []
  returnType: navigation.Menu
  signatures: [MENU.ByWeight]
---

The `ByWeight` method returns the given menu with its entries sorted by [`weight`], then by `name`, then by `identifier`. This is the default sort order.

[`weight`]: /getting-started/glossary/#weight

Consider this menu definition:

{{< code-toggle file=hugo >}}
[[menus.main]]
identifier = 'about'
name = 'About'
pageRef = '/about'
weight = 20

[[menus.main]]
identifier = 'services'
name = 'Services'
pageRef = '/services'
weight = 10

[[menus.main]]
identifier = 'contact'
name = 'Contact'
pageRef = '/contact'
weight = 30
{{< /code-toggle >}}

To sort the entries by `weight`, then by `name`, then by `identifier`:

```go-html-template
<ul>
  {{ range .Site.Menus.main.ByWeight }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```

Hugo renders this to:

```html
<ul>
  <li><a href="/services/">Services</a></li>
  <li><a href="/about/">About</a></li>
  <li><a href="/contact">Contact</a></li>
</ul>
```

{{% note %}}
In the menu definition above, note that the `identifier` property is only required when two or more menu entries have the same name, or when localizing the name using translation tables.

[details]: /content-management/menus/#properties-front-matter
{{% /note %}}

You can also sort menu entries using the [`sort`] function. For example, to sort by `weight` in descending order:

```go-html-template
<ul>
  {{ range sort .Site.Menus.main "Weight" "desc" }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```

When using the sort function with menu entries, specify any of the following keys: `Identifier`, `Name`, `Parent`, `Post`, `Pre`, `Title`, `URL`, or `Weight`.

[`sort`]: /functions/collections/sort/
