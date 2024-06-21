---
title: KeyName
description: Returns the `identifier` property of the given menu entry, falling back to its `name` property. 
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [MENUENTRY.KeyName]
---

In this menu definition, the second entry does not contain an `identifier`, so the `Identifier` method returns its `name` property instead:

{{< code-toggle file=hugo >}}
[[menus.main]]
identifier = 'about'
name = 'About'
pageRef = '/about'
weight = 10

[[menus.main]]
name = 'Contact'
pageRef = '/contact'
weight = 20
{{< /code-toggle >}}

This example uses the `KeyName` method when querying the translation table on a multilingual site, falling back the `name` property if a matching key in the translation table does not exist:

```go-html-template
<ul>
  {{ range .Site.Menus.main }}
    <li><a href="{{ .URL }}">{{ or (T (.KeyName | lower)) .Name }}</a></li>
  {{ end }}
</ul>
```

In the example above, we need to pass the value returned by `.KeyName` through the [`lower`] function because the keys in the translation table are lowercase.

[`lower`]: /functions/strings/tolower/
