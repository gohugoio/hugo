---
title: int
description: Creates an `int` from the argument passed into the function.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings,integers]
signature: ["int INPUT"]
relatedfuncs: []
---

Useful for turning strings into numbers.

```go-html-template
{{ int "123" }} → 123
```

{{% note "Usage Note" %}}
If the input string is supposed to represent a decimal number, and if it has
leading 0's, then those 0's will have to be removed before passing the string
to the `int` function, else that string will be tried to be parsed as an octal
number representation.

The `strings.TrimLeft` can be used for this purpose.

```go-html-template
{{ int ("0987" | strings.TrimLeft "0") }}
{{ int ("00987" | strings.TrimLeft "0") }}
```

### Explanation

The `int` function eventually calls the `ParseInt` function from the Go library
`strconv`.

From its [documentation](https://golang.org/pkg/strconv/#ParseInt):

> the base is implied by the string's prefix: base 16 for "0x", base 8 for "0",
> and base 10 otherwise.
{{% /note %}}
