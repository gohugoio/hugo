---
title: safeHTMLAttr
description: Declares the provided string as a safe HTML attribute.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["safeHTMLAttr INPUT"]
relatedfuncs: []
---

Given a site configuration that contains this menu entry:

{{< code-toggle file="config" >}}
[[menu.main]]
  name = "IRC"
  url = "irc://irc.freenode.net/#golang"
{{< /code-toggle >}}

Attempting to use the `url` value directly in an attribute:

```go-html-template
{{ range site.Menus.main }}
  <a href="{{ .URL }}">{{ .Name }}</a>
{{ end }}
``` 

Will produce:

```html
<a href="#ZgotmplZ">IRC</a>
```

`ZgotmplZ` is a special value, inserted by Go's [template/html] package, that indicates that unsafe content reached a CSS or URL context.

To override the safety check, use the `safeHTMLAttr` function:

```go-html-template
{{ range site.Menus.main }}
  <a {{ printf "href=%q" .URL | safeHTMLAttr }}>{{ .Name }}</a>
{{ end }}
``` 

[template/html]: https://pkg.go.dev/html/template
