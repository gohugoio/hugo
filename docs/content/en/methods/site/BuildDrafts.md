---
title: BuildDrafts
description: Reports whether the current build includes draft pages.
categories: []
keywords: []
action:
  related: []
  returnType: bool
  signatures: [SITE.BuildDrafts]
---

By default, draft pages are not published when building a site. You can change this behavior with a command line flag:

```sh
hugo --buildDrafts
```

Or by setting `buildDrafts` to `true` in your site configuration:

{{< code-toggle file=hugo >}}
buildDrafts = true
{{< /code-toggle >}}

Use the `BuildDrafts` method on a `Site` object to determine the current configuration:

```go-html-template
{{ .Site.BuildDrafts }} → true
```
