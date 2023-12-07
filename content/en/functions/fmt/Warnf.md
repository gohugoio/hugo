---
title: fmt.Warnf
description: Log a WARNING from a template.
categories: []
keywords: []
action:
  aliases: [warnf]
  related:
    - functions/fmt/Errorf
    - functions/fmt/Erroridf
  returnType: string
  signatures: ['fmt.Warnf FORMAT [INPUT]']
aliases: [/functions/warnf]
---

{{% include "functions/fmt/_common/fmt-layout.md" %}}

The `warnf` function evaluates the format string, then prints the result to the WARNING log. Hugo prints each unique message once to avoid flooding the log with duplicate warnings.

```go-html-template
{{ warnf "The %q shortcode was unable to find %s. See %s" .Name $file .Position }}
```

To prevent supression of duplicate messages when using `warnf` for debugging, make each message unique with the [`math.Counter`] function. For example:


```go-html-template
{{ range site.RegularPages }}
  {{ .Section | warnf "%#[2]v [%[1]d]" math.Counter }}
{{ end }}
```

[`math.Counter`]: /functions/math/counter
