---
title: Title
description: Returns the title of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/LinkTitle
  returnType: string
  signatures: [PAGE.Title]
---

With pages backed by a file, the `Title` method returns the `title` field as defined in front matter:

{{< code-toggle file=content/about.md fm=true >}}
title = 'About us'
{{< /code-toggle >}}

```go-html-template
{{ .Title }} → About us
```

With section pages not backed by a file, the `Title` method returns the section name, pluralized and converted to title case.

To disable [pluralization]:

{{< code-toggle file=hugo >}}
pluralizeListTitles = false
{{< /code-toggle >}}

To change the [title case style], specify one of `ap`, `chicago`, `go`, `firstupper`, or `none`:

{{< code-toggle file=hugo >}}
titleCaseStyle = "ap"
{{< /code-toggle >}}

[pluralization]: /functions/inflect/pluralize
[title case style]: /getting-started/configuration/#configure-title-case
