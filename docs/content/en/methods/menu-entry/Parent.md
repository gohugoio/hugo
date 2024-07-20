---
title: Parent
description: Returns the `parent` property of the given menu entry.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [MENUENTRY.Parent]
---

With this menu definition:

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

This template renders the nested menu, listing the `parent` property next each of the child entries:

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li>
      <a href="{{ .URL }}">{{ .Name }}</a>
      {{ if .HasChildren }}
        <ul>
          {{ range .Children }}
            <li><a href="{{ .URL }}">{{ .Name }}</a> ({{ .Parent  }})</li>
          {{ end }}
        </ul>
      {{ end }}
    </li>
  {{ end }}
</ul>
```
