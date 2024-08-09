---
title: HasMenuCurrent
description: Reports whether the given Page object matches the Page object associated with one of the child menu entries under the given menu entry in the given menu.
categories: []
keywords: []
action:
  related:
    - methods/page/IsMenuCurrent
  returnType: bool
  signatures: [PAGE.HasMenuCurrent MENU MENUENTRY]
aliases: [/functions/hasmenucurrent]
---

If the `Page` object associated with the menu entry is a section, this method also returns `true` for any descendant of that section.

```go-html-template
{{ $currentPage := . }}
{{ range site.Menus.main }}
  {{ if $currentPage.IsMenuCurrent .Menu . }}
    <a class="active" aria-current="page" href="{{ .URL }}">{{ .Name }}</a>
  {{ else if $currentPage.HasMenuCurrent .Menu . }}
    <a class="ancestor" aria-current="true" href="{{ .URL }}">{{ .Name }}</a>
  {{ else }}
    <a href="{{ .URL }}">{{ .Name }}</a>
  {{ end }}
{{ end }}
```

See [menu templates] for a complete example.

[menu templates]: /templates/menu/#example
