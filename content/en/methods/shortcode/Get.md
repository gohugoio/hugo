---
title: Get
description: Returns the value of the given parameter.
categories: []
keywords: []
action:
  related:
    - methods/shortcode/IsNamedParams
    - methods/shortcode/Params
  returnType: any
  signatures: [SHORTCODE.Get PARAM]
toc: true
---

Specify the parameter by position or by name. When calling a shortcode within markdown, use either positional or named parameters, but not both.

{{% note %}}
Some shortcodes support positional parameters, some support named parameters, and others support both. Refer to the shortcode's documentation for usage details.
{{% /note %}}

## Positional parameters

This shortcode call uses positional parameters:

{{< code file="content/about.md" lang=md copy=false >}}
{{</* myshortcode "Hello" "world" */>}}
{{< /code >}}

To retrieve parameters by position:

{{< code file="layouts/shortcodes/myshortcode.html" lang=go-html-template copy=false >}}
{{ printf "%s %s." (.Get 0) (.Get 1) }} → Hello world.
{{< /code >}}

## Named parameters

This shortcode call uses named parameters:

{{< code file="content/about.md" lang=md copy=false >}}
{{</* myshortcode greeting="Hello" firstName="world" */>}}
{{< /code >}}

To retrieve parameters by name:

{{< code file="layouts/shortcodes/myshortcode.html" lang=go-html-template copy=false >}}
{{ printf "%s %s." (.Get "greeting") (.Get "firstName") }} → Hello world.
{{< /code >}}

{{% note %}}
Parameter names are case-sensitive.
{{% /note %}}
