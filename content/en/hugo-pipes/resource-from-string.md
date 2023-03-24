---
title: FromString
linkTitle: Resource from String
description: Creates a resource from a string.
date: 2018-07-14
publishdate: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 90
weight: 90
sections_weight: 90
signature: ["resources.FromString TARGET_PATH CONTENT"]
---

## Usage

It is possible to create a resource directly from the template using `resources.FromString` which takes two arguments, the target path for the created resource and the given content string.

The following example creates a resource file containing localized variables for every project's languages.

```go-html-template
{{ $string := (printf "var rootURL = '%s'; var apiURL = '%s';" (absURL "/") (.Param "API_URL")) }}
{{ $targetPath := "js/vars.js" }}
{{ $vars := $string | resources.FromString $targetPath }}
{{ $global := resources.Get "js/global.js" | resources.Minify }}

<script src="{{ $vars.Permalink }}"></script>
<script src="{{ $global.Permalink }}"></script>
```
