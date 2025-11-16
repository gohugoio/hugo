---
title: Layout
description: Returns the layout for the given page as defined in front matter.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [PAGE.Layout]
---

Specify the `layout` field in front matter to target a particular template. See&nbsp;[details].

[details]: /templates/lookup-order/#target-a-template

{{< code-toggle file=content/contact.md fm=true >}}
title = 'Contact'
layout = 'contact'
{{< /code-toggle >}}

Hugo will render the page using contact.html.

```text
layouts/
├── baseof.html
├── contact.html
├── home.html
├── page.html
├── section.html
├── taxonomy.html
└── term.html
```

Although rarely used within a template, you can access the value with:

```go-html-template
{{ .Layout }}
```

The `Layout` method returns an empty string if the `layout` field in front matter is not defined.
