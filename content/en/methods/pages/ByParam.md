---
title: ByParam
description: Returns the given page collection sorted by the given parameter in ascending order.
categories: []
keywords: []
action:
  related:
    - methods/pages/ByTitle
    - methods/pages/ByLinkTitle
  returnType: page.Pages
  signatures: [PAGES.ByParam PARAM]
---

If the given parameter is not present in front matter, Hugo will use the matching parameter in your site configuration if present.

```go-html-template
{{ range .Pages.ByParam "author" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To sort in descending order:

```go-html-template
{{ range (.Pages.ByParam "author").Reverse }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

If the targeted parameter is nested, access the field using dot notation:

```go-html-template
{{ range .Pages.ByParam "author.last_name" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
