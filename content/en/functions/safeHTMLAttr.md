---
title: safeHTMLAttr
description: Declares the provided string as a safe HTML attribute.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: safe
relatedFuncs:
  - safe.CSS
  - safe.HTML
  - safe.HTMLAttr
  - safe.JS
  - safe.JSStr
  - safe.URL
signature:
  - safe.HTMLAttr INPUT
  - safeHTMLAttr INPUT
---

Given a site configuration that contains this menu entry:

{{< code-toggle file="hugo" >}}
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

To indicate that the HTML attribute is safe:

```go-html-template
{{ range site.Menus.main }}
  <a {{ printf "href=%q" .URL | safeHTMLAttr }}>{{ .Name }}</a>
{{ end }}
```

{{% note %}}
As demonstrated above, you must pass the HTML attribute name _and_ value through the function. Applying `safeHTMLAttr` to the attribute value has no effect.
{{% /note %}}

[template/html]: https://pkg.go.dev/html/template
