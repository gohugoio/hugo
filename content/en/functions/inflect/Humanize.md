---
title: inflect.Humanize
description: Returns the humanized version of the input with the first letter capitalized.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [humanize]
    returnType: string
    signatures: [inflect.Humanize INPUT]
---

```go-html-template
{{ humanize "my-first-post" }} → My first post
{{ humanize "myCamelPost" }} → My camel post
```

If the input is an integer or a string representation of an integer, humanize returns the number with the proper ordinal appended.

```go-html-template
{{ humanize "52" }} → 52nd
{{ humanize 103 }} → 103rd
```
