---
title: truth
linktitle: truth
description: Creates a `bool` from the truthyness of the argument passed into the function
date: 2023-01-28
publishdate: 2023-01-28
lastmod: 2023-01-28
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings,boolean,bool,truthy,falsey]
signature: ["truth INPUT"]
workson: []
hugoversion:
relatedfuncs: [bool]
deprecated: false
aliases: []
---

Useful for turning different types into booleans based on their [truthy-ness](https://developer.mozilla.org/en-US/docs/Glossary/Truthy).

It follows the same rules as [`bool`](/functions/bool), but with increased flexibility.

```
{{ truth "true" }} → true
{{ truth "false" }} → false

{{ truth "TRUE" }} → true
{{ truth "FALSE" }} → false

{{ truth "t" }} → true
{{ truth "f" }} → false

{{ truth "T" }} → true
{{ truth "F" }} → false

{{ truth "1" }} → true
{{ truth "0" }} → false

{{ truth 1 }} → true
{{ truth 0 }} → false

{{ truth true }} → true
{{ truth false }} → false

{{ truth nil }} → false

{{ truth "cheese" }} → true
{{ truth 1.67 }} → true
```

This function will not throw an error. For more strict behavior, see [`bool`](/functions/bool).
