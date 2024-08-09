---
title: Page
description: Returns the Page object associated with the given menu entry.
categories: []
keywords: []
action:
  related: []
  returnType: page.Page
  signatures: [MENUENTRY.Page]
---

Regardless of how you [define menu entries], an entry associated with a page has access to its [methods].

In this menu definition, the first two entries are associated with a page, the last entry is not:

{{< code-toggle file=hugo >}}
[[menus.main]]
pageRef = '/about'
weight = 10

[[menus.main]]
pageRef = '/contact'
weight = 20

[[menus.main]]
name = 'Hugo'
url = 'https://gohugo.io'
weight = 30
{{< /code-toggle >}}

In this example, if the menu entry is associated with a page, we use page's [`RelPermalink`] and [`LinkTitle`] when rendering the anchor element.

If the entry is not associated with a page, we use its `url` and `name` properties.

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    {{ with .Page }}
      <li><a href="{{ .RelPermalink }}">{{ .Title }}</a></li>
    {{ else }}
      <li><a href="{{ .URL }}">{{ .Name }}</a></li>
    {{ end }}
  {{ end }}
</ul>
```

See the [menu templates] section for more information.

[`LinkTitle`]: /methods/page/linktitle/
[`RelPermalink`]: /methods/page/relpermalink/
[define menu entries]: /content-management/menus/
[menu templates]: /templates/menu/#page-references
[methods]: /methods/page/
