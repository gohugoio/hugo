---
title: Mathematics in Markdown
linkTitle: Mathematics
description: Include mathematical equations and expressions in your Markdown using LaTeX or TeX typesetting syntax.
categories: [content management]
keywords: [chemical,chemistry,latex,math,mathjax,tex,typesetting]
menu:
  docs:
    parent: content-management
    weight: 270
weight: 270
toc: true
math: true
---

{{< new-in 0.122.0 >}}

\[
\begin{aligned}
KL(\hat{y} || y) &= \sum_{c=1}^{M}\hat{y}_c \log{\frac{\hat{y}_c}{y_c}} \\
JS(\hat{y} || y) &= \frac{1}{2}(KL(y||\frac{y+\hat{y}}{2}) + KL(\hat{y}||\frac{y+\hat{y}}{2}))
\end{aligned}
\]

## Overview

Mathematical equations and expressions authored in [LaTeX] or [TeX] are common in academic and scientific publications. Your browser typically renders this mathematical markup using an open-source JavaScript display engine such as [MathJax] or [KaTeX].

For example, this is the mathematical markup for the equations displayed at the top of this page:

```text
\[
\begin{aligned}
KL(\hat{y} || y) &= \sum_{c=1}^{M}\hat{y}_c \log{\frac{\hat{y}_c}{y_c}} \\
JS(\hat{y} || y) &= \frac{1}{2}(KL(y||\frac{y+\hat{y}}{2}) + KL(\hat{y}||\frac{y+\hat{y}}{2}))
\end{aligned}
\]
```

Equations and expressions can be displayed inline with other text, or as standalone blocks. Block presentation is also known as "display" mode.

Whether an equation or expression appears inline, or as a block, depends on the delimiters that surround the mathematical markup. Delimiters are defined in pairs, where each pair consists of an opening and closing delimiter. The opening and closing delimiters may be the same, or different. Common delimiter pairs are shown in [Step 1].

The approach described below avoids reliance on platform-specific features like shortcodes or code block render hooks. Instead, it utilizes a standardized markup format for mathematical equations and expressions, compatible with the rendering engines used by GitHub, GitLab, [Microsoft VS Code], [Obsidian], [Typora], and others.

## Setup

Follow these instructions to include mathematical equations and expressions in your Markdown using LaTeX or TeX typesetting syntax.

###### Step 1

Enable and configure the Goldmark [passthrough extension] in your site configuration. The passthrough extension preserves raw Markdown within delimited snippets of text, including the delimiters themselves.

{{< code-toggle file=hugo copy=true >}}
[markup.goldmark.extensions.passthrough]
enable = true

[markup.goldmark.extensions.passthrough.delimiters]
block = [['\[', '\]'], ['$$', '$$']]
inline = [['\(', '\)']]

[params]
math = true
{{< /code-toggle >}}

The configuration above enables mathematical rendering on every page unless you set the `math` parameter to `false` in front matter. To enable mathematical rendering as needed, set the `math` parameter to `false` in your site configuration, and set the `math` parameter to `true` in front matter. Use this parameter in your base template as shown in [Step 3].

{{% note %}}
The configuration above precludes the use of the `$...$` delimiter pair for inline equations. Although you can add this delimiter pair to the configuration and JavaScript, you will need to double-escape the `$` symbol when used outside of math contexts to avoid unintended formatting.

See the [inline delimiters](#inline-delimiters) section for details.
{{% /note %}}

To disable passthrough of inline snippets, omit the `inline` key from the configuration:

{{< code-toggle file=hugo >}}
[markup.goldmark.extensions.passthrough.delimiters]
block = [['\[', '\]'], ['$$', '$$']]
{{< /code-toggle >}}

You can define your own opening and closing delimiters, provided they match the delimiters that you set in [Step 2].

{{< code-toggle file=hugo >}}
[markup.goldmark.extensions.passthrough.delimiters]
block = [['@@', '@@']]
inline = [['@', '@']]
{{< /code-toggle >}}

###### Step 2

Create a partial template to load MathJax or KaTeX. The example below loads MathJax, or you can use KaTeX as described in the [engines](#engines) section.

{{< code file=layouts/partials/math.html copy=true >}}
<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-chtml.js"></script>
<script>
  MathJax = {
    tex: {
      displayMath: [['\\[', '\\]'], ['$$', '$$']],  // block
      inlineMath: [['\\(', '\\)']]                  // inline
    }
  };
</script>
{{< /code >}}

The delimiters above must match the delimiters in your site configuration.

###### Step 3

Conditionally call the partial template from the base template.

{{< code file=layouts/_default/baseof.html >}}
<head>
  ...
  {{ if .Param "math" }}
    {{ partialCached "math.html" . }}
  {{ end }}
  ...
</head>
{{< /code >}}

The example above loads the partial template if you have set the `math` parameter in front matter to `true`. If you have not set the `math` parameter in front matter, the conditional statement falls back to the `math` parameter in your site configuration.

###### Step 4

Include mathematical equations and expressions in your Markdown using LaTeX or TeX typesetting syntax.

{{< code file=content/math-examples.md copy=true >}}
This is an inline \(a^*=x-b^*\) equation.

These are block equations:

\[a^*=x-b^*\]

\[ a^*=x-b^* \]

\[
a^*=x-b^*
\]

These are block equations using alternate delimiters:

$$a^*=x-b^*$$

$$ a^*=x-b^* $$

$$
a^*=x-b^*
$$
{{< /code >}}

If you set the `math` parameter to `false` in your site configuration, you must set the `math` parameter to `true` in front matter. For example:

{{< code-toggle file=content/math-examples.md fm=true >}}
title = 'Math examples'
date = 2024-01-24T18:09:49-08:00
[params]
math = true
{{< /code-toggle >}}

## Inline delimiters

The configuration, JavaScript, and examples above use the `\(...\)` delimiter pair for inline equations. The `$...$` delimiter pair is a common alternative, but using it may result in unintended formatting if you use the `$` symbol outside of math contexts.

If you add the `$...$` delimiter pair to your configuration and JavaScript, you must double-escape the `$` when outside of math contexts, regardless of whether mathematical rendering is enabled on the page. For example:

```text
A \\$5 bill _saved_ is a \\$5 bill _earned_.
```

{{% note %}}
If you use the `$...$` delimiter pair for inline equations, and occasionally use the&nbsp;`$`&nbsp;symbol outside of math contexts, you must use MathJax instead of KaTeX to avoid unintended formatting caused by [this KaTeX limitation](https://github.com/KaTeX/KaTeX/issues/437).
{{% /note %}}

## Engines

MathJax and KaTeX are open-source JavaScript display engines. Both engines are fast, but at the time of this writing MathJax v3.2.2 is slightly faster than KaTeX v0.16.9.

{{% note %}}
If you use the `$...$` delimiter pair for inline equations, and occasionally use the&nbsp;`$`&nbsp;symbol outside of math contexts, you must use MathJax instead of KaTeX to avoid unintended formatting caused by [this KaTeX limitation](https://github.com/KaTeX/KaTeX/issues/437).

See the [inline delimiters](#inline-delimiters) section for details.
{{% /note %}}

To use KaTeX instead of MathJax, replace the partial template from [Step 2] with this:

{{< code file=layouts/partials/math.html copy=true >}}
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" integrity="sha384-n8MVd4RsNIU0tAv4ct0nTaAbDJwPJzDEaqSD1odI+WdtXRGWt2kTvGFasHpSy3SV" crossorigin="anonymous">
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js" integrity="sha384-XjKyOOlGwcjNTAIQHIpgOno0Hl1YQqzUOEleOLALmuqehneUG+vnGctmUb0ZY0l8" crossorigin="anonymous"></script>
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js" integrity="sha384-+VBxd3r6XgURycqtZ117nYw44OOcIax56Z4dCRWbxyPt0Koah1uHoK0o4+/RRE05" crossorigin="anonymous"></script>
<script>
  document.addEventListener("DOMContentLoaded", function() {
    renderMathInElement(document.body, {
      delimiters: [
        {left: '\\[', right: '\\]', display: true},   // block
        {left: '$$', right: '$$', display: true},     // block
        {left: '\\(', right: '\\)', display: false},  // inline
      ],
      throwOnError : false
    });
  });
</script>
{{< /code >}}

The delimiters above must match the delimiters in your site configuration.

## Chemistry

Both MathJax and KaTeX provide support for chemical equations. For example:

```text
$$C_p[\ce{H2O(l)}] = \pu{75.3 J // mol K}$$
```

$$C_p[\ce{H2O(l)}] = \pu{75.3 J // mol K}$$

As shown in [Step 2] above, MathJax supports chemical equations without additional configuration. To add chemistry support to KaTeX, enable the mhchem extension as described in the KaTeX [documentation](https://katex.org/docs/libs).

[KaTeX]: https://katex.org/
[LaTeX]: https://www.latex-project.org/
[MathJax]: https://www.mathjax.org/
[Microsoft VS Code]: https://code.visualstudio.com/
[Obsidian]: https://obsidian.md/
[Step 1]: #step-1
[Step 2]: #step-2
[Step 3]: #step-3
[TeX]: https://en.wikipedia.org/wiki/TeX
[Typora]: https://typora.io/
[passthrough extension]: https://github.com/gohugoio/hugo-goldmark-extensions
