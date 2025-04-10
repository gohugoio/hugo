---
title: compare.Ge
description: Returns the boolean truth of arg1 >= arg2 && arg1 >= arg3.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [ge]
    returnType: bool
    signatures: ['compare.Ge ARG1 ARG2 [ARG...]']
aliases: [/functions/ge]
---

```go-html-template
{{ ge 1 1 }} → true
{{ ge 1 2 }} → false
{{ ge 2 1 }} → true

{{ ge 1 1 1 }} → true
{{ ge 1 1 2 }} → false
{{ ge 1 2 1 }} → false
{{ ge 1 2 2 }} → false

{{ ge 2 1 1 }} → true
{{ ge 2 1 2 }} → true
{{ ge 2 2 1 }} → true
```

Use the `compare.Ge` function to compare other data types as well:

```go-html-template
{{ ge "ab" "a" }} → true
{{ ge time.Now (time.AsTime "1964-12-30") }} → true
{{ ge true false }} → true
```
