---
title: ExecuteAsTemplate
linkTitle: Resource from Template
description: Creates a resource from a template
categories: [asset management]
keywords: []
menu:
  docs:
    parent: pipes
    weight: 80
weight: 80
signature: ["resources.ExecuteAsTemplate TARGET_PATH CONTEXT RESOURCE"]
---

## Usage

In order to use Hugo Pipes function on an asset file containing Go Template magic the function `resources.ExecuteAsTemplate` must be used.

The function takes three arguments: the target path for the created resource, the template context, and the resource object. The target path is used to cache the result.

```go-html-template
// assets/sass/template.scss
$backgroundColor: {{ .Param "backgroundColor" }};
$textColor: {{ .Param "textColor" }};
body{
  background-color:$backgroundColor;
  color: $textColor;
}
// [...]
```

```go-html-template
{{ $sassTemplate := resources.Get "sass/template.scss" }}
{{ $style := $sassTemplate | resources.ExecuteAsTemplate "main.scss" . | resources.ToCSS }}
```
