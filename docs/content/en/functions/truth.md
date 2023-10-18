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

```
{{ truth "true" }} → true
{{ truth "false" }} → true

{{ truth "TRUE" }} → true
{{ truth "FALSE" }} → true

{{ truth "t" }} → true
{{ truth "f" }} → true

{{ truth "T" }} → true
{{ truth "F" }} → true

{{ truth "1" }} → true
{{ truth "0" }} → true

{{ truth 1 }} → true
{{ truth 0 }} → false

{{ truth true }} → true
{{ truth false }} → false

{{ truth nil }} → false

{{ truth "cheese" }} → true
{{ truth 1.67 }} → true
```

This function will not throw an error. For more strict behavior, see [`bool`](/functions/bool).
