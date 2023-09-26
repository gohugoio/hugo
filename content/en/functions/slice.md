---
title: slice
description: Creates a slice (array) of all passed arguments.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: collections
relatedFuncs:
  - collections.Append
  - collections.Apply
  - collections.Delimit
  - collections.In
  - collections.Reverse
  - collections.Seq
  - collections.Slice
signature:
  - collections.Slice ITEM...
  - slice ITEM...
---

One use case is the concatenation of elements in combination with the [`delimit` function]:

{{< code file="slice.html" >}}
{{ $sliceOfStrings := slice "foo" "bar" "buzz" }}
<!-- returns the slice [ "foo", "bar", "buzz"] -->
{{ delimit ($sliceOfStrings) ", " }}
<!-- returns the string "foo, bar, buzz" -->
{{< /code >}}


[`delimit` function]: /functions/delimit/
