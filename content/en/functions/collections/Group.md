---
title: collections.Group
description: Groups the given page collection by the given key.
categories: []
keywords: []
action:
  aliases: [group]
  related:
    - functions/collections/Dictionary
    - functions/collections/IndexFunction
    - functions/collections/IsSet
    - functions/collections/Where
  returnType: any
  signatures: [PAGES | collections.Group KEY]
aliases: [/functions/group]
---

```go-html-template
{{ $new := .Site.RegularPages | first 10 | group "New" }}
{{ $old := .Site.RegularPages | last 10 | group "Old" }}
{{ $groups := slice $new $old }}
{{ range $groups }}
  <h3>{{ .Key }}{{/* Prints "New", "Old" */}}</h3>
  <ul>
    {{ range .Pages }}
      <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
      </li>
    {{ end }}
  </ul>
{{ end }}
```

The page group you get from `group` is of the same type you get from the built-in [group methods](/templates/lists#group-content) in Hugo. The above example can be [paginated](/templates/pagination/#list-paginator-pages).
