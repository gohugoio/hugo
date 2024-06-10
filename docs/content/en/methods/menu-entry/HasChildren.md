---
title: HasChildren
description: Reports whether the given menu entry has child menu entries.
categories: []
keywords: []
action:
  related:
    - methods/menu-entry/Children
  returnType: bool
  signatures: [MENUENTRY.HasChildren]
---

Use the `HasChildren` method when rendering a nested menu.

With this site configuration:

{{< code-toggle file=hugo >}}
[[menus.main]]
name = 'Products'
pageRef = '/product'
weight = 10

[[menus.main]]
name = 'Product 1'
pageRef = '/products/product-1'
parent = 'Products'
weight = 1

[[menus.main]]
name = 'Product 2'
pageRef = '/products/product-2'
parent = 'Products'
weight = 2
{{< /code-toggle >}}

And this template:

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li>
      <a href="{{ .URL }}">{{ .Name }}</a>
      {{ if .HasChildren }}
        <ul>
          {{ range .Children }}
            <li><a href="{{ .URL }}">{{ .Name }}</a></li>
          {{ end }}
        </ul>
      {{ end }}
    </li>
  {{ end }}
</ul>
```

Hugo renders this HTML:

```html
<ul>
  <li>
    <a href="/products/">Products</a>
    <ul>
      <li><a href="/products/product-1/">Product 1</a></li>
      <li><a href="/products/product-2/">Product 2</a></li>
    </ul>
  </li>
</ul>
```
