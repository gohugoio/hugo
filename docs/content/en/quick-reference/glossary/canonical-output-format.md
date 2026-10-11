---
title: canonical output format
---

The _canonical output format_ is the [_output format_](g) for the current page where the format's [`rel`][] property is set to `canonical` in your project configuration, if such a format exists. If there is only one output format for the current page and it is a predefined format, Hugo automatically treats it as the canonical output format regardless of whether its `rel` property is set to `canonical`. Custom output formats are not subject to this rule, so you must explicitly set their `rel` property to `canonical`.

  By default, `html` is the only predefined output format with this setting. The `rel` property for all others is set to `alternate`. If two or more output formats for the current page have their `rel` property set to `canonical`, the canonical output format is the first one specified in either of these locations:

  - The `outputs` field in the [front matter of the current page][]
  - The `outputs` section of your [project configuration][] for the current [_page kind_](g)

  [`rel`]: /configuration/output-formats/#rel
  [front matter of the current page]: /configuration/outputs/#outputs-per-page
  [project configuration]: /configuration/outputs/#outputs-per-page-kind
