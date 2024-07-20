---
title: fmt.Errorf
description: Log an ERROR from a template.
categories: []
keywords: []
action:
  aliases: [errorf]
  related:
    - functions/fmt/Erroridf
    - functions/fmt/Warnf
    - functions/fmt/Warnidf
  returnType: string
  signatures: ['fmt.Errorf FORMAT [INPUT]']
aliases: [/functions/errorf]
---

{{% include "functions/fmt/_common/fmt-layout.md" %}}

The `errorf` function evaluates the format string, then prints the result to the ERROR log and fails the build.

```go-html-template
{{ errorf "The %q shortcode requires a src argument. See %s" .Name .Position }}
```

Use the [`erroridf`] function to allow optional suppression of specific errors.

[`erroridf`]: /functions/fmt/erroridf/
