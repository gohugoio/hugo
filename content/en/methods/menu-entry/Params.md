---
title: Params
description: Returns the `params` property of the given menu entry.
categories: []
keywords: []
action:
  related: []
  returnType: maps.Params
  signatures: [MENUENTRY.Params]
---

When you define menu entries [in site configuration] or [in front matter], you can include a `params` key to attach additional information to the entry. For example:

{{< code-toggle file=hugo >}}
[[menus.main]]
name = 'About'
pageRef = '/about'
weight = 10

[[menus.main]]
name = 'Contact'
pageRef = '/contact'
weight = 20

[[menus.main]]
name = 'Hugo'
url = 'https://gohugo.io'
weight = 30
[menus.main.params]
  rel = 'external'
{{< /code-toggle >}}

With this template:


```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li>
      <a href="{{ .URL }}" {{ with .Params.rel }}rel="{{ . }}"{{ end }}>
        {{ .Name }}
      </a>
    </li>
  {{ end }}
</ul>
```

Hugo renders:

```html
<ul>
  <li><a href="/about/">About</a></li>
  <li><a href="/contact/">Contact</a></li>
  <li><a href="https://gohugo.io" rel="external">Hugo</a></li>
</ul>
```

See the [menu templates] section for more information.

[menu templates]: /templates/menu/#menu-entry-parameters
[in front matter]: /content-management/menus/#define-in-front-matter
[in site configuration]: /content-management/menus/
