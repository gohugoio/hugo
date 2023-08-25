---
title: .Get
description: Accesses positional and ordered parameters in shortcode declaration.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [shortcodes]
signature: [".Get INDEX", ".Get KEY"]
relatedfuncs: []
---

`.Get` is specifically used when creating your own [shortcode template][sc], to access the [positional and named](/templates/shortcode-templates/#positional-vs-named-parameters) parameters passed to it. When used with a numeric INDEX, it queries positional parameters (starting with 0). With a string KEY, it queries named parameters.

When accessing named or positional parameters that do not exist, `.Get` returns an empty string instead of interrupting the build. This allows you to chain `.Get` with `if`, `with`, `default` or `cond` to check for parameter existence. For example:

```go-html-template
{{ $quality := default "100" (.Get 1) }}
```

[sc]: /templates/shortcode-templates/
