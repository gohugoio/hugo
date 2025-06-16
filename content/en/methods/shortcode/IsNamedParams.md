---
title: IsNamedParams
description: Reports whether the shortcode call uses named arguments.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: bool
    signatures: [SHORTCODE.IsNamedParams]
---

To support both positional and named arguments when calling a shortcode, use the `IsNamedParams` method to determine how the shortcode was called.

With this shortcode template:

```go-html-template {file="layouts/_shortcodes/myshortcode.html"}
{{ if .IsNamedParams }}
  {{ printf "%s %s." (.Get "greeting") (.Get "firstName") }}
{{ else }}
  {{ printf "%s %s." (.Get 0) (.Get 1) }}
{{ end }}
```

Both of these calls return the same value:

```text {file="content/about.md"}
{{</* myshortcode greeting="Hello" firstName="world" */>}}
{{</* myshortcode "Hello" "world" */>}}
```
