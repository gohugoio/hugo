---
title: PageRef
description: Returns the `pageRef` property of the given menu entry.
categories: []
keywords: []
action:
  related:
    - /methods/menu-entry/URL
  returnType: string
  signatures: [MENUENTRY.PageRef]
toc: true
---

The use case for this method is rare.

In almost also scenarios you should use the [`URL`] method instead.

## Explanation

If you specify a `pageRef` property when [defining a menu entry] in your site configuration, Hugo looks for a matching page when rendering the entry.

If a matching page is found:

- The [`URL`] method returns the page's relative permalink
- The [`Page`] method returns the corresponding `Page` object
- The [`HasMenuCurrent`] and [`IsMenuCurrent`] methods on a `Page` object return the expected values

If a matching page is not found:

- The [`URL`] method returns the entry's `url` property if set, else an empty string
- The [`Page`] method returns nil
- The [`HasMenuCurrent`] and [`IsMenuCurrent`] methods on a `Page` object return `false`

{{% note %}}
In almost also scenarios you should use the [`URL`] method instead.

[`URL`]: /methods/menu-entry/url/
{{% /note %}}

[defining a menu entry]: /content-management/menus/#define-in-site-configuration
[`Page`]: /methods/menu-entry/page/
[`URL`]: /methods/menu-entry/url/
[`IsMenuCurrent`]: /methods/page/ismenucurrent/
[`HasMenuCurrent`]: /methods/page/hasmenucurrent/
[`RelPermalink`]: /methods/page/relpermalink/

## Example

This example is contrived.

{{% note %}}
In almost also scenarios you should use the [`URL`] method instead.

[`URL`]: /methods/menu-entry/url/
{{% /note %}}


Consider this content structure:

```text
content/
├── products.md
└── _index.md
```

And this menu definition:

{{< code-toggle file=hugo >}}
[[menus.main]]
name = 'Products'
pageRef = '/products'
weight = 10
[[menus.main]]
name = 'Services'
pageRef = '/services'
weight = 20
{{< /code-toggle >}}

With this template code:

{{< code file=layouts/partials/menu.html >}}
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
{{< /code >}}

Hugo render this HTML:

```html
<ul>
  <li><a href="/products/">Products</a></li>
  <li><a href="">Services</a></li>
</ul>
```

In the above note that the `href` attribute of the second `anchor` element is blank because Hugo was unable to find the "services" page.

With this template code:


{{< code file=layouts/partials/menu.html >}}
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ or .URL .PageRef }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
{{< /code >}}

Hugo renders this HTML:

```html
<ul>
  <li><a href="/products/">Products</a></li>
  <li><a href="/services">Services</a></li>
</ul>
```

In the above note that Hugo populates the `href` attribute of the second `anchor` element with the `pageRef` property as defined in the site configuration because the template code falls back to the `PageRef` method.
