---
title: safe.HTMLAttr
description: Declares the provided string as a safe HTML attribute.
categories: []
keywords: []
action:
  aliases: [safeHTMLAttr]
  related:
    - functions/safe/CSS
    - functions/safe/HTML
    - functions/safe/JS
    - functions/safe/JSStr
    - functions/safe/URL
  returnType: template.HTMLAttr
  signatures: [safe.HTMLAttr INPUT]
aliases: [/functions/safehtmlattr]
---

Given a site configuration that contains this menu entry:

{{< code-toggle file=hugo >}}
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
