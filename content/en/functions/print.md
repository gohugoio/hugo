---
title: print
description: Prints the default representation of the given arguments using the standard `fmt.Print` function.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: fmt
relatedFuncs:
  - fmt.Print
  - fmt.Printf
  - fmt.Println
signature:
  - fmt.Print INPUT
  - print INPUT
---

```go-html-template
{{ print "foo" }} → "foo"
{{ print "foo" "bar" }} → "foobar"
{{ print (slice 1 2 3) }} → [1 2 3]
```
