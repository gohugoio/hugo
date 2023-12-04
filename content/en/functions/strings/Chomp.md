---
title: strings.Chomp
description: Returns the given string, removing all trailing newline characters and carriage returns.
categories: []
keywords: []
action:
  aliases: [chomp]
  related:
    - functions/strings/Trim
    - functions/strings/TrimLeft
    - functions/strings/TrimPrefix
    - functions/strings/TrimRight
    - functions/strings/TrimSuffix
  returnType: any
  signatures: [strings.Chomp STRING]
aliases: [/functions/chomp]
---

If the argument is of type `template.HTML`, returns `template.HTML`, else returns a `string`.

```go-html-template
{{ chomp | "foo\n" }} → foo
{{ chomp | "foo\n\n" }} → foo

{{ chomp | "foo\r\n" }} → foo
{{ chomp | "foo\r\n\r\n" }} → foo
```
