---
title: Get
description: Returns the value of the given argument.
categories: []
keywords: []
action:
  related:
    - methods/shortcode/IsNamedParams
    - methods/shortcode/Params
  returnType: any
  signatures: [SHORTCODE.Get ARG]
toc: true
---

Specify the argument by position or by name. When calling a shortcode within Markdown, use either positional or named argument, but not both.

{{% note %}}
Some shortcodes support positional arguments, some support named arguments, and others support both. Refer to the shortcode's documentation for usage details.
{{% /note %}}

## Positional arguments

This shortcode call uses positional arguments:

{{< code file=content/about.md lang=md >}}
{{</* myshortcode "Hello" "world" */>}}
{{< /code >}}

To retrieve arguments by position:

{{< code file=layouts/shortcodes/myshortcode.html  >}}
{{ printf "%s %s." (.Get 0) (.Get 1) }} → Hello world.
{{< /code >}}

## Named arguments

This shortcode call uses named arguments:

{{< code file=content/about.md lang=md >}}
{{</* myshortcode greeting="Hello" firstName="world" */>}}
{{< /code >}}

To retrieve arguments by name:

{{< code file=layouts/shortcodes/myshortcode.html  >}}
{{ printf "%s %s." (.Get "greeting") (.Get "firstName") }} → Hello world.
{{< /code >}}

{{% note %}}
Argument names are case-sensitive.
{{% /note %}}
