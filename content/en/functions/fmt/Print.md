---
title: fmt.Print
linkTitle: print
description: Prints the default representation of the given arguments using the standard `fmt.Print` function.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [print]
  returnType: string
  signatures: [fmt.Print INPUT]
relatedFunctions:
  - fmt.Print
  - fmt.Printf
  - fmt.Println
aliases: [/functions/print]
---

```go-html-template
{{ print "foo" }} → "foo"
{{ print "foo" "bar" }} → "foobar"
{{ print (slice 1 2 3) }} → [1 2 3]
```
