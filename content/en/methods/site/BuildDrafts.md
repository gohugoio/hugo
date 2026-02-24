---
title: BuildDrafts
description: Reports reports whether draft publishing is enabled for the current build.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: bool
    signatures: [SITE.BuildDrafts]
expiryDate: '2028-02-18' # deprecated 2026-02-18 in v0.156.0
---

{{< deprecated-in 0.156.0 >}}
See [details](https://discourse.gohugo.io/t/56732).
{{< /deprecated-in >}}

By default, draft pages are not published when building a site. You can change this behavior with a command line flag:

```sh
hugo build --buildDrafts
```

Or by setting `buildDrafts` to `true` in your project configuration:

{{< code-toggle file=hugo >}}
buildDrafts = true
{{< /code-toggle >}}

Use the `BuildDrafts` method on a `Site` object to determine the current configuration:

```go-html-template
{{ .Site.BuildDrafts }} → true
```
