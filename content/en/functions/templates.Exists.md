---
title: templates.Exists
description: "Checks whether a template file exists under the given path relative to the `layouts` directory."
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: templates
relatedFuncs: []
signature:
  - templates.Exists PATH

---

A template file is any file living below the `layouts` directories of either the project or any of its theme components including partials and shortcodes.

The function is particularly handy with dynamic path. The following example ensures the build will not break on a `.Type` missing its dedicated `header` partial.

```go-html-template
{{ $partialPath := printf "headers/%s.html" .Type }}
{{ if templates.Exists ( printf "partials/%s" $partialPath ) }}
  {{ partial $partialPath . }}
{{ else }}
  {{ partial "headers/default.html" . }}
{{ end }}
```
