---
title: Creating a resource from template
linkTitle: Resource from Template
description: Hugo Pipes allows the creation of a resource from an asset file using Go Template.
date: 2018-07-14
publishdate: 2018-07-14
lastmod: 2018-07-14
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 80
weight: 80
sections_weight: 80
draft: false
---

In order to use Hugo Pipes function on an asset file containing Go Template magic the function `resources.ExecuteAsTemplate` must be used.

The function takes three arguments, the resource object, the resource target path and the template context.

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