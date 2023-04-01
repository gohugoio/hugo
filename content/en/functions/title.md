---
title: title
description: Converts all characters in the provided string to title case.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature:
  - "title INPUT"
  - "strings.Title INPUT"
relatedfuncs: []
---

```go-html-template
{{ title "BatMan"}}` → "Batman"
```

Can be combined in pipes. In the following snippet, the link text is cleaned up using `humanize` to remove dashes and `title` to convert the value of `$name` to Initial Caps.

```go-html-template
{{ range $name, $items := .Site.Taxonomies.categories }}
  <li><a href="{{ printf "%s/%s" "categories" ($name | urlize | lower) | absURL }}">{{ $name | humanize | title }} ({{ len $items }})</a></li>
{{ end }}
```

## Configure Title Case

The default is AP Stylebook, but you can [configure it](/getting-started/configuration/#configure-title-case).
