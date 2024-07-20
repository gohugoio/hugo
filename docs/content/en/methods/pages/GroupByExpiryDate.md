---
title: GroupByExpiryDate
description: Returns the given page collection grouped by expiration date in descending order.
categories: []
keywords: []
action:
  related:
    - methods/pages/GroupByDate
    - methods/pages/GroupByLastMod
    - methods/pages/GroupByParamDate
    - methods/pages/GroupByPublishDate
  returnType: page.PagesGroup
  signatures: ['PAGES.GroupByExpiryDate LAYOUT [SORT]']
---

When grouping by expiration date, the value is determined by your [site configuration], defaulting to the `expiryDate` field in front matter.

The [layout string] has the same format as the layout string for the [`time.Format`] function. The resulting group key is [localized] for language and region.

[`time.Format`]: /functions/time/format/
[layout string]: #layout-string
[localized]: /getting-started/glossary/#localization
[site configuration]: /getting-started/configuration/#configure-dates

{{% include "methods/pages/_common/group-sort-order.md" %}}

To group content by year and month:

```go-html-template
{{ range .Pages.GroupByExpiryDate "January 2006" }}
  <p>{{ .Key }}</p>
  <ul>
    {{ range .Pages }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

To sort the groups in ascending order:

```go-html-template
{{ range .Pages.GroupByExpiryDate "January 2006" "asc" }}
  <p>{{ .Key }}</p>
  <ul>
    {{ range .Pages }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

The pages within each group will also be sorted by expiration date, either ascending or descending depending on your grouping option. To sort the pages within each group, use one of the sorting methods. For example, to sort the pages within each group by title:

```go-html-template
{{ range .Pages.GroupByExpiryDate "January 2006" }}
  <p>{{ .Key }}</p>
  <ul>
    {{ range .Pages.ByTitle }}
      <li><a href="{{ .RelPermalink }}">{{ .Title }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

## Layout string

{{% include "functions/_common/time-layout-string.md" %}}
