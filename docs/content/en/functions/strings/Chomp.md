---
title: strings.Chomp
description: Returns the given string, removing all trailing newline characters and carriage returns.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [chomp]
    returnType: any
    signatures: [strings.Chomp STRING]
---

If the argument is of type `template.HTML`, returns `template.HTML`, else returns a `string`.

```go-html-template
{{ chomp "foo\n" }} → foo
{{ chomp "foo\n\n" }} → foo

{{ chomp "foo\r\n" }} → foo
{{ chomp "foo\r\n\r\n" }} → foo
```
