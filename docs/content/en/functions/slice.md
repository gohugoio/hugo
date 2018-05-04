---
title: slice
# linktitle: slice
description: Creates a slice (array) of all passed arguments.
godocref:
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

One use case is the concatenation of elements in combination with the [`delimit` function][]:

{{< code file="slice.html" >}}
{{ delimit (slice "foo" "bar" "buzz") ", " }}
<!-- returns the string "foo, bar, buzz" -->
{{< /code >}}


[`delimit` function]: /functions/delimit/
