---
title: fmt.Print
description: Prints the default representation of the given arguments using the standard `fmt.Print` function.
categories: []
keywords: []
action:
  aliases: [print]
  related:
    - functions/fmt/Printf
    - functions/fmt/Println
  returnType: string
  signatures: [fmt.Print INPUT]
aliases: [/functions/print]
---

```go-html-template
{{ print "foo" }} → foo
{{ print "foo" "bar" }} → foobar
{{ print (slice 1 2 3) }} → [1 2 3]
```
