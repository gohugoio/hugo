---
title: HasShortcode
description: Reports whether the given shortcode is called by the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: bool
    signatures: [PAGE.HasShortcode NAME]
---

By example, let's use [Plotly] to render a chart:

[Plotly]: https://plotly.com/javascript/

```text {file="content/example.md"}
{{</* plotly */>}}
{
  "data": [
    {
      "x": ["giraffes", "orangutans", "monkeys"],
      "y": [20, 14, 23],
      "type": "bar"
    }
  ],
}
{{</* /plotly */>}}
```

The shortcode is simple:

```go-html-template {file="layouts/_shortcodes/plotly.html"}
{{ $id := printf "plotly-%02d" .Ordinal }}
<div id="{{ $id }}"></div>
<script>
  Plotly.newPlot(document.getElementById({{ $id }}), {{ .Inner | safeJS }});
</script>
```

Now we can selectively load the required JavaScript on pages that call the "plotly" shortcode:

```go-html-template {file="layouts/baseof.html"}
<head>
  ...
  {{ if .HasShortcode "plotly" }}
    <script src="https://cdn.plot.ly/plotly-2.28.0.min.js"></script>
  {{ end }}
  ...
</head>
```
