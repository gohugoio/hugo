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

With section, taxonomy, and term pages not backed by a file, the `Title` method returns the section name, capitalized and pluralized. You can disable these transformations by setting [`capitalizeListTitles`] and [`pluralizeListTitles`] in your site configuration. For example:

{{< code-toggle file=hugo >}}
capitalizeListTitles = false
pluralizeListTitles = false
{{< /code-toggle >}}

You can change the capitalization style in your site configuration to one of `ap`, `chicago`, `go`, `firstupper`, or `none`. For example:

{{< code-toggle file=hugo >}}
titleCaseStyle = "firstupper"
{{< /code-toggle >}}

 See [details].

[`capitalizeListTitles`]: /getting-started/configuration/#capitalizelisttitles
[`pluralizeListTitles`]: /getting-started/configuration/#pluralizelisttitles
[details]: /getting-started/configuration/#configure-title-case
