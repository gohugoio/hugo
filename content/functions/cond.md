---
title: "cond"
date: 2017-09-08
description: "Return one of two arguments, depending on the value of a third argument."
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["cond CONTROL VAR1 VAR2"]
aliases: [/functions/cond/]
hugoversion: 0.27
relatedfuncs: [default]
toc: false
draft: false
needsexamples: false
---

`cond` returns *VAR1* if *CONTROL* is true, or *VAR2* if it is not.

Example:

```
{{ cond (eq (len $geese) 1) "goose" "geese" }}
```

Would emit "goose" if the `$geese` array has exactly 1 item, or "geese" otherwise.
