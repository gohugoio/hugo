---
title: Params
description: Returns a collection of the shortcode arguments.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: any
    signatures: [SHORTCODE.Params]
---

When you call a shortcode using positional arguments, the `Params` method returns a slice.

```text {file="content/about.md"}
{{</* myshortcode "Hello" "world" */>}}
```

```go-html-template {file="layouts/_shortcodes/myshortcode.html"}
{{ index .Params 0 }} → Hello
{{ index .Params 1 }} → world
```

When you call a shortcode using named arguments, the `Params` method returns a map.

```text {file="content/about.md"}
{{</* myshortcode greeting="Hello" name="world" */>}}
```

```go-html-template {file="layouts/_shortcodes/myshortcode.html"}
{{ .Params.greeting }} → Hello
{{ .Params.name }} → world
```
