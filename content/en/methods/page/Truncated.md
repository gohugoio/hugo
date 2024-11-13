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

You can define a [summary] manually, in front matter, or automatically. A manual summary takes precedence over a front matter summary, and a front matter summary takes precedence over an automatic summary.

[summary]: /content-management/summaries/

The `Truncated` method returns `true` if the content length exceeds the summary length. This is useful for conditionally rendering a "read more" link:

```go-html-template
{{ range .Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ .Summary }}
  {{ if .Truncated }}
    <a href="{{ .RelPermalink }}">Read more...</a>
  {{ end }}
{{ end }}
```

{{% note %}}
The `Truncated` method returns `false` if you define the summary in front matter.
{{% /note %}}
