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

By example, let's use [MathJax] to render a LaTeX mathematical expression:

[MathJax]: https://www.mathjax.org/

{{< code file="contents/physics/lesson-1.md" lang=markdown >}}
Albert Einstein’s theory of special relativity expresses
the fact that mass and energy are the same physical entity
and can be changed into each other.

{{</* math */>}}
$$
E=mc^2
$$
{{</* /math */>}}

In the equation, the increased relativistic mass (m) of a
body times the speed of light squared (c2) is equal to
the kinetic energy (E) of that body.
{{< /code >}}

The shortcode is simple:

{{< code file="layouts/shortcodes/math.html" lang=go-html-template >}}
{{ trim .Inner "\r\n" }}
{{< /code >}}

Now we can selectively load the required CSS and JavaScript on pages that call the "math" shortcode:


{{< code file="layouts/baseof.html" lang=go-html-template >}}
<head>
  ...
  {{ if .HasShortcode "math" }}
    <script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
    <script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
  {{ end }}
  ...
</head>
{{< /code >}}
