---
title: merge
description: "`merge` deep merges two maps and returns the resulting map."
date: 2018-09-14
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [dictionary]
signature: ["$params :=  merge $default_params $user_params", $params := $default_params | merge $user_params]
workson: []
hugoversion: "0.49"
relatedfuncs: [append, reflect.IsMap, reflect.IsSlice]
aliases: []
---

An example merging two maps.

```go-html-template
{{ $default_params := dict "color" "blue" "width" "50%" "height" "25%" }}
{{ $user_params := dict "color" "red" "extra" (dict "duration" 2) }}
{{ $params := $default_params | merge $user_params }}
```

Resulting __$params__:

```
"color": "blue"
"extra":
  "duration": 2
"height": "25%"
"icon": "mail"
"width": "50%"
```

{{% note %}}
  Regardless of depth, merging only applies to maps. For slices, use [append]({{< ref "functions/append" >}})
{{% /note %}}


