---
title: Truncated
description: Reports whether the content length exceeds the summary length.
categories: []
keywords: []
action:
  related:
    - methods/page/Summary
  returnType: bool
  signatures: [PAGE.Truncated]
---

There are three ways to define the [content summary]:

1. Let Hugo create the summary based on the first 70 words. You can change the number of words by setting the `summaryLength` in your site configuration.
2. Manually split the content with a `<--more-->` tag in Markdown. Everything before the tag is included in the summary.
3. Create a `summary` field in front matter.

{{% note %}}
The `Truncated` method returns `false` if you define the summary in front matter.
{{% /note %}}

The `Truncated` method returns `true` if the content length exceeds the summary length. This is useful for rendering a "read more" link:

```go-html-template
{{ range .Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ .Summary }}
  {{ if .Truncated }}
    <a href="{{ .RelPermalink }}">Read more...</a>
  {{ end }}
{{ end }}
```

[content summary]: /content-management/summaries/
