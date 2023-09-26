---
title: strings.TrimRight
description: Returns a slice of a given string with all trailing characters contained in the cutset removed.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
signature:
  - strings.TrimRight CUTSET STRING
---

Given the string `"abba"`, trailing `"a"`'s can be removed a follows:

```go-html-template
{{ strings.TrimRight "a" "abba" }} → "abb"
```

Numbers can be handled as well:

```go-html-template
{{ strings.TrimRight 12 1221341221 }} → "122134"
```
