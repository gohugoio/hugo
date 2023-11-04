---
title: .HasMenuCurrent
description: Reports whether the given page object matches the page object associated with one of the child menu entries under the given menu entry in the given menu.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: bool
  signatures: [PAGE.HasMenuCurrent MENU MENUENTRY]
relatedFunctions:
  - .HasMenuCurrent
  - .IsMenuCurrent
---

If the page object associated with the menu entry is a section, this method also returns `true` for any descendant of that section.

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

[menu templates]: /templates/menu-templates/#example
