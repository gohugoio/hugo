---
title: Params
description: Returns a collection of the shortcode parameters.
categories: []
keywords: []
action:
  related:
    - methods/shortcode/Get
  returnType: any
  signatures: [SHORTCODE.Params]
---

When you call a shortcode using positional parameters, the `Params` method returns a slice.

{{< code file=content/about.md lang=md >}}
{{</* myshortcode "Hello" "world" */>}}
{{< /code >}}

{{< code file=layouts/shortcodes/myshortcode.html  >}}
{{ index .Params 0 }} → Hello
{{ index .Params 1 }} → world
{{< /code >}}

When you call a shortcode using named parameters, the `Params` method returns a map.

{{< code file=content/about.md lang=md >}}
{{</* myshortcode greeting="Hello" name="world" */>}}
{{< /code >}}

{{< code file=layouts/shortcodes/myshortcode.html  >}}
{{ .Params.greeting }} → Hello
{{ .Params.name }} → world
{{< /code >}}
