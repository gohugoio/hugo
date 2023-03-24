---
title: Concat
linkTitle: Concatenating assets
description: Bundle any number of assets into one resource.
date: 2018-07-14
publishdate: 2018-07-14
lastmod: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 60
weight: 60
sections_weight: 60
signature: ["resources.Concat TARGET_PATH SLICE_RESOURCES"]
---

## Usage

Asset files of the same MIME type can be bundled into one resource using `resources.Concat` which takes two arguments, the target path for the created resource bundle and a slice of resource objects to be concatenated.

```go-html-template
{{ $plugins := resources.Get "js/plugins.js" }}
{{ $global := resources.Get "js/global.js" }}
{{ $js := slice $plugins $global | resources.Concat "js/bundle.js" }}
```
