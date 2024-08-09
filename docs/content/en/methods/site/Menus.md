---
title: Menus
description: Returns a collection of menu objects for the given site.
categories: []
keywords: []
action:
  related:
    - methods/page/IsMenuCurrent
    - methods/page/HasMenuCurrent
  returnType: navigation.Menus
  signatures: [SITE.Menus]
---

The `Menus` method on a `Site` object returns a collection of menus, where each menu contains one or more entries, either flat or nested. Each entry points to a page within the site, or to an external resource.

{{% note %}}
Menus can be defined and localized in several ways. Please see the [menus] section for a complete explanation and examples.

[menus]: /content-management/menus/
{{% /note %}}

A site can have multiple menus. For example, a main menu and a footer menu:

{{< code-toggle file=hugo >}}
[[menus.main]]
name = 'Home'
pageRef = '/'
weight = 10

[[menus.main]]
name = 'Books'
pageRef = '/books'
weight = 20

[[menus.main]]
name = 'Films'
pageRef = '/films'
weight = 30

[[menus.footer]]
name = 'Legal'
pageRef = '/legal'
weight = 10

[[menus.footer]]
name = 'Privacy'
pageRef = '/privacy'
weight = 20
{{< /code-toggle >}}

This template renders the main menu:

```go-html-template
{{ with site.Menus.main }}
  <nav class="menu">
    {{ range . }}
      {{ if $.IsMenuCurrent .Menu . }}
        <a class="active" aria-current="page" href="{{ .URL }}">{{ .Name }}</a>
      {{ else }}
        <a href="{{ .URL }}">{{ .Name }}</a>
      {{ end }}
    {{ end }}
  </nav>
{{ end }}
```

When viewing the home page, the result is:

```html
<nav class="menu">
  <a class="active" aria-current="page" href="/">Home</a>
  <a href="/books/">Books</a>
  <a href="/films/">Films</a>
</nav>
```

When viewing the "books" page, the result is:

```html
<nav class="menu">
  <a href="/">Home</a>
  <a class="active" aria-current="page" href="/books/">Books</a>
  <a href="/films/">Films</a>
</nav>
```

You will typically render a menu using a partial template. As the active menu entry will be different on each page, use the [`partial`] function to call the template. Do not use the [`partialCached`] function.

The example above is simplistic. Please see the [menu templates] section for more information.

[menu templates]: /templates/menu/

[`partial`]: /functions/partials/include/
[`partialCached`]: /functions/partials/includecached/
