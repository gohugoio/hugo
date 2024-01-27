---
title: FromString
linkTitle: Resource from string
description: Creates a resource from a string.
categories: [asset management]
keywords: []
menu:
  docs:
    parent: hugo-pipes
    weight: 110
weight: 110
action:
  aliases: []
  returnType: resource.Resource
  signatures: [resources.FromString TARGETPATH STRING]
---

## Usage

It is possible to create a resource directly from the template using `resources.FromString` which takes two arguments, the target path for the created resource and the given content string.

The result is cached using the target path as the cache key.

The following example creates a resource file containing localized variables for every project's languages.

```go-html-template
{{ $string := (printf "var rootURL = '%s'; var apiURL = '%s';" (absURL "/") (.Param "API_URL")) }}
{{ $targetPath := "js/vars.js" }}
{{ $vars := $string | resources.FromString $targetPath }}
{{ $global := resources.Get "js/global.js" | resources.Minify }}

<script src="{{ $vars.Permalink }}"></script>
<script src="{{ $global.Permalink }}"></script>
```
