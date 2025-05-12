---
title: math.MaxInt64
description: Returns the maximum value for a signed 64-bit integer.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: int64
    signatures: [math.MaxInt64]
---

{{< new-in 0.147.3 />}}

```go-html-template
{{ math.MaxInt64 }} → 9223372036854775807
```

This function is helpful for simulating a loop that continues indefinitely until a break condition is met. For example:

```go-html-template
{{ range math.MaxInt64 }}
  {{ if eq . 42 }}
    {{ break }}
  {{ end }}
{{ end }}
```
