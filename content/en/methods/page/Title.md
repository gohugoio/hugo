---
title: Title
description: Returns the title of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
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

When a page is not backed by a file, the value returned by the `Title` method depends on the page [kind](g).

Page kind|Page title when the page is not backed by a file
:--|:--
home|site title
section|section name (capitalized and pluralized)
taxonomy|taxonomy name (capitalized and pluralized)
term|term name (capitalized and pluralized)

You can disable automatic capitalization and pluralization in your site configuration:

{{< code-toggle file=hugo >}}
capitalizeListTitles = false
pluralizeListTitles = false
{{< /code-toggle >}}

You can change the capitalization style in your site configuration to one of `ap`, `chicago`, `go`, `firstupper`, or `none`. For example:

{{< code-toggle file=hugo >}}
titleCaseStyle = "firstupper"
{{< /code-toggle >}}

See&nbsp;[details].

[details]: /configuration/all/#title-case-style
