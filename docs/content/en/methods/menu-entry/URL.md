---
title: URL
description: Returns the relative permalink of the page associated with the given menu entry, else its `url` property.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [MENUENTRY.URL]
---

For menu entries associated with a page, the `URL` method returns the page's [`RelPermalink`], otherwise it returns the entry's `url` property.


```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
```

[`RelPermalink`]: /methods/page/relpermalink/
