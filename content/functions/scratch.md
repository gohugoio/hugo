---
title: .Scratch
description: Acts as a "scratchpad" to allow for writable page- or shortcode-scoped variables.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [iteration]
categories: [functions]
menu:
  docs:
    parent: "functions"
toc:
signature: []
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: [/extras/scratch/,/doc/scratch/]
---

In most cases you can do well without `Scratch`, but there are some use cases that aren't solvable with Go's templates without `Scratch`'s help, due to scoping issues.

{{% note %}}
See [this Go issue](https://github.com/golang/go/issues/10608) for the main motivation behind Scratch.
{{% /note %}}

`Scratch` is added to both `Page` and `Shortcode` -- with following methods:

* `Set` and `Add` takes a `key` and the `value` to add.
* `Get` returns the `value` for the `key` given.
* `SetInMap` takes a `key`, `mapKey` and `value`
* `GetSortedMapValues` returns array of values from `key` sorted by `mapKey`

`Set` and `SetInMap` can store values of any type.

For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

The scope of the backing data is global for the given `Page` or `Shortcode`, and spans partial and shortcode includes.

Note that `.Scratch` from a shortcode will return the shortcode's `Scratch`, which in most cases is what you want. If you want to store it in the page scroped Scratch, then use `.Page.Scratch`.

## Sample usage

The usage is best illustrated with some samples:

```
{{ $.Scratch.Add "a1" 12 }}
{{ $.Scratch.Get "a1" }} {{/* => 12 */}}
{{ $.Scratch.Add "a1" 1 }}
{{ $.Scratch.Get "a1" }} // {{/* => 13 */}}

{{ $.Scratch.Add "a2" "AB" }}
{{ $.Scratch.Get "a2" }} {{/* => AB */}}
{{ $.Scratch.Add "a2" "CD" }}
{{ $.Scratch.Get "a2" }} {{/* => ABCD */}}

{{ $.Scratch.Add "l1" (slice "A" "B") }}
{{ $.Scratch.Get "l1" }} {{/* => [A B]  */}}
{{ $.Scratch.Add "l1" (slice "C" "D") }}
{{ $.Scratch.Get "l1" }} {{/* => [A B C D] */}}

{{ $.Scratch.Set "v1" 123 }}
{{ $.Scratch.Get "v1" }}  {{/* => 123 */}}

{{ $.Scratch.SetInMap "a3" "b" "XX" }}
{{ $.Scratch.SetInMap "a3" "a" "AA" }}
{{ $.Scratch.SetInMap "a3" "c" "CC" }}
{{ $.Scratch.SetInMap "a3" "b" "BB" }}
{{ $.Scratch.GetSortedMapValues "a3" }} {{/* => []interface {}{"AA", "BB", "CC"} */}}
```

{{% note %}}
The examples above uses the special `$` variable, which refers to the top-level node. This is the behavior you most likely want, and will help remove some confusion when using `Scratch` inside page range loops -- and you start inadvertently calling the wrong `Scratch`. But there may be use cases for `{{ .Scratch.Add "key" "some value" }}`.
{{% /note %}}


[pagevars]: /variables/page/
