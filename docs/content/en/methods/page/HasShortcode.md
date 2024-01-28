---
title: HasShortcode
description: Reports whether the given shortcode is called by the given page.
categories: []
keywords: []
action:
  related: []
  returnType: bool
  signatures: [PAGE.HasShortcode NAME]
---

By example, let's use [Plotly] to render a chart:

[Plotly]: https://plotly.com/javascript/

{{< code file=contents/example.md lang=markdown >}}
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
{{< /code >}}

The shortcode is simple:

{{< code file=layouts/shortcodes/plotly.html  >}}
{{ $id := printf "plotly-%02d" .Ordinal }}
<div id="{{ $id }}"></div>
<script>
  Plotly.newPlot(document.getElementById({{ $id }}), {{ .Inner | safeJS }});
</script>
{{< /code >}}

Now we can selectively load the required JavaScript on pages that call the "plotly" shortcode:

{{< code file=layouts/baseof.html  >}}
<head>
  ...
  {{ if .HasShortcode "plotly" }}
    <script src="https://cdn.plot.ly/plotly-2.28.0.min.js"></script>
  {{ end }}
  ...
</head>
{{< /code >}}
