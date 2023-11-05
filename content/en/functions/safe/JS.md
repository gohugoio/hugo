---
title: safe.JS
description: Declares the provided string as a safe JavaScript string.
categories: []
keywords: []
action:
  aliases: [safeJS]
  related:
    - functions/safe/CSS
    - functions/safe/HTML
    - functions/safe/HTMLAttr
    - functions/safe/JSStr
    - functions/safe/URL
  returnType: template.JS
  signatures: [safe.JS INPUT]
aliases: [/functions/safejs]
---

In this context, *safe* means the string encapsulates a known safe EcmaScript5 Expression (e.g., `(x + y * z())`).

Template authors are responsible for ensuring that typed expressions do not break the intended precedence and that there is no statement/expression ambiguity as when passing an expression like `{ foo:bar() }\n['foo']()`, which is both a valid expression and a valid program with a very different meaning.

Example: Given `hash = "619c16f"` defined in the front matter of your `.md` file:

* <span class="good">`<script>var form_{{ .Params.hash | safeJS }};…</script>` &rarr; `<script>var form_619c16f;…</script>`</span>
* <span class="bad">`<script>var form_{{ .Params.hash }};…</script>` &rarr; `<script>var form_"619c16f";…</script>`</span>
