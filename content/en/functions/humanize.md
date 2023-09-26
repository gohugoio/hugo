---
title: humanize
description: Returns the humanized version of an argument with the first letter capitalized.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: inflect
relatedFuncs:
  - inflect.Humanize
  - inflect.Pluralize
  - inflect.Singularize
signature:
 - inflect.Humanize INPUT
 - humanize INPUT
---

If the input is either an int64 value or the string representation of an integer, humanize returns the number with the proper ordinal appended.


```go-html-template
{{ humanize "my-first-post" }} → "My first post"
{{ humanize "myCamelPost" }} → "My camel post"
{{ humanize "52" }} → "52nd"
{{ humanize 103 }} → "103rd"
```
