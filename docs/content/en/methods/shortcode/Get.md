---
title: Get
description: Returns the value of the given argument.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: any
    signatures: [SHORTCODE.Get ARG]
---

Specify the argument by position or by name. When calling a shortcode within Markdown, use either positional or named argument, but not both.

> [!note]
> Some shortcodes support positional arguments, some support named arguments, and others support both. Refer to the shortcode's documentation for usage details.

## Positional arguments

This shortcode call uses positional arguments:

```text {file="content/about.md"}
{{</* myshortcode "Hello" "world" */>}}
```

To retrieve arguments by position:

```go-html-template {file="layouts/_shortcodes/myshortcode.html"}
{{ printf "%s %s." (.Get 0) (.Get 1) }} → Hello world.
```

## Named arguments

This shortcode call uses named arguments:

```text {file="content/about.md"}
{{</* myshortcode greeting="Hello" firstName="world" */>}}
```

To retrieve arguments by name:

```go-html-template {file="layouts/_shortcodes/myshortcode.html"}
{{ printf "%s %s." (.Get "greeting") (.Get "firstName") }} → Hello world.
```

> [!note]
> Argument names are case-sensitive.
