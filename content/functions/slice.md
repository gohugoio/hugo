---
title: slice
linktitle: slice
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [slice, array, interface]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
toc: false
needsexamples: true
---

`slice` allows you to create an array (`[]interface{}`) of all arguments that you pass to this function.

One use case is the concatenation of elements in combination with the [`delimit` function][]:

{{% code "slice.html" %}}
```html
{{ delimit (slice "foo" "bar" "buzz") ", " }}
<!-- returns the string "foo, bar, buzz" -->
```
{{% /code %}}


[`delimit` function]: /functions/delimit/