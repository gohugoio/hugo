---
title: .Render
description: Takes a view to apply when rendering content.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace:
relatedFuncs: []
signature: 
  - .Render LAYOUT
---

The view is an alternative layout and should be a file name that points to a template in one of the locations specified in the documentation for [Content Views](/templates/views).

This function is only available when applied to a single piece of content within a [list context].

This example could render a piece of content using the content view located at `/layouts/_default/summary.html`:

```go-html-template
{{ range .Pages }}
  {{ .Render "summary" }}
{{ end }}
```

[list context]: /templates/lists/
