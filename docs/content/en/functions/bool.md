---
title: bool
linktitle: bool
description: Creates a `bool` from the argument passed into the function.
date: 2023-01-28
publishdate: 2023-01-28
lastmod: 2023-01-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings,boolean,bool]
signature: ["bool INPUT"]
workson: []
hugoversion:
relatedfuncs: [truth]
deprecated: false
aliases: []
---

Useful for turning ints, strings, and nil into booleans.

```
{{ bool "true" }} → true
{{ bool "false" }} → false

{{ bool "TRUE" }} → true
{{ bool "FALSE" }} → false

{{ truth "t" }} → true
{{ truth "f" }} → false

{{ truth "T" }} → true
{{ truth "F" }} → false

{{ bool "1" }} → true
{{ bool "0" }} → false

{{ bool 1 }} → true
{{ bool 0 }} → false

{{ bool true }} → true
{{ bool false }} → false

{{ bool nil }} → false
```

This function will throw a type-casting error for most other types or strings. For less strict behavior, see [`truth`](/functions/truth).
