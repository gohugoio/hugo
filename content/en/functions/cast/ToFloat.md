---
title: cast.ToFloat
linkTitle: float
description: Casts a value to a decimal (base 10) floating point value.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [float]
  returnType: float64
  signatures: [cast.ToFloat INPUT]
relatedFunctions:
  - cast.ToFloat
  - cast.ToInt
  - cast.ToString
aliases: [/functions/float]
---

With a decimal (base 10) input:

```go-html-template
{{ float 11 }} → 11 (float64)
{{ float "11" }} → 11 (float64)

{{ float 11.1 }} → 11.1 (float64)
{{ float "11.1" }} → 11.1 (float64)

{{ float 11.9 }} → 11.9 (float64)
{{ float "11.9" }} → 11.9 (float64)
```

With a binary (base 2) input:

```go-html-template
{{ float 0b11 }} → 3 (float64)
```

With an octal (base 8) input (use either notation):

```go-html-template
{{ float 011 }} → 9 (float64)
{{ float "011" }} → 11 (float64)

{{ float 0o11 }} → 9 (float64)
```

With a hexadecimal (base 16) input:

```go-html-template
{{ float 0x11 }} → 17 (float64)
```
