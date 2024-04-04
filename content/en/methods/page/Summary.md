---
title: Summary
description: Returns the content summary of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Truncated
    - methods/page/Description
  returnType: template.HTML
  signatures: [PAGE.Summary]
---

<!-- Do not remove the manual summary divider below. -->
<!-- If you do, you will break its first literal usage on this page. -->
<!--more-->

There are three ways to define the [content summary]:

1. Let Hugo create the summary based on the first 70 words. You can change the number of words by setting the `summaryLength` in your site configuration.
2. Manually split the content with a `<!--more-->` tag in Markdown. Everything before the tag is included in the summary.
3. Create a `summary` field in front matter.

To list the pages in a section with a summary beneath each link:

```go-html-template
{{ range .Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ .Summary }}
{{ end }}
```

[content summary]: /content-management/summaries/
