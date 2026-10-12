---
title: partial decorator
params:
  reference: /templates/partial-decorators/
---

A _partial decorator_ is a specific type of [_partial template_](g) that functions as a [_wrapper component_](g). While a standard partial template renders data within a fixed template, a decorator uses composition to enclose an entire block of content. It uses the [`templates.Inner`][] function as a placeholder to define exactly where that external content should be injected within the wrapper's layout.

  [`templates.Inner`]: /functions/templates/inner/
