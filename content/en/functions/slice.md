---
title: slice
# linktitle: slice
description: Creates a slice (array) of all passed arguments.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [slice, array, interface]
signature: ["slice ITEM..."]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
toc: false
---

One use case is the concatenation of elements in combination with the [`delimit` function]:

{{< code file="slice.html" >}}
{{ $sliceOfStrings := slice "foo" "bar" "buzz" }}
<!-- returns the slice [ "foo", "bar", "buzz"] -->
{{ delimit ($sliceOfStrings) ", " }}
<!-- returns the string "foo, bar, buzz" -->
{{< /code >}}


[`delimit` function]: /functions/delimit/
