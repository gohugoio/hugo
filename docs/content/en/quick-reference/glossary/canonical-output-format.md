---
title: canonical output format
---

The _canonical output format_ is the [_output format_](g) for the current page where the format's [`rel`][] property is set to `canonical` in your project configuration, if such a format exists. If there is only one _output format_ for the current page and it is a predefined format, Hugo automatically treats it as the _canonical output format_ regardless of whether its `rel` property is set to `canonical`. Custom output formats are not subject to this rule; `rel` must be explicitly set to `canonical`.

  By default, `html` is the only predefined _output format_ with this setting; the `rel` property for all others is set to `alternate`. If two or more _output formats_ for the current page have their `rel` property set to `canonical`, the _canonical output format_ is the first one specified in:

  - The [`outputs`][outputs_front_matter] front matter field of the current page, or
  - The [`outputs`][outputs_project_config] section of your project configuration for the current [_page kind_](g).

  [`rel`]: /configuration/output-formats/#rel
  [outputs_front_matter]: /configuration/outputs/#outputs-per-page
  [outputs_project_config]: /configuration/outputs/#outputs-per-page-kind
