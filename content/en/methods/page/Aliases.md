---
title: Aliases
description: Returns the URL aliases as defined in front matter.
categories: []
keywords: []
action:
  related: []
  returnType: '[]string'
  signatures: [PAGE.Aliases]
---

The `Aliases` method on a `Page` object returns the URL [aliases] as defined in front matter.

For example:

{{< code-toggle file=content/about.md fm=true >}}
title = 'About'
aliases = ['/old-url','/really-old-url']
{{< /code-toggle >}}

To list the aliases:

```go-html-template
{{ range .Aliases }}
  {{ . }}
{{ end }}
```

[aliases]: /content-management/urls/#aliases
