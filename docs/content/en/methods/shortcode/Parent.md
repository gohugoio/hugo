---
title: Parent
description:  Returns the parent shortcode context in nested shortcodes.
categories: []
keywords: []
action:
  related: []
  returnType: hugolib.ShortcodeWithPage
  signatures: [SHORTCODE.Parent]
---

This is useful for inheritance of common shortcode arguments from the root.

In this contrived example, the "greeting" shortcode is the parent, and the "now" shortcode is child.

{{< code file=content/welcome.md lang=md >}}
{{</* greeting dateFormat="Jan 2, 2006" */>}}
Welcome. Today is {{</* now */>}}.
{{</* /greeting */>}}
{{< /code >}}

{{< code file=layouts/shortcodes/greeting.html  >}}
<div class="greeting">
  {{ trim .Inner "\r\n" | .Page.RenderString }}
</div>
{{< /code >}}

{{< code file=layouts/shortcodes/now.html  >}}
{{- $dateFormat := "January 2, 2006 15:04:05" }}

{{- with .Params }}
  {{- with .dateFormat }}
    {{- $dateFormat = . }}
  {{- end }}
{{- else }}
  {{- with .Parent.Params }}
    {{- with .dateFormat }}
      {{- $dateFormat = . }}
    {{- end }}
  {{- end }}
{{- end }}

{{- now | time.Format $dateFormat -}}
{{< /code >}}

The "now" shortcode formats the current time using:

1. The `dateFormat` argument passed to the "now" shortcode, if present
2. The `dateFormat` argument passed to the "greeting" shortcode, if present
3. The default layout string defined at the top of the shortcode
