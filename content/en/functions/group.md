---
title: group
description: "`group` groups a list of pages."
date: 2018-09-14
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [collections]
signature: ["PAGES | group KEY"]
hugoversion: "0.49"
---

{{< code file="layouts/partials/groups.html" >}}
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
{{< /code >}}



The page group you get from `group` is of the same type you get from the built-in [group methods](/templates/lists#group-content) in Hugo. The above example can even be [paginated](/templates/pagination/#list-paginator-pages).




