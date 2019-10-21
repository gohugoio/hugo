---
title: title
# linktitle:
description: Converts all characters in the provided string to title case.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["title INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---


```
{{title "BatMan"}}` → "Batman"
```

Can be combined in pipes. In the following snippet, the link text is cleaned up using `humanize` to remove dashes and `title` to convert the value of `$name` to Initial Caps.

```
{{ range $name, $items := .Site.Taxonomies.categories }}
    <li><a href="{{ printf "%s/%s" "categories" ($name | urlize | lower) | absURL }}">{{ $name | humanize | title }} ({{ len $items }})</a></li>
{{ end }}
```

## Configure Title Case

The default is AP Stylebook, but you can [configure it](/getting-started/configuration/#configure-title-case).
