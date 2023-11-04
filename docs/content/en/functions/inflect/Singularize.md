---
title: inflect.Singularize
linkTitle: singularize
description: Converts a word according to a set of common English singularization rules.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [singularize]
  returnType: string
  signatures: [inflect.Singularize INPUT]
relatedFunctions:
  - inflect.Humanize
  - inflect.Pluralize
  - inflect.Singularize
aliases: [/functions/singularize]
---

```go-html-template
{{ "cats" | singularize }} → "cat"
```

See also the `.Data.Singular` [taxonomy variable](/variables/taxonomy/) for singularizing taxonomy names.
