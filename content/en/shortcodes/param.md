---
title: Param
description: Insert a parameter from front matter or site configuration into your content using the param shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
---

{{% note %}}
To override Hugo's embedded `param` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl param %}}
{{% /note %}}

The `param` shortcode renders a parameter from front matter, falling back to a site parameter of the same name. The shortcode throws an error if the parameter does not exist.

{{< code file=example.md lang=text >}}
---
title: Example
date: 2025-01-15T23:29:46-08:00
params:
  color: red
  size: medium
---

We found a {{</* param "color" */>}} shirt.
{{< /code >}}

Hugo renders this to:

```html
<p>We found a red shirt.</p>
```

Access nested values by [chaining](g) the [identifiers](g):

```text
{{</* param my.nested.param */>}}
```
