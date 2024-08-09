---
title: Menu
description: Returns the identifier of the menu that contains the given menu entry.
categories: []
keywords: []
action:
  related:
    - methods/page/IsMenuCurrent
    - methods/page/HasMenuCurrent
  returnType: string
  signatures: [MENUENTRY.Menu]
---

```go-html-template
{{ range .Site.Menus.main }}
  {{ .Menu }} → main
{{ end }}
```

Use this method with the [`IsMenuCurrent`] and [`HasMenuCurrent`] methods on a `Page` object to set "active" and "ancestor" classes on a rendered entry. See [this example].

[`HasMenuCurrent`]: /methods/page/hasmenucurrent/
[`IsMenuCurrent`]: /methods/page/ismenucurrent/
[this example]: /templates/menu/#example
