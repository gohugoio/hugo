---
title: int
linktitle: int
description: Creates an `int` from the argument passed into the function.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings,integers]
signature: ["int INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Useful for turning strings into numbers.

```
{{ int "123" }} → 123
```

{{% note "Usage Note" %}}
If the input string is supposed to represent a decimal number, and if it has
leading 0's, then those 0's will have to be removed before passing the string
to the `int` function, else that string will be tried to be parsed as an octal
number representation.

The [`strings.TrimLeft` function](/functions/strings.trimleft/) can be used for
this purpose.

```
{{ int ("0987" | strings.TrimLeft "0") }}
{{ int ("00987" | strings.TrimLeft "0") }}
```

**Explanation**

The `int` function eventually calls the `ParseInt` function from the Go library
`strconv`.

From its [documentation](https://golang.org/pkg/strconv/#ParseInt):

> the base is implied by the string's prefix: base 16 for "0x", base 8 for "0",
> and base 10 otherwise.
{{% /note %}}
